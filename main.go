package main

import (
	"fmt"
	"os"
	"path/filepath"
	"pocket-todo/internal/scanner"
	"pocket-todo/internal/storage"
)

func main() {
	storage.Init()

	wd_path, err := filepath.Abs(".")
	if err != nil {
		fmt.Println("Error - Someting wrong with paths")
	}

	project_name := filepath.Base(wd_path)
	storage.CreateNewProjectFile(project_name, wd_path)

	project := storage.ReadProjectFile(project_name)

	new_data, err := scanner.ScanProject(project.Project_path)
	storage.SaveProjectFileData(project.Name, new_data)

	os.Exit(0)
}
