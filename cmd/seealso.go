package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/nwalker85/dt/internal/osascript"
)

var seeAlsoLimit int

var seeAlsoCmd = &cobra.Command{
	Use:   "see-also <uuid>",
	Short: "Find related documents",
	Long: `Use DEVONthink's "See Also" feature to find documents related to a specific item.

Examples:
  dt see-also ABC123-DEF456
  dt see-also ABC123-DEF456 --limit 5`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		uuid := args[0]

		script := fmt.Sprintf(`
tell application "DEVONthink"
	set theRecord to get record with uuid "%s"
	if theRecord is missing value then
		return "null"
	end if
	set relatedDocs to compare record theRecord
	set results to {}
	set maxItems to %d
	set itemCount to 0
	repeat with r in relatedDocs
		if itemCount ≥ maxItems then exit repeat
		set end of results to (name of r) & "	" & (uuid of r) & "	" & (path of r)
		set itemCount to itemCount + 1
	end repeat
	return results
end tell`, uuid, seeAlsoLimit)

		out, err := osascript.Run(script)
		if err != nil {
			return fmt.Errorf("see-also failed: %w", err)
		}

		if out == "null" {
			return fmt.Errorf("item not found: %s", uuid)
		}

		if jsonOutput {
			// Parse tab-separated output into JSON
			lines := strings.Split(out, ", ")
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
			lines := strings.Split(out, ", ")
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

func init() {
	rootCmd.AddCommand(seeAlsoCmd)
	seeAlsoCmd.Flags().IntVar(&seeAlsoLimit, "limit", 10, "Maximum number of related items to return")
}
