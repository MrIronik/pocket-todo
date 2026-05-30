package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Project struct {
	Name         string
	Project_path string
	Total_todos  int
	Files        []File
}

type File struct {
	Name            string
	Path            string
	Number_of_todos int
	Todos           []Todo
}

type Todo struct {
	Line    int
	Content string
}

/* Variables for Data Manager */
var Data_dir_path string

func Init() error {
	err := formatDataFolderPath()
	if err != nil {
		fmt.Println("Init Error - Getting Path")
		return err
	}

	err = os.MkdirAll(Data_dir_path, 0755)
	if err != nil {
		return err
	}

	return nil
}

func CreateNewProjectFile(project_name string, project_path string) error {
	project_file_path := formatProjectFilePath(project_name)

	_, err := os.Stat(project_file_path)
	if err == nil {
		fmt.Println("Error - Project alredy added!")
		return err
	}

	file, err := os.Create(project_file_path)
	if err != nil {
		fmt.Println("Error - New project file not created")
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	template := createJsonTemplate(project_name, project_path)

	encoder.Encode(template)

	fmt.Println("Created new JSON file")

	return nil
}

func DeleteProjectFile(data Project) error {
	file_name := data.Name + "_project.json"
	project_file_path := filepath.Join(Data_dir_path, file_name)

	err := os.Remove(project_file_path)
	if err != nil {
		fmt.Println("Error - Project file not deleted")
		return err
	}

	fmt.Printf("%s deleted\n", data.Name)

	return nil
}

func SaveProjectFileData(project_name string, new_files []File) error {
	project := ReadProjectFile(project_name)
	project_file_path := formatProjectFilePath(project_name)

	_, err := os.Stat(project_file_path)
	if err != nil {
		fmt.Println("Error - No Project File Path!")
		return err
	}

	project.Total_todos = 0
	project.Files = new_files

	for _, file := range project.Files {
		project.Total_todos += file.Number_of_todos
	}

	file, err := os.Create(project_file_path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	encoder.Encode(project)

	return nil
}

func ReadProjectFile(project_name string) Project {
	project_file_path := formatProjectFilePath(project_name)

	file, err := os.Open(project_file_path)
	if err != nil {
		fmt.Println("Error - file not open")
		// TODO: Handle os.Open error correctly
	}
	defer file.Close()

	var project Project

	err = json.NewDecoder(file).Decode(&project)
	if err != nil {
		panic(err)
	}

	return project
}

func createJsonTemplate(project_name string, project_path string) *Project {
	return &Project{
		Name:         project_name,
		Project_path: project_path,
		Total_todos:  0,
	}
}

func formatDataFolderPath() error {
	data_folder_name := ".pocketTODO"

	home_path, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	Data_dir_path = filepath.Join(home_path, data_folder_name)

	return nil
}

func formatProjectFilePath(project_name string) string {
	file_name := project_name + "_project.json"
	project_file_path := filepath.Join(Data_dir_path, file_name)

	return project_file_path
}
