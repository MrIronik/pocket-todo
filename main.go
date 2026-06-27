package main

import (
	"os"
	"pocket-todo/cmd/cli"
	"pocket-todo/internal/storage"
	"pocket-todo/internal/ui"
)

func main() {
	storage.Init()

	cli.FlagHandle(cli.FlagParse())

	ui.Run()

	os.Exit(0)
}
