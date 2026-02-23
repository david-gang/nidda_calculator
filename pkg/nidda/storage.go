package nidda

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hebcal/hdate"
)

// periodEntryJSON is used for custom marshal/unmarshal. We store "Adar I" and
// "Adar II" (hebcal's expected format) so MonthFromName parses correctly.
type periodEntryJSON struct {
	Date struct {
		Year  int    `json:"hy"`
		Month string `json:"hm"`
		Day   int    `json:"hd"`
	} `json:"Date"`
	Onah Onah `json:"Onah"`
}

func (pej *periodEntryJSON) toPeriodEntry() (PeriodEntry, error) {
	month, err := hdate.MonthFromName(pej.Date.Month)
	if err != nil {
		return PeriodEntry{}, err
	}
	hd := hdate.New(pej.Date.Year, month, pej.Date.Day)
	return PeriodEntry{Date: hd, Onah: pej.Onah}, nil
}

// monthForStorage returns the month string for JSON storage. Uses "Adar I" and
// "Adar II" (hebcal format) so round-trip works.
func monthForStorage(hd hdate.HDate) string {
	m := hd.Month()
	if m == hdate.Adar1 {
		return "Adar I"
	}
	if m == hdate.Adar2 {
		return "Adar II"
	}
	return m.String()
}

func (e PeriodEntry) toPeriodEntryJSON() periodEntryJSON {
	var pej periodEntryJSON
	pej.Date.Year = e.Date.Year()
	pej.Date.Month = monthForStorage(e.Date)
	pej.Date.Day = e.Date.Day()
	pej.Onah = e.Onah
	return pej
}

// SaveToFile saves the history to a JSON file.
func (m *NiddaManager) SaveToFile(filename string) error {
	raw := make([]periodEntryJSON, len(m.History))
	for i, e := range m.History {
		raw[i] = e.toPeriodEntryJSON()
	}
	data, err := json.MarshalIndent(raw, "", "  ")
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
	var raw []periodEntryJSON
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	m.History = make([]PeriodEntry, 0, len(raw))
	for _, pej := range raw {
		entry, err := pej.toPeriodEntry()
		if err != nil {
			return fmt.Errorf("failed to parse entry: %v", err)
		}
		m.History = append(m.History, entry)
	}
	return nil
}
