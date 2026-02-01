package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var jsonOutput bool

var rootCmd = &cobra.Command{
	Use:   "dt",
	Short: "DEVONthink CLI - Automate DEVONthink via AppleScript/JXA",
	Long: `dt is a command-line interface for DEVONthink automation.

It wraps DEVONthink's AppleScript and JXA capabilities to provide
quick access to search, tagging, import, and export operations.

DEVONthink Query Syntax:
  kind:pdf          Search by file type
  tag:mytag         Search by tag
  name:filename     Search by filename
  content:term      Search within content
  date:today        Search by date
  size:>1mb         Search by size

  Boolean operators: AND, OR, NOT

Examples:
  dt search "kind:pdf tag:work"
  dt stats
  dt tags list
  dt tag "kind:pdf" important review
  dt export important --dest ~/exports
  dt import ~/Documents/reports --db "Work"`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")
}
