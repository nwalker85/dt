package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/nwalker85/dt/internal/osascript"
)

var untagCmd = &cobra.Command{
	Use:   "untag <query> <tags...>",
	Short: "Remove tags from items matching a query",
	Long: `Remove one or more tags from all items matching a DEVONthink query.

Examples:
  dt untag "kind:pdf" old-tag
  dt untag "name:report" draft review`,
	Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]
		tagsToRemove := args[1:]

		// Build AppleScript tag list
		tagList := `{"` + strings.Join(tagsToRemove, `", "`) + `"}`

		script := fmt.Sprintf(`
tell application "DEVONthink"
	set results to search "%s"
	set removeTags to %s
	set untaggedCount to 0
	repeat with r in results
		set currentTags to tags of r
		set newTags to {}
		repeat with t in currentTags
			set shouldKeep to true
			repeat with rt in removeTags
				if t as string is equal to rt as string then
					set shouldKeep to false
					exit repeat
				end if
			end repeat
			if shouldKeep then
				set end of newTags to t
			end if
		end repeat
		set tags of r to newTags
		set untaggedCount to untaggedCount + 1
	end repeat
	return untaggedCount
end tell`, strings.ReplaceAll(query, `"`, `\"`), tagList)

		out, err := osascript.Run(script)
		if err != nil {
			return fmt.Errorf("untag failed: %w", err)
		}

		fmt.Printf("Removed tags [%s] from %s items\n", strings.Join(tagsToRemove, ", "), out)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(untagCmd)
}
