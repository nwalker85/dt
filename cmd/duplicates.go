package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/walkern/dt/internal/osascript"
)

var duplicatesCmd = &cobra.Command{
	Use:   "duplicates",
	Short: "Find duplicate items",
	Long:  `Find and list duplicate items across all databases.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		script := `
var app = Application("DEVONthink");
var dbs = app.databases();
var seen = {};
var duplicates = [];

for (var i = 0; i < dbs.length; i++) {
	var contents = dbs[i].contents();
	for (var j = 0; j < contents.length; j++) {
		var item = contents[j];
		var dupes = item.duplicates();
		if (dupes.length > 0) {
			var key = item.uuid();
			if (!seen[key]) {
				seen[key] = true;
				var group = {
					name: item.name(),
					uuid: item.uuid(),
					path: item.path(),
					duplicateCount: dupes.length,
					duplicates: []
				};
				for (var k = 0; k < dupes.length; k++) {
					seen[dupes[k].uuid()] = true;
					group.duplicates.push({
						name: dupes[k].name(),
						uuid: dupes[k].uuid(),
						path: dupes[k].path()
					});
				}
				duplicates.push(group);
			}
		}
	}
}
JSON.stringify(duplicates, null, 2);`

		out, err := osascript.RunJS(script)
		if err != nil {
			return fmt.Errorf("duplicates failed: %w", err)
		}

		fmt.Println(out)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(duplicatesCmd)
}
