package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/walkern/dt/internal/osascript"
)

var summaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "Show database summary as JSON",
	Long:  `Display a detailed summary of all DEVONthink databases in JSON format.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		script := `
var app = Application("DEVONthink");
var dbs = app.databases();
var result = [];
for (var i = 0; i < dbs.length; i++) {
	var db = dbs[i];
	result.push({
		name: db.name(),
		uuid: db.uuid(),
		path: db.path(),
		itemCount: db.contents().length
	});
}
JSON.stringify(result, null, 2);`

		out, err := osascript.RunJS(script)
		if err != nil {
			return fmt.Errorf("summary failed: %w", err)
		}

		fmt.Println(out)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(summaryCmd)
}
