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

// version is set at build time via ldflags by goreleaser.
// See .goreleaser.yaml: -X main.version={{.Version}}
var version = "dev"

// These variables hold our CLI flag values.
// In Go, package-level variables are declared outside functions.
var (
	rawOutput   bool
	allOutput   bool
	showVersion bool
	tafOutput   bool
)

func main() {
	// Create the root command - this is what runs when user types "go-metar"
	rootCmd := &cobra.Command{
		Use:   "go-metar [ICAO...]",       // How to use the command
		Short: "Fetch METAR weather data", // Brief description
		Long: `go-metar fetches METAR aviation weather reports for any airport.

Examples:
  go-metar KJFK              # Get decoded METAR for JFK airport
  go-metar KJFK KLAX EGLL    # Get METARs for multiple airports
  go-metar EGLL --raw        # Get raw METAR for London Heathrow
  go-metar KJFK KLAX --all   # Get both raw and decoded for multiple airports
  go-metar KJFK --taf        # Include TAF forecast`,

		// Run is the function that executes when the command is called.
		// It receives the command itself and the positional arguments (args).
		Run: func(cmd *cobra.Command, args []string) {
			// Handle --version flag
			if showVersion {
				fmt.Printf("go-metar %s\n", version)
				return
			}

			// Validate that we have at least 1 argument when not showing version
			if len(args) < 1 {
				fmt.Fprintln(os.Stderr, "Error: requires at least 1 ICAO code")
				cmd.Usage()
				os.Exit(1)
			}

			// Validate mutually exclusive flags
			if rawOutput && allOutput {
				fmt.Fprintln(os.Stderr, "Error: cannot use both --raw and --all flags")
				os.Exit(1)
			}

			// Fetch METAR data for all airports
			metars, err := metar.FetchMultiple(args)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			// Handle output based on flags
			for i, data := range metars {
				if rawOutput {
					fmt.Println(data.Raw)
				} else if allOutput {
					if i > 0 {
						fmt.Println() // Blank line between airports
					}
					fmt.Printf("Raw METAR (%s):\n", data.StationID)
					fmt.Println(data.Raw)
					fmt.Println("\nDecoded:")
					fmt.Println(metar.Decode(data))
				} else {
					// Default: show decoded output
					if i > 0 {
						fmt.Println() // Blank line between airports
					}
					fmt.Println(metar.Decode(data))
				}
			}

			// Fetch and display TAF if requested
			if tafOutput {
				tafs, err := metar.FetchMultipleTAF(args)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error fetching TAF: %v\n", err)
					os.Exit(1)
				}

				fmt.Println() // Blank line before TAF section
				for i, taf := range tafs {
					if rawOutput {
						fmt.Println(taf.RawTAF)
					} else {
						if i > 0 {
							fmt.Println()
						}
						fmt.Println(metar.DecodeTAF(taf))
					}
				}
			}
		},
	}

	// Add flags to the command
	// Flags().BoolVarP connects a boolean variable to a flag
	// Parameters: variable pointer, long name, short name, default value, description
	rootCmd.Flags().BoolVarP(&rawOutput, "raw", "r", false, "Show raw METAR string only")
	rootCmd.Flags().BoolVarP(&allOutput, "all", "a", false, "Show both raw and decoded output")
	rootCmd.Flags().BoolVarP(&tafOutput, "taf", "t", false, "Include TAF forecast")
	rootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "Show version information")

	// Execute the command - this parses arguments and runs the appropriate function
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
