package handler

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"task-tracker/internal/task"
	"text/tabwriter"
)

func HandleAdd(svc *task.Service, args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Error: missing task description")
		fmt.Fprintln(os.Stderr, "Usage: task-cli add <description>")
		os.Exit(1)
	}

	t, err := svc.Add(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error adding task: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Task added successfully (ID: %d)\n", t.ID)
}

func HandleUpdate(svc *task.Service, args []string) {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "Error: missing arguments")
		fmt.Fprintln(os.Stderr, "Usage: task-cli update <id> <new description>")
		os.Exit(1)
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: invalid task ID %q (must be a number)\n", args[0])
		os.Exit(1)
	}

	if err := svc.Update(id, args[1]); err != nil {
		handleServiceError(err, id)
	}

	fmt.Printf("Task updated successfully (ID: %d)\n", id)
}

func HandleDelete(svc *task.Service, args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Error: missing task ID")
		fmt.Fprintln(os.Stderr, "Usage: task-cli delete <id>")
		os.Exit(1)
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: invalid task ID %q (must be a number)\n", args[0])
		os.Exit(1)
	}

	if err := svc.Delete(id); err != nil {
		handleServiceError(err, id)
	}

	fmt.Printf("Task deleted successfully (ID: %d)\n", id)
}

func HandleMark(svc *task.Service, args []string, status task.Status) {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Error: missing task ID\n")
		fmt.Fprintf(os.Stderr, "Usage: task-cli mark-%s <id>\n", status)
		os.Exit(1)
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: invalid task ID %q (must be a number)\n", args[0])
		os.Exit(1)
	}

	if err := svc.MarkStatus(id, status); err != nil {
		handleServiceError(err, id)
	}

	fmt.Printf("Task marked as %s (ID: %d)\n", status, id)
}

func HandleList(svc *task.Service, args []string) {
	var filter task.Status
	if len(args) > 0 {
		filter = task.Status(args[0])
	}

	tasks, err := svc.List(filter)
	if err != nil {
		if errors.Is(err, task.ErrInvalidStatus) {
			fmt.Fprintf(os.Stderr, "Error: invalid status filter %q (allowed: todo, in-progress, done)\n", filter)
		} else {
			fmt.Fprintf(os.Stderr, "Error listing tasks: %v\n", err)
		}
		os.Exit(1)
	}

	if len(tasks) == 0 {
		if filter != "" {
			fmt.Printf("No tasks found with status %q.\n", filter)
		} else {
			fmt.Println("No tasks found. Add one using: task-cli add <description>")
		}
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tSTATUS\tDESCRIPTION\tCREATED AT\tUPDATED AT")
	fmt.Fprintln(w, "--\t------\t-----------\t----------\t----------")

	for _, t := range tasks {
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n",
			t.ID,
			t.Status,
			t.Description,
			t.CreatedAt.Format("2006-01-02 15:04"),
			t.UpdatedAt.Format("2006-01-02 15:04"),
		)
	}
	w.Flush()
}

func handleServiceError(err error, id int) {
	if errors.Is(err, task.ErrTaskNotFound) {
		fmt.Fprintf(os.Stderr, "Error: task with ID %d not found\n", id)
	} else if errors.Is(err, task.ErrEmptyDescription) {
		fmt.Fprintln(os.Stderr, "Error: description cannot be empty")
	} else {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
	os.Exit(1)
}
