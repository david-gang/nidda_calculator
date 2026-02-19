# Quick Start Guide

## Launch the Calendar UI

```bash
./nidda_calculator ui
```

## First Time Setup

1. **Build the program**:
   ```bash
   go build -o nidda_calculator
   ```

2. **Add your first period**:
   ```bash
   # Hebrew date
   ./nidda_calculator add "23 Tishrei 5784" night
   
   # OR Gregorian date
   ./nidda_calculator add "2023-10-09" night
   ```

3. **Launch the UI**:
   ```bash
   ./nidda_calculator ui
   ```

## Using the Calendar

### Navigate
- **↑ ↓ ← →**: Move between dates
- **[** or **PgUp**: Previous month
- **]** or **PgDn**: Next month
- **n**: Jump to today
- **t**: Toggle view (Hebrew/Gregorian/Dual)

### View Modes
Press **t** to cycle through:
1. **Hebrew Calendar**: Hebrew months (Tishrei, Cheshvan, etc.)
2. **Gregorian Calendar**: Standard calendar months
3. **Dual View**: Both calendars side-by-side

### View Information
- Select any date to see details in the panel below
- Colors indicate:
  - **Cyan/Blue**: Recorded period
  - **Yellow**: Predicted concern date
  - **Orange**: Ohr Zarua (preceding onah)
  - **Green**: Today's date

### Exit
- Press **q**, **Esc**, or **Ctrl+C**

## Command Line Tools

```bash
# View all periods
./nidda_calculator list

# See predictions
./nidda_calculator predict

# Remove a period (by number from list)
./nidda_calculator rm 1

# Show help
./nidda_calculator help
```

## Tips

- Press **t** to toggle between Hebrew, Gregorian, and dual calendar views
- The calendar displays Hebrew months by default (Tishrei, Cheshvan, etc.)
- In dual view mode, both calendars are shown side-by-side
- Dates are shown with both Hebrew and Gregorian equivalents in the details panel
- All data is saved automatically to `nidda_history.json`
- Press **?** in the UI for full keyboard shortcuts

