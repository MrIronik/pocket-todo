package main

import (
	"os"
	"pocket-todo/internal/storage"
	"pocket-todo/internal/ui"
)

func main() {
	storage.Init()

	// TODO: Add arguments hadnling to add new project by the -a <path> flag

	// wd_path, err := filepath.Abs(".")
	// if err != nil {
	// 	fmt.Println("Error - Someting wrong with paths")
	// }

	// project_name := filepath.Base(wd_path)
	// storage.CreateNewProjectFile(project_name, wd_path)

	// project := storage.ReadProjectFile(project_name)

	// new_data, err := scanner.ScanProject(project.Project_path)
	// storage.SaveProjectFileData(project.Name, new_data)

	ui.Run()

	os.Exit(0)
}
