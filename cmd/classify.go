package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/walkern/dt/internal/osascript"
)

var classifyApply bool

var classifyCmd = &cobra.Command{
	Use:   "classify <query>",
	Short: "Auto-classify items using DEVONthink AI",
	Long: `Use DEVONthink's AI to suggest or apply classification for items.

By default, shows suggested locations without moving items.
Use --apply to actually move items to suggested locations.

Examples:
  dt classify "tag:inbox"
  dt classify "kind:pdf" --apply`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]

		if classifyApply {
			script := fmt.Sprintf(`
tell application "DEVONthink"
	set results to search "%s"
	set classifiedCount to 0
	repeat with r in results
		set proposedGroup to classify record r
		if proposedGroup is not missing value then
			move record r to proposedGroup
			set classifiedCount to classifiedCount + 1
		end if
	end repeat
	return classifiedCount
end tell`, strings.ReplaceAll(query, `"`, `\"`))

			out, err := osascript.Run(script)
			if err != nil {
				return fmt.Errorf("classify failed: %w", err)
			}

			fmt.Printf("Classified and moved %s items\n", out)
		} else {
			script := fmt.Sprintf(`
var app = Application("DEVONthink");
var results = app.search("%s");
var suggestions = [];
for (var i = 0; i < results.length; i++) {
	var item = results[i];
	var proposed = app.classify(item);
	suggestions.push({
		name: item.name(),
		uuid: item.uuid(),
		suggestedLocation: proposed ? proposed.name() : null
	});
}
JSON.stringify(suggestions, null, 2);`, strings.ReplaceAll(query, `"`, `\"`))

			out, err := osascript.RunJS(script)
			if err != nil {
				return fmt.Errorf("classify failed: %w", err)
			}

			fmt.Println(out)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(classifyCmd)
	classifyCmd.Flags().BoolVar(&classifyApply, "apply", false, "Actually move items to suggested locations")
}
