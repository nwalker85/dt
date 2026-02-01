package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/walkern/dt/internal/osascript"
)

var openCmd = &cobra.Command{
	Use:   "open <query>",
	Short: "Open matching items in DEVONthink",
	Long: `Open items matching a query in DEVONthink.

Examples:
  dt open "name:report"
  dt open "kind:pdf tag:important"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]

		script := fmt.Sprintf(`
tell application "DEVONthink"
	activate
	set results to search "%s"
	set openedCount to 0
	repeat with r in results
		open tab for record r
		set openedCount to openedCount + 1
	end repeat
	return openedCount
end tell`, strings.ReplaceAll(query, `"`, `\"`))

		out, err := osascript.Run(script)
		if err != nil {
			return fmt.Errorf("open failed: %w", err)
		}

		fmt.Printf("Opened %s items\n", out)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(openCmd)
}
