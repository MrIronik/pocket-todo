package scanner

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"pocket-todo/internal/storage"
	"strings"
)

/* API */

func ScanProject(project_path string) ([]storage.File, error) {
	all_project_files_paths, err := scanDirForSourceFiles(project_path)
	if err != nil {
		fmt.Println("Error - Could not scan project folder path")
		return nil, err
	}

	var all_files []storage.File

	for _, file_path := range all_project_files_paths {
		new_file, err := scanFileForTodos(file_path)
		if err != nil {
			return nil, err
		} else if new_file != nil {
			all_files = append(all_files, *new_file)
		}
	}

	return all_files, nil
}

func ScanDataDir() ([]string, error) {
	data_dir_path := storage.Data_dir_path
	var data_file_names []string

	temp_file_names, err := os.ReadDir(data_dir_path)
	if err != nil {
		return nil, err
	}

	for _, file := range temp_file_names {
		if strings.HasSuffix(file.Name(), "_project.json") {
			clean_name := strings.TrimSuffix(file.Name(), "_project.json")
			data_file_names = append(data_file_names, clean_name)
		}
	}
	return data_file_names, nil
}

/* Helper Functions */

func isSourceFile(path string) bool {
	source_file := map[string]struct{}{
		".c":  {},
		".h":  {},
		".py": {},
	}

	file_ext := filepath.Ext(path)

	_, isSource := source_file[file_ext]
	return isSource
}

func scanDirForSourceFiles(dir_path string) ([]string, error) {
	var source_files_path []string

	filepath.WalkDir(dir_path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if isSourceFile(path) {
			source_files_path = append(source_files_path, path)
		}

		return nil
	})

	return source_files_path, nil
}

func scanFileForTodos(file_path string) (*storage.File, error) {
	scanned_file := storage.File{
		Name:            filepath.Base(file_path),
		Path:            file_path,
		Number_of_todos: 0,
	}

	file, err := os.Open(scanned_file.Path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	line_number := 0

	for scanner.Scan() {
		line_number++
		line := scanner.Text()

		if strings.Contains(line, "TODO: ") {
			new_todo := storage.Todo{
				Line:    line_number,
				Content: strings.TrimPrefix(line, "TODO: "),
			}

			scanned_file.Todos = append(scanned_file.Todos, new_todo)
			scanned_file.Number_of_todos++
		}
	}

	err = scanner.Err()
	if err != nil {
		return nil, err
	}

	if len(scanned_file.Todos) == 0 {
		return nil, nil
	}

	return &scanned_file, nil
}
