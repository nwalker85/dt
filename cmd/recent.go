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
		// Use DEVONthink's search with date filter - much faster than iterating
		script := fmt.Sprintf(`
tell application "DEVONthink"
	set cutoffDate to (current date) - (%d * days)
	set dateStr to (year of cutoffDate as string) & "-"
	set m to (month of cutoffDate as integer)
	if m < 10 then set dateStr to dateStr & "0"
	set dateStr to dateStr & (m as string) & "-"
	set d to (day of cutoffDate as integer)
	if d < 10 then set dateStr to dateStr & "0"
	set dateStr to dateStr & (d as string)

	set results to search "date:>=" & dateStr
	set recentItems to {}
	set maxItems to 100
	set itemCount to 0
	repeat with r in results
		if itemCount ≥ maxItems then exit repeat
		set end of recentItems to (name of r) & "	" & (uuid of r) & "	" & (path of r)
		set itemCount to itemCount + 1
	end repeat
	return recentItems
end tell`, recentDays)

		out, err := osascript.Run(script)
		if err != nil {
			return fmt.Errorf("recent failed: %w", err)
		}

		if out == "" {
			fmt.Println("No recent items found")
			return nil
		}

		lines := strings.Split(out, ", ")
		if jsonOutput {
			var items []map[string]string
			for _, line := range lines {
				parts := strings.Split(line, "\t")
				if len(parts) >= 3 {
					items = append(items, map[string]string{
						"name": parts[0],
						"uuid": parts[1],
						"path": parts[2],
					})
				}
			}
			jsonBytes, _ := json.MarshalIndent(items, "", "  ")
			fmt.Println(string(jsonBytes))
		} else {
			for _, line := range lines {
				parts := strings.Split(line, "\t")
				if len(parts) >= 1 {
					fmt.Println(parts[0])
				}
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
