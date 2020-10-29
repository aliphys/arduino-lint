package main

import (
	"fmt"
	"os"

	"github.com/arduino/arduino-check/check"
	"github.com/arduino/arduino-check/configuration"
	"github.com/arduino/arduino-check/project"
	"github.com/arduino/arduino-check/result"
	"github.com/arduino/arduino-check/result/feedback"
)

func main() {
	configuration.Initialize()
	// Must be called after configuration.Initialize()
	result.Results.Initialize()

	projects, err := project.FindProjects()
	if err != nil {
		feedback.Errorf("Error while finding projects: %v", err)
		os.Exit(1)
	}

	for _, project := range projects {
		check.RunChecks(project)
	}

	// All projects have been checked, so summarize their check results in the report.
	result.Results.AddSummary()

	if configuration.OutputFormat() == "text" {
		if len(projects) > 1 {
			// There are multiple projects, print the summary of check results for all projects.
			fmt.Print(result.Results.SummaryText())
		}
	} else {
		// Print the complete JSON formatted report.
		fmt.Println(result.Results.JSONReport())
	}

	if configuration.ReportFilePath() != nil {
		// Write report file.
		result.Results.WriteReport()
	}

	if !result.Results.Passed() {
		os.Exit(1)
	}
}
