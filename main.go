package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"nidda_calculator/pkg/nidda"
)

const storageFile = "nidda_history.json"

func main() {
	manager := &nidda.NiddaManager{}
	if err := manager.LoadFromFile(storageFile); err != nil {
		log.Fatalf("Error loading history: %v", err)
	}

	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]

	switch command {
	case "add":
		handleAdd(manager)
	case "rm":
		handleRemove(manager)
	case "list":
		handleList(manager)
	case "predict":
		handlePredict(manager)
	case "ui":
		handleUI(manager)
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Nidda Calendar Manager CLI")
	fmt.Println("Usage:")
	fmt.Println("  ui                        - Launch interactive calendar UI")
	fmt.Println("  add \"<date>\" <day|night>  - Add a new period. Date can be Hebrew or Gregorian.")
	fmt.Println("                               (e.g., add \"1 Tishrei 5784\" day OR add \"2023-10-07\" day)")
	fmt.Println("  rm <index>                - Remove a period by its list index (e.g., rm 1)")
	fmt.Println("  list                      - List all recorded periods")
	fmt.Println("  predict                   - Show upcoming concern dates in chronological order")
	fmt.Println("  help                      - Show this help message")
}

func handleAdd(m *nidda.NiddaManager) {
	if len(os.Args) < 4 {
		fmt.Println("Usage: add \"<date>\" <day|night>")
		return
	}

	dateStr := os.Args[2]
	onahStr := os.Args[3]

	onah, err := nidda.ParseOnah(onahStr)
	if err != nil {
		fmt.Printf("Error parsing onah: %v\n", err)
		return
	}

	hDate, err := nidda.ParseDate(dateStr, onah)
	if err != nil {
		fmt.Printf("Error parsing date: %v\n", err)
		return
	}

	// Basic validation: No future dates
	if nidda.IsPeriodInFuture(hDate, onah, time.Now()) {
		fmt.Println("Error: Cannot add a period in the future.")
		return
	}

	m.AddPeriod(hDate, onah)
	if err := m.SaveToFile(storageFile); err != nil {
		fmt.Printf("Error saving history: %v\n", err)
		return
	}

	fmt.Printf("Successfully added period: %s (%s)\n", hDate.String(), onah)
}

func handleRemove(m *nidda.NiddaManager) {
	if len(os.Args) < 3 {
		fmt.Println("Usage: rm <index>")
		return
	}

	indexStr := os.Args[2]
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		fmt.Printf("Invalid index: %s\n", indexStr)
		return
	}

	// Convert 1-based index to 0-based
	if err := m.RemovePeriod(index - 1); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if err := m.SaveToFile(storageFile); err != nil {
		fmt.Printf("Error saving history: %v\n", err)
		return
	}

	fmt.Printf("Successfully removed period #%d\n", index)
}

func handleList(m *nidda.NiddaManager) {
	if len(m.History) == 0 {
		fmt.Println("No periods recorded.")
		return
	}

	fmt.Println("Recorded Periods:")
	for i, entry := range m.History {
		fmt.Printf("%d. %s (%s)\n", i+1, entry.Date.String(), entry.Onah)
	}
}

func handlePredict(m *nidda.NiddaManager) {
	if len(m.History) == 0 {
		fmt.Println("No periods recorded yet. Add one first.")
		return
	}

	predictions := m.GetAllPredictions()
	if len(predictions) == 0 {
		fmt.Println("No upcoming concern dates.")
		return
	}

	fmt.Println("\nUpcoming Halachic Concern Dates (Timeline):")
	for _, p := range predictions {
		printPrediction(p)
	}
}

func printPrediction(p nidda.Prediction) {
	year, month, day := p.Date.Greg()
	gregDate := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	gregStr := gregDate.Format("2006-01-02 (Mon)")
	fmt.Printf("\n* %s: %s (%s) [%s]\n", p.Type, p.Date.String(), p.Onah, gregStr)
	if p.IncludesOhrZarua {
		oz := p.OhrZarua()
		ozYear, ozMonth, ozDay := oz.Date.Greg()
		ozGregDate := time.Date(ozYear, ozMonth, ozDay, 0, 0, 0, 0, time.UTC)
		ozGregStr := ozGregDate.Format("2006-01-02 (Mon)")
		fmt.Printf("  -> %s: %s (%s) [%s]\n", oz.Type, oz.Date.String(), oz.Onah, ozGregStr)
	}
}

func handleUI(m *nidda.NiddaManager) {
	if err := nidda.RunCalendarUI(m); err != nil {
		log.Fatalf("Error running UI: %v", err)
	}
}
