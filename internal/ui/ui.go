package ui

import (
	"fmt"

	"pocket-todo/internal/scanner"
	"pocket-todo/internal/storage"

	tea "github.com/charmbracelet/bubbletea"
)

type screen int

const (
	projectsScreen screen = iota
	filesScreen
	todosScreen
)

type model struct {
	screen screen

	projects []storage.Project
	files    []storage.File
	todos    []storage.Todo

	cursorProjects int
	cursorFiles    int
	cursorTodos    int
}

func initialModel() model {
	return model{
		projects: loadProjects(),
	}
}

func loadProjects() []storage.Project {
	data_file_names, _ := scanner.ScanDataDir() // TODO: Handle err correctly
	var project []storage.Project

	for _, name := range data_file_names {
		project = append(project, storage.ReadProjectFile(name))
	}

	return project
}

func loadFiles() {

}

func loadTodos() {

}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "up":
			m.moveUp()

		case "down":
			m.moveDown()

		case "enter", "space":
			m.selectItem()

		case "a":
			m.addNewProject()

		case "d":
			m.deleteCurrentProject()
		}
	}
	return m, nil
}

func (m model) moveUp() {
	switch m.screen {
	case projectsScreen:
		if m.cursorProjects > 0 {
			m.cursorProjects--
		} else {
			m.cursorProjects = len(m.projects) - 1
		}
	case filesScreen:
		if m.cursorFiles > 0 {
			m.cursorFiles--
		} else {
			m.cursorFiles = len(m.files) - 1
		}
	case todosScreen:
		if m.cursorTodos > 0 {
			m.cursorTodos--
		} else {
			m.cursorTodos = len(m.todos) - 1
		}
	}
}

func (m model) moveDown() {
	switch m.screen {
	case projectsScreen:
		if m.cursorProjects < len(m.projects)-1 {
			m.cursorProjects++
		} else {
			m.cursorProjects = 0
		}
	case filesScreen:
		if m.cursorFiles < len(m.files)-1 {
			m.cursorFiles++
		} else {
			m.cursorFiles = 0
		}
	case todosScreen:
		if m.cursorTodos < len(m.todos) {
			m.cursorTodos++
		} else {
			m.cursorTodos = 0
		}
	}
}

func (m model) selectItem() {
	switch m.screen {
	case projectsScreen:
		// m.files = loadFiles(m.projects[m.cursorProjects])
		m.screen = filesScreen

	case filesScreen:
		// m.todos = loadTodos(m.files[m.cursorFiles])
		m.screen = todosScreen
	}
}

func (m model) View() string {
	switch m.screen {
	case projectsScreen:
		return m.renderProjects()

	case filesScreen:
		return m.renderFiles()

	case todosScreen:
		return m.renderTodos()
	}

	return "unknown screen"
}

func (m model) renderProjects() string {
	s := "PROJECTS:\n\n"

	for i, p := range m.projects {
		cursor := " "
		if i == m.cursorProjects {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, p.Name)
	}

	return s
}

func (m model) renderFiles() string {
	s := "FILES:\n\n"

	for i, f := range m.files {
		cursor := " "
		if i == m.cursorFiles {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, f.Name)
	}

	return s
}

func (m model) renderTodos() string {
	s := "TODO:\n\n"

	for i, t := range m.todos {
		cursor := " "
		if i == m.cursorTodos {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, t.Content)
	}

	return s
}
