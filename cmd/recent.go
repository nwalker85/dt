package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/nwalker85/dt/internal/osascript"
)

var recentDays int

var recentCmd = &cobra.Command{
	Use:   "recent",
	Short: "List recently added/modified items",
	Long:  `List items that were recently added or modified across all databases.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		script := fmt.Sprintf(`
tell application "DEVONthink"
	set cutoffDate to (current date) - (%d * days)
	set recentItems to {}
	repeat with db in databases
		set items to contents of db
		repeat with item in items
			if modification date of item > cutoffDate then
				set end of recentItems to {name:name of item, path:path of item, uuid:uuid of item, modified:(modification date of item as string)}
			end if
		end repeat
	end repeat
	return recentItems
end tell`, recentDays)

		out, err := osascript.Run(script)
		if err != nil {
			return fmt.Errorf("recent failed: %w", err)
		}

		if out == "" || out == "{}" {
			return nil
		}

		if jsonOutput {
			// Parse AppleScript record list into JSON
			items := parseRecentItems(out)
			jsonBytes, err := json.MarshalIndent(items, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(jsonBytes))
		} else {
			items := parseRecentItems(out)
			for _, item := range items {
				fmt.Printf("%s\t%s\n", item["name"], item["path"])
			}
		}

		return nil
	},
}

func parseRecentItems(out string) []map[string]string {
	var items []map[string]string
	// AppleScript returns: {{name:x, path:y, uuid:z, modified:w}, ...}
	// Simple parsing - split by "}, {" pattern
	out = strings.TrimPrefix(out, "{{")
	out = strings.TrimSuffix(out, "}}")
	if out == "" {
		return items
	}

	records := strings.Split(out, "}, {")
	for _, record := range records {
		item := make(map[string]string)
		parts := strings.Split(record, ", ")
		for _, part := range parts {
			kv := strings.SplitN(part, ":", 2)
			if len(kv) == 2 {
				item[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
			}
		}
		if len(item) > 0 {
			items = append(items, item)
		}
	}
	return items
}

func init() {
	rootCmd.AddCommand(recentCmd)
	recentCmd.Flags().IntVar(&recentDays, "days", 7, "Number of days to look back")
}
