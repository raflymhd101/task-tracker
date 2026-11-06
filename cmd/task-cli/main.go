package main

import (
	"fmt"
	"os"

	"task-tracker/internal/handler"
	"task-tracker/internal/storage"
	"task-tracker/internal/task"
)

const defaultStoragePath = "tasks.json"

func main() {
	store := storage.New(defaultStoragePath)
	svc := task.NewService(store)

	args := os.Args[1:]
	if len(args) == 0 {
		printUsage()
		os.Exit(1)
	}

	command := args[0]

	switch command {
	case "add":
		handler.HandleAdd(svc, args[1:])
	case "update":
		handler.HandleUpdate(svc, args[1:])
	case "delete":
		handler.HandleDelete(svc, args[1:])
	case "mark-in-progress":
		handler.HandleMark(svc, args[1:], task.StatusInProgress)
	case "mark-done":
		handler.HandleMark(svc, args[1:], task.StatusDone)
	case "list":
		handler.HandleList(svc, args[1:])
	case "help", "--help", "-h":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Error: unknown command %q\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`Task Tracker CLI - roadmap.sh

Usage:
  task-cli add <description>                Add a new task
  task-cli update <id> <description>        Update task description
  task-cli delete <id>                      Delete a task
  task-cli mark-in-progress <id>            Mark task as in-progress
  task-cli mark-done <id>                   Mark task as done
  task-cli list                             List all tasks
  task-cli list <done|todo|in-progress>     List tasks filtered by status
  task-cli help                             Show this help message`)
}
