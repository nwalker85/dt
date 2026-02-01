package cmd

import (
	"fmt"

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
var app = Application("DEVONthink");
var record = app.getRecordWithUuid("%s");
if (!record) {
	"null";
} else {
	var related = app.seeAlso({record: record});
	var results = [];
	var limit = %d;
	for (var i = 0; i < Math.min(related.length, limit); i++) {
		var r = related[i];
		results.push({
			name: r.name(),
			uuid: r.uuid(),
			path: r.path(),
			score: r.score ? r.score() : null
		});
	}
	JSON.stringify(results, null, 2);
}`, uuid, seeAlsoLimit)

		out, err := osascript.RunJS(script)
		if err != nil {
			return fmt.Errorf("see-also failed: %w", err)
		}

		if out == "null" {
			return fmt.Errorf("item not found: %s", uuid)
		}

		fmt.Println(out)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(seeAlsoCmd)
	seeAlsoCmd.Flags().IntVar(&seeAlsoLimit, "limit", 10, "Maximum number of related items to return")
}
