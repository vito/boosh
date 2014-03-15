package main

import (
	"fmt"
	"strings"

	"github.com/dynport/gocloud/aws/cloudformation"
	"github.com/hotei/ansiterm"
	"github.com/mgutz/ansi"
)

func renderEvents(events <-chan *cloudformation.StackEvent) bool {
	ok := true

	inFlight := make(map[string]Process)

	printedLines := 0

	for ev := range events {
		status := ev.ResourceStatus

		for i := 0; i < printedLines; i++ {
			fmt.Print("\x1b[1A")
			ansiterm.ClearLine()
		}

		var process Process

		if strings.Contains(status, "ROLLBACK") {
			process = RollbackProcess
		} else if strings.Contains(status, "UPDATE") {
			process = UpdateProcess
		} else if strings.Contains(status, "DELETE") {
			process = DeleteProcess
		} else if strings.Contains(status, "CREATE") {
			process = CreateProcess
		}

		if strings.HasSuffix(ev.ResourceStatus, "IN_PROGRESS") {
			inFlight[ev.LogicalResourceId] = process
		} else {
			delete(inFlight, ev.LogicalResourceId)

			label := process.CompletedLabel()

			if strings.HasSuffix(status, "FAILED") {
				ok = false
				fmt.Printf(ansi.Color("%s: %s (%s)\n", "red"), label, ev.LogicalResourceId, ev.ResourceStatusReason)
			} else {
				fmt.Printf(ansi.Color("%s: %s (%s)\n", "green"), label, ev.LogicalResourceId, ev.PhysicalResourceId)
			}
		}

		printedLines = renderInFlight(inFlight)
	}

	return ok
}

func renderInFlight(inFlight map[string]Process) int {
	byProcess := make(map[Process][]string)

	for id, process := range inFlight {
		byProcess[process] = append(byProcess[process], id)
	}

	printedLines := 0

	for process, ids := range byProcess {
		prefix := process.ActiveLabel() + ": "

		for i, id := range ids {
			if i == 0 {
				fmt.Print(prefix)
			} else {
				fmt.Print(strings.Repeat(" ", len(prefix)))
			}

			fmt.Println(id)

			printedLines++
		}
	}

	return printedLines
}
