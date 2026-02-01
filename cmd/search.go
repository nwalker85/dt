package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/nwalker85/dt/internal/osascript"
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search DEVONthink and return file paths",
	Long: `Search across all DEVONthink databases using DEVONthink query syntax.

Returns the file paths of matching items, one per line.

Examples:
  dt search "kind:pdf"
  dt search "tag:important AND kind:markdown"
  dt search "content:project date:thisweek"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]

		script := fmt.Sprintf(`
tell application "DEVONthink"
	set results to search "%s"
	set paths to {}
	repeat with r in results
		set end of paths to path of r
	end repeat
	return paths
end tell`, strings.ReplaceAll(query, `"`, `\"`))

		out, err := osascript.Run(script)
		if err != nil {
			return fmt.Errorf("search failed: %w", err)
		}

		if out == "" {
			return nil
		}

		// AppleScript returns comma-separated list
		paths := strings.Split(out, ", ")

		if jsonOutput {
			jsonBytes, err := json.MarshalIndent(paths, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(jsonBytes))
		} else {
			for _, p := range paths {
				fmt.Println(p)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
