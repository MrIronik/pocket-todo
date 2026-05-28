package main

import (
	"pocket-todo/internal/storage"
)

func main() {
	storage.Init()

	storage.CreateNewProjectFile("test")

}
