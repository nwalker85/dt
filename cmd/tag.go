package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/nwalker85/dt/internal/osascript"
)

var tagCmd = &cobra.Command{
	Use:   "tag <query> <tags...>",
	Short: "Batch tag items matching a query",
	Long: `Add one or more tags to all items matching a DEVONthink query.

Examples:
  dt tag "kind:pdf" important
  dt tag "name:report" work review urgent`,
	Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]
		tags := args[1:]

		// Build AppleScript tag list
		tagList := `{"` + strings.Join(tags, `", "`) + `"}`

		script := fmt.Sprintf(`
tell application "DEVONthink"
	set results to search "%s"
	set newTags to %s
	set taggedCount to 0
	repeat with r in results
		set currentTags to tags of r
		repeat with t in newTags
			if currentTags does not contain t then
				set end of currentTags to t
			end if
		end repeat
		set tags of r to currentTags
		set taggedCount to taggedCount + 1
	end repeat
	return taggedCount
end tell`, strings.ReplaceAll(query, `"`, `\"`), tagList)

		out, err := osascript.Run(script)
		if err != nil {
			return fmt.Errorf("tagging failed: %w", err)
		}

		fmt.Printf("Tagged %s items with: %s\n", out, strings.Join(tags, ", "))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(tagCmd)
}
