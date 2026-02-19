package nidda

import (
	"encoding/json"
	"fmt"
	"os"
)

// SaveToFile saves the history to a JSON file.
func (m *NiddaManager) SaveToFile(filename string) error {
	data, err := json.MarshalIndent(m.History, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal history: %v", err)
	}
	return os.WriteFile(filename, data, 0644)
}

// LoadFromFile loads the history from a JSON file.
func (m *NiddaManager) LoadFromFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			m.History = []PeriodEntry{}
			return nil
		}
		return fmt.Errorf("failed to read file: %v", err)
	}
	return json.Unmarshal(data, &m.History)
}
