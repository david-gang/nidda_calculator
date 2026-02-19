# Nidda Calculator

A halachic calendar application for tracking menstrual periods and calculating important dates according to Jewish law (niddah/family purity laws).

> **Disclaimer:** This project is in beta. Predictions may contain errors. This tool is intended as a convenience aid and **does not substitute consulting a competent posek (halachic authority)**. Always verify results with your rabbi.

## Features

### Interactive Calendar UI
Launch a beautiful terminal-based calendar interface:
```bash
./nidda_calculator ui
```

The calendar displays:
- 🔵 **Blue**: Recorded period dates
- 🟡 **Yellow**: Predicted concern dates (Onah Beinonit, Haflaga, Vesset HaChodesh)
- 🟠 **Orange**: Ohr Zarua dates (preceding onah for predictions)
- 🟢 **Green**: Current date
- 🔵 **Cyan background**: Selected date cursor

#### Keyboard Controls
- **Arrow keys** or **h/j/k/l**: Navigate dates
- **PgUp/PgDn** or **[/]**: Switch months
- **t**: Toggle view mode (Hebrew → Gregorian → Dual → Hebrew)
- **n**: Jump to today
- **?**: Toggle help
- **q/Esc**: Exit to CLI

### View Modes

Press **t** to cycle through three view modes:

1. **Hebrew Calendar** (default): Shows the Hebrew month (Tishrei, Cheshvan, etc.)
2. **Gregorian Calendar**: Shows the corresponding Gregorian month
3. **Dual View**: Shows both calendars side-by-side for easy comparison

All period dates and predictions are displayed in all view modes, making it easy to see important dates in whichever calendar system you prefer.

### Command-Line Interface

#### Add a Period
```bash
# Using Hebrew date
./nidda_calculator add "1 Tishrei 5784" day

# Using Gregorian date
./nidda_calculator add "2023-10-07" night
```

#### List All Periods
```bash
./nidda_calculator list
```

#### Show Predictions
```bash
./nidda_calculator predict
```

#### Remove a Period
```bash
./nidda_calculator rm 1
```

## Halachic Calculations

The calculator computes three main concern dates:

1. **Onah Beinonit**: The 30th day from the last period
2. **Haflaga**: Based on the interval between the last two periods
3. **Vesset HaChodesh**: Same Hebrew day in the next Hebrew month

For each prediction, it also calculates the **Ohr Zarua** (the immediately preceding onah).

## Installation

### Prerequisites
- Go 1.24.3 or later

### Build
```bash
go build -o nidda_calculator
```

## Data Storage

Period history is stored in `nidda_history.json` in the current directory.

## Hebrew Calendar Support

The calculator fully supports the Hebrew calendar, including:
- Leap years (with Adar I and Adar II)
- Variable month lengths (29 or 30 days)
- Conversion between Hebrew and Gregorian dates
- Proper handling of night/day onah transitions

## Technologies

- **Go**: Main programming language
- **[hebcal/hdate](https://github.com/hebcal/hdate)**: Hebrew calendar library
- **[bubbletea](https://github.com/charmbracelet/bubbletea)**: Terminal UI framework
- **[lipgloss](https://github.com/charmbracelet/lipgloss)**: Terminal styling

## Contributing

Contributions are welcome! Here are some ways you can help:

- **Bug reports** — If you find a calculation error or unexpected behavior, please [open an issue](https://github.com/david-gang/nidda_calculator/issues). Halachic calculation bugs are especially important to catch.
- **Halachic review** — If you have halachic knowledge, we'd appreciate review of the prediction logic in `pkg/nidda/` for correctness.
- **Feature requests** — Ideas for new features (additional kavua patterns, support for different minhagim, etc.) are welcome as issues.
- **Code contributions** — Fork the repo, make your changes, and open a pull request. Please run `go test ./pkg/nidda/...` before submitting.

## License

This project is licensed under the [MIT License](LICENSE).

