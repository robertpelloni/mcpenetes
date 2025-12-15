package ui

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"sync"

	"github.com/tuannvm/mcpenetes/internal/config"
	"github.com/tuannvm/mcpenetes/internal/core"
	"github.com/tuannvm/mcpenetes/internal/log"
	"github.com/tuannvm/mcpenetes/internal/registry"
	"github.com/tuannvm/mcpenetes/internal/search"
	"github.com/tuannvm/mcpenetes/internal/util"
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

	addr := fmt.Sprintf("localhost:%d", s.Port)
	log.Success("Starting Web UI at http://%s", addr)
	return http.ListenAndServe(addr, mux)
}

// Response structs
type ConfigDataResponse struct {
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

type ApplyResponse struct {
	Results []core.ApplyResult `json:"results"`
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
	var wg sync.WaitGroup

	// If no clients specified, apply to all in config
	targetClients := req.ClientNames
	if len(targetClients) == 0 {
		for name := range cfg.Clients {
			targetClients = append(targetClients, name)
		}
	}

	for _, name := range targetClients {
		clientConf, ok := cfg.Clients[name]
		if !ok {
			continue
		}

		wg.Add(1)
		go func(n string, c config.Client) {
			defer wg.Done()
			res := manager.ApplyToClient(n, c)
			mu.Lock()
			results = append(results, res)
			mu.Unlock()
		}(name, clientConf)
	}

	wg.Wait()

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
