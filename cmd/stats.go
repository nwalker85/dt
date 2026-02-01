package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/nwalker85/dt/internal/osascript"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show item count per database",
	Long:  `Display the number of items in each DEVONthink database.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		script := `
tell application "DEVONthink"
	set output to ""
	repeat with db in databases
		set dbName to name of db
		set itemCount to count of contents of db
		set output to output & dbName & ":" & itemCount & linefeed
	end repeat
	return output
end tell`

		out, err := osascript.Run(script)
		if err != nil {
			return fmt.Errorf("stats failed: %w", err)
		}

		if jsonOutput {
			stats := make(map[string]int)
			lines := strings.Split(strings.TrimSpace(out), "\n")
			for _, line := range lines {
				if line == "" {
					continue
				}
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					count, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
					stats[strings.TrimSpace(parts[0])] = count
				}
			}
			jsonBytes, err := json.MarshalIndent(stats, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(jsonBytes))
		} else {
			lines := strings.Split(strings.TrimSpace(out), "\n")
			for _, line := range lines {
				if line == "" {
					continue
				}
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					fmt.Printf("%s: %s items\n", strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
				}
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
}
