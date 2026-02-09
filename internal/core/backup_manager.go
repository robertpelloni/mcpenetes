package core

// RemoveBackup deletes a backup file for a given client
func (m *Manager) RemoveBackup(clientName, filename string) error {
	return m.Trans.DeleteBackup(clientName, filename)
}
