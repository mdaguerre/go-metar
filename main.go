// Package main is the entry point for the go-metar CLI application.
// In Go, every executable must have a "main" package with a "main" function.
package main

import (
	"fmt"
	"os"

	// Cobra is the most popular library for building CLI apps in Go.
	// It handles argument parsing, flags, help text, and subcommands.
	"github.com/spf13/cobra"

	// This imports our own "metar" package from this project.
	// The path matches what we defined in go.mod + the folder name.
	"github.com/mdaguerre/go-metar/metar"
)

// These variables hold our CLI flag values.
// In Go, package-level variables are declared outside functions.
var (
	rawOutput bool
	allOutput bool
)

func main() {
	// Create the root command - this is what runs when user types "go-metar"
	rootCmd := &cobra.Command{
		Use:   "go-metar [ICAO]",          // How to use the command
		Short: "Fetch METAR weather data", // Brief description
		Long: `go-metar fetches METAR aviation weather reports for any airport.

Examples:
  go-metar KJFK        # Get decoded METAR for JFK airport
  go-metar EGLL --raw  # Get raw METAR for London Heathrow
  go-metar LFPG --all  # Get both raw and decoded for Paris CDG`,

		// Args: cobra.ExactArgs(1) means the command requires exactly 1 argument
		Args: cobra.ExactArgs(1),

		// Run is the function that executes when the command is called.
		// It receives the command itself and the positional arguments (args).
		Run: func(cmd *cobra.Command, args []string) {
			icao := args[0] // First argument is the ICAO code

			// Fetch METAR data from the API
			data, err := metar.Fetch(icao)
			if err != nil {
				// fmt.Fprintf writes formatted output to a specific writer (os.Stderr here)
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1) // Exit with error code
			}

			// Handle output based on flags
			if rawOutput {
				fmt.Println(data.Raw)
			} else if allOutput {
				fmt.Println("Raw METAR:")
				fmt.Println(data.Raw)
				fmt.Println("\nDecoded:")
				fmt.Println(metar.Decode(data))
			} else {
				// Default: show decoded output
				fmt.Println(metar.Decode(data))
			}
		},
	}

	// Add flags to the command
	// Flags().BoolVarP connects a boolean variable to a flag
	// Parameters: variable pointer, long name, short name, default value, description
	rootCmd.Flags().BoolVarP(&rawOutput, "raw", "r", false, "Show raw METAR string only")
	rootCmd.Flags().BoolVarP(&allOutput, "all", "a", false, "Show both raw and decoded output")

	// Execute the command - this parses arguments and runs the appropriate function
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
