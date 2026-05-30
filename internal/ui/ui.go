package ui

import (
	"fmt"
	"path/filepath"
	"strings"

	"pocket-todo/internal/scanner"
	"pocket-todo/internal/storage"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type focusArea int

const (
	focusProject focusArea = iota
	focusFile
	focusLine
)

type model struct {
	// Data
	projects []storage.Project
	files    []storage.File
	todos    []storage.Todo

	// Cursors
	cursorProject int
	cursorFile    int
	cursorLine    int

	// Focus
	focus focusArea

	// Screen dimensions
	width  int
	height int
}

func initialModel() model {
	return model{
		projects:      loadProjects(),
		cursorProject: 0,
		cursorFile:    0,
		cursorLine:    0,
		focus:         focusProject,
	}
}

func loadProjects() []storage.Project {
	data_file_names, _ := scanner.ScanDataDir()
	var projects []storage.Project

	for _, name := range data_file_names {
		func() {
			defer func() {
				if r := recover(); r != nil {
					// Skip projects with JSON parsing errors
				}
			}()
			project := storage.ReadProjectFile(name)
			projects = append(projects, project)
		}()
	}

	return projects
}

// BACKEND TODO: Implement dynamic file loading when project changes
func (m *model) loadFilesForProject() {
	if len(m.projects) == 0 {
		m.files = []storage.File{}
		return
	}
	m.files = m.projects[m.cursorProject].Files
	m.cursorFile = 0
	m.loadTodosForFile()
}

// BACKEND TODO: Implement dynamic todos loading when file changes
func (m *model) loadTodosForFile() {
	if len(m.files) == 0 {
		m.todos = []storage.Todo{}
		return
	}
	m.todos = m.files[m.cursorFile].Todos
	m.cursorLine = 0
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "esc":
			if m.focus == focusProject {
				return m, tea.Quit
			}
			m.focus = focusProject

		case "up":
			m.moveUp()

		case "down":
			m.moveDown()

		case "enter", "space":
			m.selectItem()

		case "left":
			m.moveFocusLeft()

		case "right":
			m.moveFocusRight()

		// BACKEND TODO: Implement add/delete functionality
		case "a":
			m.addNewProject()

		case "d":
			m.deleteCurrentProject()

		case "x":
			// m.gotoCurrentProject()
		}
	}
	return m, nil
}

/* TODO: Move function with controller to different file */

func (m *model) addNewProject() {
	wd_path, err := filepath.Abs(".")
	if err != nil {
		fmt.Println("Error - Someting wrong with paths")
	}

	project_name := filepath.Base(wd_path)
	storage.CreateNewProjectFile(project_name, wd_path)

	project := storage.ReadProjectFile(project_name)

	new_data, err := scanner.ScanProject(project.Project_path)
	storage.SaveProjectFileData(project.Name, new_data)

	m.projects = loadProjects()
	m.cursorProject = 0
}

func (m *model) deleteCurrentProject() {
	storage.DeleteProjectFile(m.projects[m.cursorProject])

	m.projects = loadProjects()
	m.cursorProject = 0
}

func (m *model) moveUp() {
	switch m.focus {
	case focusProject:
		if m.cursorProject > 0 {
			m.cursorProject--
		} else {
			m.cursorProject = len(m.projects) - 1
		}
		m.loadFilesForProject()

	case focusFile:
		if m.cursorFile > 0 {
			m.cursorFile--
		} else {
			m.cursorFile = len(m.files) - 1
		}
		m.loadTodosForFile()

	case focusLine:
		if m.cursorLine > 0 {
			m.cursorLine--
		}
	}
}

func (m *model) moveDown() {
	switch m.focus {
	case focusProject:
		if m.cursorProject < len(m.projects)-1 {
			m.cursorProject++
		} else {
			m.cursorProject = 0
		}
		m.loadFilesForProject()

	case focusFile:
		if m.cursorFile < len(m.files)-1 {
			m.cursorFile++
		} else {
			m.cursorFile = 0
		}
		m.loadTodosForFile()

	case focusLine:
		if m.cursorLine < len(m.todos)-1 {
			m.cursorLine++
		}
	}
}

func (m *model) selectItem() {
	switch m.focus {
	case focusProject:
		m.focus = focusFile

	case focusFile:
		m.focus = focusLine
	}
}

func (m *model) moveFocusLeft() {
	if m.focus > focusProject {
		m.focus--
	}
}

func (m *model) moveFocusRight() {
	if m.focus < focusLine {
		m.focus++
	}
}

func (m *model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	// Calculate pane heights: 40% for columns, 60% for content
	columnHeight := int(float64(m.height) * 0.4)
	contentHeight := m.height - columnHeight - 2 // -2 for separators

	if columnHeight < 3 {
		columnHeight = 3
	}
	if contentHeight < 3 {
		contentHeight = 3
	}

	columnWidth := m.width / 3

	// Top section: 3 columns (Project | File | Line) - 40%
	topSection := m.renderTopSection(columnHeight, columnWidth)

	// Separator lines
	separator := lipgloss.NewStyle().
		Foreground(lipgloss.Color("239")).
		Render(strings.Repeat("─", m.width))

	// Bottom section: TODO content view - 60%
	bottomSection := m.renderBottomSection(contentHeight)

	return topSection + "\n" + separator + "\n" + bottomSection
}

func (m *model) renderTopSection(height int, columnWidth int) string {
	projectCol := m.renderProjectColumn(height-1, columnWidth)
	fileCol := m.renderFileColumn(height-1, columnWidth)
	lineCol := m.renderLineColumn(height-1, columnWidth)

	// Combine columns side by side
	projectLines := strings.Split(projectCol, "\n")
	fileLines := strings.Split(fileCol, "\n")
	lineLines := strings.Split(lineCol, "\n")

	displayHeight := height - 1

	// Pad all columns to the same height
	for len(projectLines) < displayHeight {
		projectLines = append(projectLines, "")
	}
	for len(fileLines) < displayHeight {
		fileLines = append(fileLines, "")
	}
	for len(lineLines) < displayHeight {
		lineLines = append(lineLines, "")
	}

	// Trim to exact height to avoid overflow
	projectLines = projectLines[:displayHeight]
	fileLines = fileLines[:displayHeight]
	lineLines = lineLines[:displayHeight]

	result := ""
	sep := "│"

	for i := 0; i < displayHeight; i++ {
		// Ensure each line is exactly columnWidth
		p := m.padString(projectLines[i], columnWidth)
		f := m.padString(fileLines[i], columnWidth)
		l := m.padString(lineLines[i], columnWidth)

		line := p + " " + sep + " " + f + " " + sep + " " + l
		result += line + "\n"
	}

	// Add controls at the bottom
	controlsStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("239")).
		Faint(true)
	controls := controlsStyle.Render("[q]uit [a]dd [x]open [d]el")
	controlLine := m.padString(controls, m.width)
	result += controlLine

	return result
}

func (m *model) renderProjectColumn(height int, width int) string {
	count := len(m.projects)
	title := m.renderColumnTitle("Project", focusProject, width, count)

	items := ""
	for i, project := range m.projects {
		line := m.renderItem(project.Name, i == m.cursorProject && m.focus == focusProject, width-2)
		items += line + "\n"
	}

	// Calculate viewport for scrolling
	itemLines := strings.Split(strings.TrimSuffix(items, "\n"), "\n")
	visibleLines := height - 1

	// Calculate scroll position to keep cursor visible
	startIdx := 0
	if m.cursorProject >= visibleLines {
		startIdx = m.cursorProject - visibleLines + 1
	}

	if len(itemLines) > visibleLines {
		if startIdx+visibleLines <= len(itemLines) {
			itemLines = itemLines[startIdx : startIdx+visibleLines]
		} else {
			itemLines = itemLines[len(itemLines)-visibleLines:]
		}
	}

	return title + "\n" + strings.Join(itemLines, "\n")
}

func (m *model) renderFileColumn(height int, width int) string {
	count := len(m.files)
	title := m.renderColumnTitle("File", focusFile, width, count)

	items := ""
	for i, file := range m.files {
		line := m.renderItem(file.Name, i == m.cursorFile && m.focus == focusFile, width-2)
		items += line + "\n"
	}

	// Calculate viewport for scrolling
	itemLines := strings.Split(strings.TrimSuffix(items, "\n"), "\n")
	visibleLines := height - 1

	// Calculate scroll position to keep cursor visible
	startIdx := 0
	if m.cursorFile >= visibleLines {
		startIdx = m.cursorFile - visibleLines + 1
	}

	if len(itemLines) > visibleLines {
		if startIdx+visibleLines <= len(itemLines) {
			itemLines = itemLines[startIdx : startIdx+visibleLines]
		} else {
			itemLines = itemLines[len(itemLines)-visibleLines:]
		}
	}

	return title + "\n" + strings.Join(itemLines, "\n")
}

func (m *model) renderLineColumn(height int, width int) string {
	count := len(m.todos)
	title := m.renderColumnTitle("Line", focusLine, width, count)

	items := ""
	for i, todo := range m.todos {
		lineStr := fmt.Sprintf("%d", todo.Line)
		line := m.renderItem(lineStr, i == m.cursorLine && m.focus == focusLine, width-2)
		items += line + "\n"
	}

	// Calculate viewport for scrolling
	itemLines := strings.Split(strings.TrimSuffix(items, "\n"), "\n")
	visibleLines := height - 1

	// Calculate scroll position to keep cursor visible
	startIdx := 0
	if m.cursorLine >= visibleLines {
		startIdx = m.cursorLine - visibleLines + 1
	}

	if len(itemLines) > visibleLines {
		if startIdx+visibleLines <= len(itemLines) {
			itemLines = itemLines[startIdx : startIdx+visibleLines]
		} else {
			itemLines = itemLines[len(itemLines)-visibleLines:]
		}
	}

	return title + "\n" + strings.Join(itemLines, "\n")
}

func (m *model) renderColumnTitle(title string, focus focusArea, width int, count int) string {
	countStr := fmt.Sprintf(" (%d)", count)
	fullTitle := title + countStr

	if m.focus == focus {
		return lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("51")).
			Background(lipgloss.Color("243")).
			Render("► " + fullTitle)
	}
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("117")).
		Render("  " + fullTitle)
}

func (m *model) renderItem(text string, selected bool, width int) string {
	if selected {
		return lipgloss.NewStyle().
			Background(lipgloss.Color("237")).
			Foreground(lipgloss.Color("15")).
			Bold(true).
			Render("> " + m.truncateString(text, width-2))
	}
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("247")).
		Render("  " + m.truncateString(text, width-2))
}

func (m *model) renderBottomSection(height int) string {
	if len(m.todos) == 0 {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("242")).
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("238")).
			Padding(1).
			Render("✨ No TODOs in this file\n\nNavigate to a file with content to see details here")
	}

	currentTodo := m.todos[m.cursorLine]

	// Format the header with colors
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("18")).
		Bold(true)
	header := headerStyle.Render(fmt.Sprintf("📝 Line %d", currentTodo.Line))

	// Wrap and truncate content to fit in box
	contentLines := strings.Split(currentTodo.Content, "\n")
	maxLines := height - 4
	if len(contentLines) > maxLines {
		contentLines = contentLines[:maxLines]
		contentLines = append(contentLines, "...")
	}

	content := strings.Join(contentLines, "\n")

	// Create styled box
	display := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("33")).
		Foreground(lipgloss.Color("251")).
		Padding(1).
		Render(header + "\n\n" + content)

	return display
}

func (m *model) renderFooter() string {
	keys := lipgloss.NewStyle().
		Foreground(lipgloss.Color("117")).
		Render("q-quit  a-add  x-open  d-delete")

	return keys
}

func (m *model) padString(s string, width int) string {
	// Count only visible characters (not ANSI codes)
	visibleLen := m.visibleLength(s)
	if visibleLen >= width {
		return s
	}
	return s + strings.Repeat(" ", width-visibleLen)
}

func (m *model) visibleLength(s string) int {
	// Remove ANSI escape sequences and count visible chars
	visibleLen := 0
	inEscape := false
	for _, r := range s {
		if r == '\x1b' { // ESC character
			inEscape = true
		} else if inEscape && (r == 'm' || r == 'K') { // End of escape sequence
			inEscape = false
		} else if !inEscape {
			visibleLen++
		}
	}
	return visibleLen
}

func (m *model) truncateString(s string, maxLen int) string {
	visibleLen := m.visibleLength(s)
	if visibleLen > maxLen {
		// Truncate at maxLen-1 to account for ellipsis
		return s[:len(s)-5] + "…"
	}
	return s
}

func Run() {
	m := initialModel()
	p := tea.NewProgram(&m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
	}
}
