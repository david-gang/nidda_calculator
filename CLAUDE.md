# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run

```bash
go build -o nidda_calculator     # Build binary
go test ./pkg/nidda/...           # Run all tests
go test ./pkg/nidda/ -run TestName  # Run a single test
./nidda_calculator ui             # Launch TUI
./nidda_calculator add "1 Tishrei 5784" day   # Add period (Hebrew date)
./nidda_calculator add "2023-10-07" night      # Add period (Gregorian date)
```

## Architecture

Single Go module with CLI entry point (`main.go`) and core logic in `pkg/nidda/`.

### Core Package (`pkg/nidda/`)

- **types.go** — Domain types: `Onah` (Day/Night), `PeriodEntry`, `Prediction`, `NiddaManager` (central struct holding `History []PeriodEntry`)
- **predictions.go** — Halachic prediction calculations: Onah Beinonit (30th and 31st day), Haflagah (interval-based), Vesset HaChodesh (same Hebrew day next month). `GetAllPredictions()` is the main entry point — it checks cycle status first and returns either kavua-only or all three prediction types.
- **kavua.go** — Fixed cycle (וסת קבוע) detection and prediction. Checks 5 kavua types in priority order: Haflagah > Chodesh > Shavua > DilugChodesh > DilugHaflagah. Each requires 3 consecutive matching occurrences (same onah required).
- **parsing.go** — Date parsing supporting both Hebrew ("1 Tishrei 5784") and Gregorian ("YYYY-MM-DD") formats. Night onahs on Gregorian dates shift +1 Hebrew day (Hebrew day starts at sunset).
- **storage.go** — JSON serialization to `nidda_history.json`
- **calendar_ui.go** — Bubbletea TUI with Hebrew/Gregorian/Dual calendar views

### Key Dependencies

- `github.com/hebcal/hdate` — Hebrew calendar dates, month arithmetic, conversions. Uses `hdate.HDate` throughout; `Abs()` returns the Rata Die (absolute day number) used for date arithmetic.
- `github.com/charmbracelet/bubbletea` + `lipgloss` — Terminal UI framework and styling

## Domain Concepts

- **Onah**: Day or Night portion of a Hebrew date. Night precedes Day in the same Hebrew date.
- **Ohr Zarua**: The onah immediately preceding a predicted concern date. For Day predictions, it's the Night of the same date. For Night predictions, it's the Day of the previous date.
- **Kavua (fixed cycle)**: Established after 3 consecutive occurrences of the same pattern. When kavua is established, only the kavua prediction is returned (not the standard three).
- Hebrew month transitions: year increments at Elul→Tishrei boundary, with special handling for leap years (Adar I/II) via `hdate.MonthsInYear()`.
