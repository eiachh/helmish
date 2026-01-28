package main

import (
	"flag"
	"os"

	"helmish/pkg/helmishlib"
)

// parseConfig parses command-line flags and environment variables to build Options
func parseConfig() helmishlib.Options {
	// Define flags
	chartPathFlag := flag.String("chart-path", "", "Path to the Helm chart")
	profileNameFlag := flag.String("profile", "", "Profile name")

	flag.Parse()

	// Get from env vars first
	chartPath := os.Getenv("HELMISH_CHART_PATH")
	profileName := os.Getenv("HELMISH_PROFILE")

	// Flags take precedence over env vars
	if *chartPathFlag != "" {
		chartPath = *chartPathFlag
	}
	if *profileNameFlag != "" {
		profileName = *profileNameFlag
	}

	// Default values if still empty
	if chartPath == "" {
		chartPath = "/default/chart/path"
	}
	if profileName == "" {
		profileName = "default"
	}

	return helmishlib.Options{
		Chart: helmishlib.Chart{
			Path: chartPath,
		},
		Profile: helmishlib.Profile{
			Name: profileName,
		},
	}
}