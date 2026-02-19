package nidda

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hebcal/hdate"
)

// ViewMode represents the calendar display mode
type ViewMode int

const (
	ViewHebrew ViewMode = iota
	ViewGregorian
	ViewDual
)

func (v ViewMode) String() string {
	switch v {
	case ViewHebrew:
		return "Hebrew"
	case ViewGregorian:
		return "Gregorian"
	case ViewDual:
		return "Dual"
	default:
		return "Hebrew"
	}
}

// Key bindings for the calendar UI
type keyMap struct {
	Up        key.Binding
	Down      key.Binding
	Left      key.Binding
	Right     key.Binding
	NextMonth key.Binding
	PrevMonth key.Binding
	Today     key.Binding
	Toggle    key.Binding
	Help      key.Binding
	Quit      key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.PrevMonth, k.NextMonth, k.Today},
		{k.Toggle, k.Help, k.Quit},
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "right"),
	),
	PrevMonth: key.NewBinding(
		key.WithKeys("pgup", "["),
		key.WithHelp("pgup/[", "prev month"),
	),
	NextMonth: key.NewBinding(
		key.WithKeys("pgdown", "]"),
		key.WithHelp("pgdn/]", "next month"),
	),
	Today: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "go to today"),
	),
	Toggle: key.NewBinding(
		key.WithKeys("t"),
		key.WithHelp("t", "toggle view"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q/esc", "quit"),
	),
}

// Styles for the calendar
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			Padding(0, 1)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("99")).
			Align(lipgloss.Center).
			Width(4)

	dayStyle = lipgloss.NewStyle().
			Width(4).
			Align(lipgloss.Center)

	selectedStyle = lipgloss.NewStyle().
			Width(4).
			Align(lipgloss.Center).
			Foreground(lipgloss.Color("0")).
			Background(lipgloss.Color("45")).
			Bold(true)

	periodStyle = lipgloss.NewStyle().
			Width(4).
			Align(lipgloss.Center).
			Foreground(lipgloss.Color("51")).
			Bold(true)

	predictionStyle = lipgloss.NewStyle().
			Width(4).
			Align(lipgloss.Center).
			Foreground(lipgloss.Color("226")).
			Bold(true)

	ohrZaruaStyle = lipgloss.NewStyle().
			Width(4).
			Align(lipgloss.Center).
			Foreground(lipgloss.Color("214")).
			Bold(true)

	todayStyle = lipgloss.NewStyle().
			Width(4).
			Align(lipgloss.Center).
			Foreground(lipgloss.Color("46")).
			Bold(true)

	detailsStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2).
			MarginTop(1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginTop(1)
)

// DateInfo holds information about a specific date
type DateInfo struct {
	IsPeriod     bool
	PeriodOnah   Onah
	IsPrediction bool
	Predictions  []Prediction
	IsOhrZarua   bool
	OhrZaruaFor  []Prediction
}

// CalendarModel represents the state of the calendar UI
type CalendarModel struct {
	manager      *NiddaManager
	currentMonth int
	currentYear  int
	selectedDate hdate.HDate
	dateInfo     map[int64]DateInfo // keyed by Abs() day
	viewMode     ViewMode
	help         help.Model
	showHelp     bool
	width        int
	height       int
}

// NewCalendarModel creates a new calendar model
func NewCalendarModel(manager *NiddaManager) CalendarModel {
	now := hdate.FromTime(time.Now())
	m := CalendarModel{
		manager:      manager,
		currentMonth: int(now.Month()),
		currentYear:  now.Year(),
		selectedDate: now,
		dateInfo:     make(map[int64]DateInfo),
		viewMode:     ViewHebrew,
		help:         help.New(),
		showHelp:     false,
	}
	m.buildDateInfo()
	return m
}

// buildDateInfo populates the dateInfo map with periods and predictions
func (m *CalendarModel) buildDateInfo() {
	m.dateInfo = make(map[int64]DateInfo)

	// Add period dates
	for _, entry := range m.manager.History {
		abs := entry.Date.Abs()
		info := m.dateInfo[abs]
		info.IsPeriod = true
		info.PeriodOnah = entry.Onah
		m.dateInfo[abs] = info
	}

	// Add predictions
	predictions := m.manager.GetAllPredictions()
	for _, pred := range predictions {
		abs := pred.Date.Abs()
		info := m.dateInfo[abs]
		info.IsPrediction = true
		info.Predictions = append(info.Predictions, pred)
		m.dateInfo[abs] = info

		// Add Ohr Zarua dates
		oz := pred.OhrZarua()
		ozAbs := oz.Date.Abs()
		ozInfo := m.dateInfo[ozAbs]
		ozInfo.IsOhrZarua = true
		ozInfo.OhrZaruaFor = append(ozInfo.OhrZaruaFor, pred)
		m.dateInfo[ozAbs] = ozInfo
	}
}

func (m CalendarModel) Init() tea.Cmd {
	return nil
}

func (m CalendarModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, keys.Help):
			m.showHelp = !m.showHelp
			return m, nil

		case key.Matches(msg, keys.Toggle):
			// Cycle through view modes: Hebrew -> Gregorian -> Dual -> Hebrew
			m.viewMode = (m.viewMode + 1) % 3
			return m, nil

		case key.Matches(msg, keys.Today):
			now := hdate.FromTime(time.Now())
			m.selectedDate = now
			m.currentMonth = int(now.Month())
			m.currentYear = now.Year()
			return m, nil

		case key.Matches(msg, keys.NextMonth):
			m.nextMonth()
			return m, nil

		case key.Matches(msg, keys.PrevMonth):
			m.prevMonth()
			return m, nil

		case key.Matches(msg, keys.Left):
			m.selectedDate = hdate.FromRD(m.selectedDate.Abs() - 1)
			m.ensureSelectedInView()
			return m, nil

		case key.Matches(msg, keys.Right):
			m.selectedDate = hdate.FromRD(m.selectedDate.Abs() + 1)
			m.ensureSelectedInView()
			return m, nil

		case key.Matches(msg, keys.Up):
			m.selectedDate = hdate.FromRD(m.selectedDate.Abs() - 7)
			m.ensureSelectedInView()
			return m, nil

		case key.Matches(msg, keys.Down):
			m.selectedDate = hdate.FromRD(m.selectedDate.Abs() + 7)
			m.ensureSelectedInView()
			return m, nil
		}
	}

	return m, nil
}

// ensureSelectedInView updates the current month/year if selected date is out of view
func (m *CalendarModel) ensureSelectedInView() {
	if int(m.selectedDate.Month()) != m.currentMonth || m.selectedDate.Year() != m.currentYear {
		m.currentMonth = int(m.selectedDate.Month())
		m.currentYear = m.selectedDate.Year()
	}
}

// nextMonth advances to the next Hebrew month
func (m *CalendarModel) nextMonth() {
	monthsInYear := hdate.MonthsInYear(m.currentYear)
	if m.currentMonth < monthsInYear {
		m.currentMonth++
	} else {
		m.currentMonth = 1
		m.currentYear++
	}
}

// prevMonth goes back to the previous Hebrew month
func (m *CalendarModel) prevMonth() {
	if m.currentMonth > 1 {
		m.currentMonth--
	} else {
		m.currentYear--
		m.currentMonth = hdate.MonthsInYear(m.currentYear)
	}
}

func (m CalendarModel) View() string {
	var b strings.Builder

	// Title with view mode indicator
	var title string
	switch m.viewMode {
	case ViewHebrew:
		monthName := hdate.HMonth(m.currentMonth).String()
		title = fmt.Sprintf("Nidda Calendar - %s %d [Hebrew]", monthName, m.currentYear)
	case ViewGregorian:
		gregDate := hdate.New(m.currentYear, hdate.HMonth(m.currentMonth), 1).Gregorian()
		title = fmt.Sprintf("Nidda Calendar - %s %d [Gregorian]", gregDate.Format("January"), gregDate.Year())
	case ViewDual:
		monthName := hdate.HMonth(m.currentMonth).String()
		// In dual view, show the Gregorian month containing the selected date
		gregDate := m.selectedDate.Gregorian()
		title = fmt.Sprintf("Nidda Calendar - %s %d / %s %d [Dual]",
			monthName, m.currentYear, gregDate.Format("January"), gregDate.Year())
	}
	b.WriteString(titleStyle.Render(title))
	b.WriteString("\n\n")

	// Calendar grid
	b.WriteString(m.renderCalendar())
	b.WriteString("\n")

	// Details panel
	b.WriteString(m.renderDetails())
	b.WriteString("\n")

	// Help
	if m.showHelp {
		b.WriteString(m.help.View(keys))
	} else {
		b.WriteString(helpStyle.Render("Press t to toggle view • ? for help • q to quit"))
	}

	return b.String()
}

// renderCalendar renders the calendar grid
func (m CalendarModel) renderCalendar() string {
	switch m.viewMode {
	case ViewHebrew:
		return m.renderHebrewCalendar(true)
	case ViewGregorian:
		return m.renderGregorianCalendar(true)
	case ViewDual:
		return m.renderDualCalendar()
	default:
		return m.renderHebrewCalendar(true)
	}
}

// renderHebrewCalendar renders a Hebrew calendar month
func (m CalendarModel) renderHebrewCalendar(showSelection bool) string {
	var b strings.Builder

	// Day headers
	headers := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	for _, h := range headers {
		b.WriteString(headerStyle.Render(h))
		b.WriteString(" ")
	}
	b.WriteString("\n")

	// Get the first day of the month
	firstDay := hdate.New(m.currentYear, hdate.HMonth(m.currentMonth), 1)
	daysInMonth := hdate.DaysInMonth(hdate.HMonth(m.currentMonth), m.currentYear)

	// Find what day of the week the month starts on (0 = Sunday)
	gregTime := firstDay.Gregorian()
	firstDayOfWeek := int(gregTime.Weekday())

	// Render empty cells before the first day
	for i := 0; i < firstDayOfWeek; i++ {
		b.WriteString(dayStyle.Render(""))
		b.WriteString(" ")
	}

	// Render each day of the month
	today := hdate.FromTime(time.Now())
	for day := 1; day <= daysInMonth; day++ {
		date := hdate.New(m.currentYear, hdate.HMonth(m.currentMonth), day)
		dayStr := fmt.Sprintf("%d", day)

		isSelected := showSelection && date.Abs() == m.selectedDate.Abs()
		isToday := date.Abs() == today.Abs()
		info := m.dateInfo[date.Abs()]

		style := m.getStyleForDate(isSelected, isToday, info)
		b.WriteString(style.Render(dayStr))
		b.WriteString(" ")

		// New line after Saturday
		if (firstDayOfWeek+day)%7 == 0 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

// renderGregorianCalendar renders a Gregorian calendar month
func (m CalendarModel) renderGregorianCalendar(showSelection bool) string {
	return m.renderGregorianCalendarForDate(m.selectedDate, showSelection)
}

// renderGregorianCalendarForDate renders a Gregorian month containing the specified date
func (m CalendarModel) renderGregorianCalendarForDate(refDate hdate.HDate, showSelection bool) string {
	var b strings.Builder

	// Day headers
	headers := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	for _, h := range headers {
		b.WriteString(headerStyle.Render(h))
		b.WriteString(" ")
	}
	b.WriteString("\n")

	// In dual view, show the Gregorian month containing the selected date
	// Otherwise, show the month that overlaps with the Hebrew month start
	gregRef := refDate.Gregorian()
	gregYear := gregRef.Year()
	gregMonth := gregRef.Month()

	gregMonthFirst := time.Date(gregYear, gregMonth, 1, 0, 0, 0, 0, time.UTC)
	daysInGregMonth := time.Date(gregYear, gregMonth+1, 0, 0, 0, 0, 0, time.UTC).Day()

	firstDayOfWeek := int(gregMonthFirst.Weekday())

	// Render empty cells before the first day
	for i := 0; i < firstDayOfWeek; i++ {
		b.WriteString(dayStyle.Render(""))
		b.WriteString(" ")
	}

	// Render each day of the Gregorian month
	today := hdate.FromTime(time.Now())
	for day := 1; day <= daysInGregMonth; day++ {
		gregDate := time.Date(gregYear, gregMonth, day, 0, 0, 0, 0, time.UTC)
		hebrewDate := hdate.FromTime(gregDate)
		dayStr := fmt.Sprintf("%d", day)

		isSelected := showSelection && hebrewDate.Abs() == m.selectedDate.Abs()
		isToday := hebrewDate.Abs() == today.Abs()
		info := m.dateInfo[hebrewDate.Abs()]

		style := m.getStyleForDate(isSelected, isToday, info)
		b.WriteString(style.Render(dayStr))
		b.WriteString(" ")

		// New line after Saturday
		if (firstDayOfWeek+day)%7 == 0 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

// renderDualCalendar renders Hebrew and Gregorian calendars side by side
func (m CalendarModel) renderDualCalendar() string {
	var b strings.Builder

	// In dual view, show selection on BOTH calendars (same date in both systems)
	hebrewCal := m.renderHebrewCalendar(true)
	gregorianCal := m.renderGregorianCalendar(true)

	// Split into lines and ensure consistent handling
	hebrewLines := strings.Split(strings.TrimRight(hebrewCal, "\n"), "\n")
	gregorianLines := strings.Split(strings.TrimRight(gregorianCal, "\n"), "\n")

	// Side-by-side labels (Gregorian on left, Hebrew on right)
	gregLabel := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("99")).Render("Gregorian Calendar")
	hebrewLabel := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("99")).Render("Hebrew Calendar")

	b.WriteString(gregLabel)
	b.WriteString(strings.Repeat(" ", 13))
	b.WriteString(hebrewLabel)
	b.WriteString("\n")

	// Combine both calendars side by side (Gregorian left, Hebrew right)
	maxLines := len(hebrewLines)
	if len(gregorianLines) > maxLines {
		maxLines = len(gregorianLines)
	}

	for i := 0; i < maxLines; i++ {
		// Gregorian calendar line
		if i < len(gregorianLines) {
			gregLine := gregorianLines[i]
			b.WriteString(gregLine)
			// Calculate visible width and pad accordingly
			visibleWidth := lipgloss.Width(gregLine)
			if visibleWidth < 35 {
				b.WriteString(strings.Repeat(" ", 35-visibleWidth))
			}
		} else {
			b.WriteString(strings.Repeat(" ", 35)) // padding if Gregorian calendar is shorter
		}

		b.WriteString("   ") // spacing between calendars

		// Hebrew calendar line
		if i < len(hebrewLines) {
			b.WriteString(hebrewLines[i])
		}

		b.WriteString("\n")
	}

	return b.String()
}

// getStyleForDate returns the appropriate style for a date based on its properties
func (m CalendarModel) getStyleForDate(isSelected, isToday bool, info DateInfo) lipgloss.Style {
	if isSelected {
		return selectedStyle
	} else if info.IsPeriod {
		return periodStyle
	} else if info.IsOhrZarua {
		return ohrZaruaStyle
	} else if info.IsPrediction {
		return predictionStyle
	} else if isToday {
		return todayStyle
	}
	return dayStyle
}

// renderDetails renders the details panel for the selected date
func (m CalendarModel) renderDetails() string {
	var b strings.Builder

	// Date header
	gregTime := m.selectedDate.Gregorian()
	b.WriteString(fmt.Sprintf("Selected: %s (%s)\n\n",
		m.selectedDate.String(),
		gregTime.Format("Mon Jan 2, 2006")))

	info := m.dateInfo[m.selectedDate.Abs()]

	if !info.IsPeriod && !info.IsPrediction && !info.IsOhrZarua {
		b.WriteString("No events on this date.")
		return detailsStyle.Render(b.String())
	}

	// Show period information
	if info.IsPeriod {
		b.WriteString(fmt.Sprintf("🔴 Period: %s\n", info.PeriodOnah))
	}

	// Show predictions
	if info.IsPrediction {
		for _, pred := range info.Predictions {
			b.WriteString(fmt.Sprintf("⚠️  %s (%s)\n", pred.Type, pred.Onah))
		}
	}

	// Show Ohr Zarua
	if info.IsOhrZarua {
		for _, pred := range info.OhrZaruaFor {
			b.WriteString(fmt.Sprintf("🌙 Ohr Zarua for %s (%s)\n", pred.Type, info.OhrZaruaFor[0].Onah))
		}
	}

	return detailsStyle.Render(b.String())
}

// RunCalendarUI launches the calendar UI
func RunCalendarUI(manager *NiddaManager) error {
	p := tea.NewProgram(NewCalendarModel(manager), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
