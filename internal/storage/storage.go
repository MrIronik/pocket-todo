package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Project struct {
	Name        string
	Total_todos int
	Files       File
}

type File struct {
	Name  string
	Todos []Todo
}

type Todo struct {
	Line    int
	Content string
}

/* Variables for Data Manager */
var data_folder_path string

func Init() error {
	err := formatDataFolderPath()
	if err != nil {
		fmt.Println("Init Error - Getting Path")
		return err
	}

	err = os.MkdirAll(data_folder_path, 0755)
	if err != nil {
		return err
	}

	return nil
}

func CreateNewProjectFile(project_name string) error {
	file_name := project_name + "_project.json"
	new_project_path := filepath.Join(data_folder_path, file_name)

	_, err := os.Stat(new_project_path)
	if err == nil {
		fmt.Println("Error - Project alredy added!")
		return err
	}

	file, err := os.Create(new_project_path)
	if err != nil {
		fmt.Println("Error - New project file not created")
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	template := createJsonTemplate(project_name)

	encoder.Encode(template)

	fmt.Println("Created new JSON file")

	return nil
}

func createJsonTemplate(project_name string) *Project {
	return &Project{
		Name:        project_name,
		Total_todos: 0,
	}
}

func formatDataFolderPath() error {
	data_folder_name := ".pocketTODO"

	home_path, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	data_folder_path = filepath.Join(home_path, data_folder_name)

	return nil
}
