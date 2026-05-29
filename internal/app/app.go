package app

import (
	"pocket-todo/internal/scanner"
	"pocket-todo/internal/storage"
)

func Init() error {
	storage.Init()

	return nil
}

func gatProjectNames() ([]string, error) {
	all_project_names, err := scanner.ScanDataDir()
	if err != nil {
		return nil, err
	}

	return all_project_names, nil
}
