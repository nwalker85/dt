package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/nwalker85/dt/internal/osascript"
)

var weeklyReportCmd = &cobra.Command{
	Use:   "weekly-report",
	Short: "Generate a summary of items added this week",
	Long:  `Generate a report showing items added to DEVONthink in the past week, grouped by database.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Use search with date filter - much faster than iterating all items
		script := `
tell application "DEVONthink"
	set cutoffDate to (current date) - (7 * days)
	set dateStr to (year of cutoffDate as string) & "-"
	set m to (month of cutoffDate as integer)
	if m < 10 then set dateStr to dateStr & "0"
	set dateStr to dateStr & (m as string) & "-"
	set d to (day of cutoffDate as integer)
	if d < 10 then set dateStr to dateStr & "0"
	set dateStr to dateStr & (d as string)

	set results to search "date:>=" & dateStr
	set dbCounts to {}

	repeat with r in results
		set dbName to name of database of r
		set found to false
		repeat with i from 1 to count of dbCounts
			if item 1 of item i of dbCounts is dbName then
				set item 2 of item i of dbCounts to (item 2 of item i of dbCounts) + 1
				set found to true
				exit repeat
			end if
		end repeat
		if not found then
			set end of dbCounts to {dbName, 1}
		end if
	end repeat

	set output to ""
	repeat with entry in dbCounts
		set output to output & (item 1 of entry) & "	" & (item 2 of entry) & ","
	end repeat
	return output
end tell`

		out, err := osascript.Run(script)
		if err != nil {
			return fmt.Errorf("weekly-report failed: %w", err)
		}

		if out == "" {
			fmt.Println("No items added in the past week")
			return nil
		}

		// Parse output and format as report
		type dbSummary struct {
			Name  string `json:"name"`
			Count int    `json:"count"`
		}
		type report struct {
			Period    string      `json:"period"`
			Databases []dbSummary `json:"databases"`
			Total     int         `json:"total"`
		}

		r := report{Period: "last 7 days"}
		entries := strings.Split(strings.TrimSuffix(out, ","), ",")
		for _, entry := range entries {
			parts := strings.Split(entry, "\t")
			if len(parts) >= 2 {
				count := 0
				fmt.Sscanf(parts[1], "%d", &count)
				r.Databases = append(r.Databases, dbSummary{Name: parts[0], Count: count})
				r.Total += count
			}
		}

		if jsonOutput {
			jsonBytes, _ := json.MarshalIndent(r, "", "  ")
			fmt.Println(string(jsonBytes))
		} else {
			fmt.Printf("Weekly Report (last 7 days)\n")
			fmt.Printf("===========================\n")
			for _, db := range r.Databases {
				fmt.Printf("  %s: %d items\n", db.Name, db.Count)
			}
			fmt.Printf("---------------------------\n")
			fmt.Printf("  Total: %d items\n", r.Total)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(weeklyReportCmd)
}
