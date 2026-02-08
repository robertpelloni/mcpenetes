package ui

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"runtime"
	"sync"

	"github.com/tuannvm/mcpenetes/internal/client"
	"github.com/tuannvm/mcpenetes/internal/config"
	"github.com/tuannvm/mcpenetes/internal/core"
	"github.com/tuannvm/mcpenetes/internal/doctor"
	"github.com/tuannvm/mcpenetes/internal/log"
	"github.com/tuannvm/mcpenetes/internal/registry"
	"github.com/tuannvm/mcpenetes/internal/registry/manager"
	"github.com/tuannvm/mcpenetes/internal/search"
	"github.com/tuannvm/mcpenetes/internal/util"
	"github.com/tuannvm/mcpenetes/internal/version"
)

//go:embed static/*
var staticFiles embed.FS

// Server represents the web UI server
type Server struct {
	Port int
}

// NewServer creates a new UI server instance
func NewServer(port int) *Server {
	return &Server{
		Port: port,
	}
}

// Start starts the web server
func (s *Server) Start() error {
	mux := http.NewServeMux()

	// Serve static files
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return err
	}
	mux.Handle("/", http.FileServer(http.FS(staticFS)))

	// API endpoints
	mux.HandleFunc("/api/data", s.handleGetData)
	mux.HandleFunc("/api/apply", s.handleApply)
	mux.HandleFunc("/api/search", s.handleSearch)
	mux.HandleFunc("/api/install", s.handleInstall)
	mux.HandleFunc("/api/server/update", s.handleUpdateServer)
	mux.HandleFunc("/api/server/remove", s.handleRemoveServer)
	mux.HandleFunc("/api/doctor", s.handleDoctor)
	mux.HandleFunc("/api/registry/add", s.handleAddRegistry)
	mux.HandleFunc("/api/registry/remove", s.handleRemoveRegistry)
	mux.HandleFunc("/api/server/inspect", s.handleInspectServer)
	mux.HandleFunc("/api/backups", s.handleGetBackups)
	mux.HandleFunc("/api/restore", s.handleRestoreBackup)
	mux.HandleFunc("/api/import", s.handleImportConfig)
	mux.HandleFunc("/api/logs", s.handleGetLogs)
	mux.HandleFunc("/api/clients/custom", s.handleCustomClients)
	mux.HandleFunc("/api/client/config", s.handleGetClientConfig)

	addr := fmt.Sprintf("localhost:%d", s.Port)
	log.Success("Starting Web UI at http://%s", addr)
	return http.ListenAndServe(addr, mux)
}

// Response structs
type ConfigDataResponse struct {
	Version    string                   `json:"version"`
	Clients    map[string]config.Client `json:"clients"`
	MCPServers map[string]config.MCPServer `json:"mcpServers"`
	Registries []config.Registry `json:"registries"`
}

type ApplyRequest struct {
	ClientNames []string `json:"clients"`
}

type InstallRequest struct {
	ServerID string            `json:"serverId"`
	Config   *config.MCPServer `json:"config,omitempty"` // Optional override
}

type UpdateServerRequest struct {
	ServerID string           `json:"serverId"`
	Config   config.MCPServer `json:"config"`
}

type RemoveServerRequest struct {
	ServerID string `json:"serverId"`
}

type AddRegistryRequest struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type RemoveRegistryRequest struct {
	Name string `json:"name"`
}

type InspectRequest struct {
	ServerID string           `json:"serverId"`
	Config   config.MCPServer `json:"config"`
}

type ApplyResponse struct {
	Results []core.ApplyResult `json:"results"`
}

type RestoreRequest struct {
	ClientName string `json:"client"`
	BackupFile string `json:"file"`
}

type ImportRequest struct {
	Config string `json:"config"`
}

type AddCustomClientRequest struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	ConfigFormat string `json:"configFormat"`
	ConfigKey    string `json:"configKey"`
	BaseDir      string `json:"baseDir"` // "home", "appdata", "userprofile"
	Path         string `json:"path"`    // Relative path
}

type RemoveCustomClientRequest struct {
	ID string `json:"id"`
}

func (s *Server) handleGetData(w http.ResponseWriter, r *http.Request) {
	cfg, err := config.LoadConfig()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading config: %v", err), http.StatusInternalServerError)
		return
	}

	mcpCfg, err := config.LoadMCPConfig()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading MCP config: %v", err), http.StatusInternalServerError)
		return
	}

	// Detect clients if none configured
	if len(cfg.Clients) == 0 {
		detected, err := util.DetectMCPClients()
		if err == nil {
			cfg.Clients = detected
		}
	}

	resp := ConfigDataResponse{
		Version:    version.Version,
		Clients:    cfg.Clients,
		MCPServers: mcpCfg.MCPServers,
		Registries: cfg.Registries,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleApply(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ApplyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading config: %v", err), http.StatusInternalServerError)
		return
	}

	mcpCfg, err := config.LoadMCPConfig()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading MCP config: %v", err), http.StatusInternalServerError)
		return
	}

	manager := core.NewManager(cfg, mcpCfg)
	var results []core.ApplyResult
	var mu sync.Mutex

	// If no clients specified, apply to all in config
	targetClients := req.ClientNames
	if len(targetClients) == 0 {
		for name := range cfg.Clients {
			targetClients = append(targetClients, name)
		}
	}

	// Group clients by config path to avoid race conditions
	// Map: expandedConfigPath -> []ClientInfo
	type ClientInfo struct {
		Name   string
		Config config.Client
	}
	groupedClients := make(map[string][]ClientInfo)

	for _, name := range targetClients {
		clientConf, ok := cfg.Clients[name]
		if !ok {
			continue
		}

		expandedPath, err := util.ExpandPath(clientConf.ConfigPath)
		if err != nil {
			// If expansion fails, just use original path as key (unlikely to collide if invalid)
			expandedPath = clientConf.ConfigPath
		}

		groupedClients[expandedPath] = append(groupedClients[expandedPath], ClientInfo{Name: name, Config: clientConf})
	}

	// Process each file group sequentially to prevent race conditions on the same file
	// We can still process different files concurrently if we wanted,
	// but for safety and simplicity, we'll do everything sequentially in this iteration.
	// If performance becomes an issue, we can parallelize the outer loop over `groupedClients`.

	for _, clientList := range groupedClients {
		for _, clientInfo := range clientList {
			res := manager.ApplyToClient(clientInfo.Name, clientInfo.Config)
			mu.Lock()
			results = append(results, res)
			mu.Unlock()
		}
	}

	type JSONResult struct {
		ClientName string `json:"clientName"`
		Success    bool   `json:"success"`
		BackupPath string `json:"backupPath"`
		Error      string `json:"error,omitempty"`
	}

	var jsonResults []JSONResult
	for _, res := range results {
		jr := JSONResult{
			ClientName: res.ClientName,
			Success:    res.Success,
			BackupPath: res.BackupPath,
		}
		if res.Error != nil {
			jr.Error = res.Error.Error()
		}
		jsonResults = append(jsonResults, jr)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"results": jsonResults})
}

func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	cfg, err := config.LoadConfig()
	if err != nil {
		http.Error(w, "Failed to load config", http.StatusInternalServerError)
		return
	}

	var allServers []registry.ServerData
	for _, reg := range cfg.Registries {
		servers, err := registry.FetchMCPServersWithCache(reg.URL, false)
		if err == nil {
			allServers = append(allServers, servers...)
		}
	}

	var filtered []registry.ServerData
	if query == "" {
		filtered = allServers
	} else {
		for _, s := range allServers {
			if util.CaseInsensitiveContains(s.Name, query) || util.CaseInsensitiveContains(s.Description, query) {
				filtered = append(filtered, s)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filtered)
}

func (s *Server) handleInstall(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req InstallRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ServerID == "" {
		http.Error(w, "Server ID is required", http.StatusBadRequest)
		return
	}

	// Pass the optional config override
	err := search.AddServerToMCPConfig(req.ServerID, nil, req.Config)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to install server: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Server added to configuration"})
}

func (s *Server) handleUpdateServer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req UpdateServerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ServerID == "" {
		http.Error(w, "Server ID is required", http.StatusBadRequest)
		return
	}

	mcpCfg, err := config.LoadMCPConfig()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading MCP config: %v", err), http.StatusInternalServerError)
		return
	}

	mcpCfg.MCPServers[req.ServerID] = req.Config

	if err := config.SaveMCPConfig(mcpCfg); err != nil {
		http.Error(w, fmt.Sprintf("Error saving MCP config: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Server configuration updated"})
}

func (s *Server) handleRemoveServer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RemoveServerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ServerID == "" {
		http.Error(w, "Server ID is required", http.StatusBadRequest)
		return
	}

	mcpCfg, err := config.LoadMCPConfig()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading MCP config: %v", err), http.StatusInternalServerError)
		return
	}

	if _, ok := mcpCfg.MCPServers[req.ServerID]; !ok {
		http.Error(w, "Server not found", http.StatusNotFound)
		return
	}

	delete(mcpCfg.MCPServers, req.ServerID)

	if err := config.SaveMCPConfig(mcpCfg); err != nil {
		http.Error(w, fmt.Sprintf("Error saving MCP config: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Server removed from configuration"})
}

func (s *Server) handleDoctor(w http.ResponseWriter, r *http.Request) {
	results := doctor.RunChecks()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (s *Server) handleAddRegistry(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AddRegistryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.URL == "" {
		http.Error(w, "Name and URL are required", http.StatusBadRequest)
		return
	}

	if err := manager.AddRegistry(req.Name, req.URL); err != nil {
		http.Error(w, fmt.Sprintf("Failed to add registry: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Registry added"})
}

func (s *Server) handleRemoveRegistry(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RemoveRegistryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Registry name is required", http.StatusBadRequest)
		return
	}

	if err := manager.RemoveRegistry(req.Name); err != nil {
		http.Error(w, fmt.Sprintf("Failed to remove registry: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Registry removed"})
}

func (s *Server) handleInspectServer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req InspectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	cmdStr := fmt.Sprintf("npx @modelcontextprotocol/inspector %s", req.Config.Command)
	for _, arg := range req.Config.Args {
		cmdStr += fmt.Sprintf(" %s", arg)
	}

	if len(req.Config.Env) > 0 {
		envStr := ""
		for k, v := range req.Config.Env {
			envStr += fmt.Sprintf("%s='%s' ", k, v)
		}
		cmdStr = envStr + cmdStr
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
		"command": cmdStr,
		"message": "To inspect this server, run the following command in your terminal:",
	})
}

func (s *Server) handleGetBackups(w http.ResponseWriter, r *http.Request) {
	cfg, err := config.LoadConfig()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading config: %v", err), http.StatusInternalServerError)
		return
	}

	// We don't strictly need MCPConfig for listing backups, but manager expects it.
	mcpCfg := &config.MCPConfig{}

	manager := core.NewManager(cfg, mcpCfg)
	backups, err := manager.ListBackups()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error listing backups: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(backups)
}

func (s *Server) handleRestoreBackup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RestoreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ClientName == "" || req.BackupFile == "" {
		http.Error(w, "Client name and backup file are required", http.StatusBadRequest)
		return
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading config: %v", err), http.StatusInternalServerError)
		return
	}
	mcpCfg := &config.MCPConfig{} // Not needed for restore

	manager := core.NewManager(cfg, mcpCfg)
	err = manager.RestoreClient(req.ClientName, req.BackupFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error restoring backup: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": fmt.Sprintf("Restored backup for %s", req.ClientName)})
}

func (s *Server) handleImportConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ImportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Config == "" {
		http.Error(w, "Config content is required", http.StatusBadRequest)
		return
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading config: %v", err), http.StatusInternalServerError)
		return
	}

	mcpCfg, err := config.LoadMCPConfig()
	if err != nil {
		// Create new if missing
		mcpCfg = &config.MCPConfig{MCPServers: make(map[string]config.MCPServer)}
	}

	manager := core.NewManager(cfg, mcpCfg)
	count, err := manager.ImportConfig(req.Config)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error importing config: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
		"message": fmt.Sprintf("Imported %d servers", count),
	})
}

func (s *Server) handleGetLogs(w http.ResponseWriter, r *http.Request) {
	logs := log.GetRecentLogs()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

func (s *Server) handleCustomClients(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		s.handleGetCustomClients(w, r)
	} else if r.Method == http.MethodPost {
		s.handleAddCustomClient(w, r)
	} else if r.Method == http.MethodDelete {
		s.handleRemoveCustomClient(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleGetCustomClients(w http.ResponseWriter, r *http.Request) {
	clients, err := client.LoadCustomClients()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading clients: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clients)
}

func (s *Server) handleAddCustomClient(w http.ResponseWriter, r *http.Request) {
	var req AddCustomClientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create definition
	def := client.ClientDefinition{
		ID:           req.ID,
		Name:         req.Name,
		ConfigFormat: client.ConfigFormatEnum(req.ConfigFormat),
		ConfigKey:    req.ConfigKey,
		Paths: map[string][]client.PathDefinition{
			runtime.GOOS: {
				{Base: client.BaseDirEnum(req.BaseDir), Path: req.Path},
			},
		},
	}

	if err := client.AddCustomClient(def); err != nil {
		http.Error(w, fmt.Sprintf("Error adding client: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Client added"})
}

func (s *Server) handleRemoveCustomClient(w http.ResponseWriter, r *http.Request) {
	var req RemoveCustomClientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := client.RemoveCustomClient(req.ID); err != nil {
		http.Error(w, fmt.Sprintf("Error removing client: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Client removed"})
}

func (s *Server) handleGetClientConfig(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("id")
	if clientID == "" {
		http.Error(w, "Client ID required", http.StatusBadRequest)
		return
	}

	// To be safe, we only read configs of detected/configured clients
	cfg, err := config.LoadConfig()
	if err != nil {
		http.Error(w, "Failed to load config", http.StatusInternalServerError)
		return
	}

	// If client is in config, use that path
	c, ok := cfg.Clients[clientID]
	if !ok {
		// Try detection
		detected, _ := util.DetectMCPClients()
		if d, found := detected[clientID]; found {
			c = config.Client{ConfigPath: d.ConfigPath}
			ok = true
		}
	}

	if !ok {
		http.Error(w, "Client not found or not configured", http.StatusNotFound)
		return
	}

	expandedPath, err := util.ExpandPath(c.ConfigPath)
	if err != nil {
		expandedPath = c.ConfigPath
	}

	content, err := os.ReadFile(expandedPath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "Config file does not exist yet", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Error reading file: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write(content)
}
