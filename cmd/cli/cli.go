package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"pocket-todo/internal/scanner"
	"pocket-todo/internal/storage"
	"strings"
)

func FlagParse() Flags {
	var flg Flags

	flag.StringVar(&flg.Add, "add", "", "Add project path <project_path>\n")
	flag.StringVar(&flg.Delete, "delete", "", "Delete project from pocket-todo <project_path>\n")
	flag.BoolVar(&flg.Help, "help", false, "Show Commands\n")

	flag.Parse()

	return flg
}

func FlagHandle(flg Flags) {
	if flg.Add != "" {
		addFlagHandle(flg.Add)
	}

	if flg.Delete != "" {
		deleteFlagHandle()
	}

	if flg.Help {
		helpFlagHandle()
	}
}

func addFlagHandle(path string) {
	new_project_path, _ := normalizePath(path)
	new_project_name := filepath.Base(new_project_path)
	storage.CreateNewProjectFile(new_project_name, new_project_path)

	project := storage.ReadProjectFile(new_project_name)

	new_data, _ := scanner.ScanProject(project.Project_path)
	storage.SaveProjectFileData(project.Name, new_data)

	fmt.Printf("Added %s to projects\n", new_project_name)

	os.Exit(0)
}

func deleteFlagHandle() {
	// TODO: Implement delete flag handle
}

func helpFlagHandle() {
	fmt.Printf("Pocket-todo all commands and arguments:\n\n")
	flag.PrintDefaults()

	os.Exit(0)
}

func normalizePath(p string) (string, error) {
	// 1. expand ~
	if p == "~" || strings.HasPrefix(p, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		if p == "~" {
			p = home
		} else {
			p = filepath.Join(home, p[2:])
		}
	}

	// 2. clean (usuwa ../, ./ itd.)
	p = filepath.Clean(p)

	// 3. zamiana na absolute (jeśli nie jest)
	abs, err := filepath.Abs(p)
	if err != nil {
		return "", err
	}

	return abs, nil
}
