package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/nwalker85/dt/internal/osascript"
)

var createDB string
var createType string
var createTags []string
var createContent string

var createCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new document",
	Long: `Create a new document in DEVONthink.

Content can be provided via --content flag or piped from stdin.

Examples:
  dt create "Meeting Notes" --db "Work"
  dt create "Quick Note" --type markdown --tags work,meeting
  echo "Hello world" | dt create "Note" --type txt`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		// Get content from flag or stdin
		content := createContent
		if content == "" {
			stat, _ := os.Stdin.Stat()
			if (stat.Mode() & os.ModeCharDevice) == 0 {
				bytes, err := io.ReadAll(os.Stdin)
				if err != nil {
					return fmt.Errorf("failed to read stdin: %w", err)
				}
				content = string(bytes)
			}
		}

		// Escape content for AppleScript
		content = strings.ReplaceAll(content, `\`, `\\`)
		content = strings.ReplaceAll(content, `"`, `\"`)

		// Build tags list
		tagsScript := ""
		if len(createTags) > 0 {
			tagsScript = fmt.Sprintf(`set tags of newRecord to {"%s"}`, strings.Join(createTags, `", "`))
		}

		var script string
		if createDB != "" {
			script = fmt.Sprintf(`
tell application "DEVONthink"
	set targetDB to database "%s"
	if targetDB is missing value then
		error "Database not found: %s"
	end if
	set newRecord to create record with {name:"%s", type:%s, content:"%s"} in root of targetDB
	%s
	return uuid of newRecord
end tell`, createDB, createDB, name, getRecordType(createType), content, tagsScript)
		} else {
			script = fmt.Sprintf(`
tell application "DEVONthink"
	set targetDB to first database
	set newRecord to create record with {name:"%s", type:%s, content:"%s"} in root of targetDB
	%s
	return uuid of newRecord
end tell`, name, getRecordType(createType), content, tagsScript)
		}

		out, err := osascript.Run(script)
		if err != nil {
			return fmt.Errorf("create failed: %w", err)
		}

		if jsonOutput {
			fmt.Printf(`{"uuid": "%s", "name": "%s"}%s`, out, name, "\n")
		} else {
			fmt.Printf("Created: %s (uuid: %s)\n", name, out)
		}
		return nil
	},
}

func getRecordType(t string) string {
	switch strings.ToLower(t) {
	case "md", "markdown":
		return "markdown"
	case "txt", "text":
		return "txt"
	case "rtf":
		return "rtf"
	case "html":
		return "html"
	case "bookmark":
		return "bookmark"
	case "sheet":
		return "sheet"
	default:
		return "markdown"
	}
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVar(&createDB, "db", "", "Target database name")
	createCmd.Flags().StringVar(&createType, "type", "markdown", "Document type (markdown, txt, rtf, html, bookmark, sheet)")
	createCmd.Flags().StringSliceVar(&createTags, "tags", nil, "Tags to apply (comma-separated)")
	createCmd.Flags().StringVar(&createContent, "content", "", "Document content")
}
