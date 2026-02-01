package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/nwalker85/dt/internal/osascript"
)

var inboxCmd = &cobra.Command{
	Use:   "inbox",
	Short: "List inbox items",
	Long:  `List items in the global inbox and database inboxes.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		script := `
var app = Application("DEVONthink");
var inboxItems = [];

// Global inbox
var globalInbox = app.inbox();
if (globalInbox) {
	var contents = globalInbox.contents();
	for (var i = 0; i < contents.length; i++) {
		var item = contents[i];
		inboxItems.push({
			name: item.name(),
			uuid: item.uuid(),
			kind: item.kind(),
			database: "Global Inbox",
			additionDate: item.additionDate() ? item.additionDate().toISOString() : null
		});
	}
}

// Database inboxes
var dbs = app.databases();
for (var d = 0; d < dbs.length; d++) {
	var db = dbs[d];
	var dbInbox = db.inbox();
	if (dbInbox) {
		var dbContents = dbInbox.contents();
		for (var j = 0; j < dbContents.length; j++) {
			var item = dbContents[j];
			inboxItems.push({
				name: item.name(),
				uuid: item.uuid(),
				kind: item.kind(),
				database: db.name(),
				additionDate: item.additionDate() ? item.additionDate().toISOString() : null
			});
		}
	}
}

JSON.stringify(inboxItems, null, 2);`

		out, err := osascript.RunJS(script)
		if err != nil {
			return fmt.Errorf("inbox failed: %w", err)
		}

		if jsonOutput {
			fmt.Println(out)
		} else {
			var items []map[string]interface{}
			if err := json.Unmarshal([]byte(out), &items); err != nil {
				fmt.Println(out)
				return nil
			}

			if len(items) == 0 {
				fmt.Println("Inbox is empty")
				return nil
			}

			for _, item := range items {
				fmt.Printf("[%s] %s (%s)\n", item["database"], item["name"], item["kind"])
			}
			fmt.Printf("\nTotal: %d items\n", len(items))
		}

		return nil
	},
}

var inboxProcessCmd = &cobra.Command{
	Use:   "process",
	Short: "Process inbox items",
	Long: `Interactively process inbox items by classifying and moving them.

For each inbox item, shows classification suggestions and moves
the item to the suggested location.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		script := `
tell application "DEVONthink"
	set inboxItems to contents of inbox
	set processedCount to 0
	repeat with item in inboxItems
		set proposedGroup to classify record item
		if proposedGroup is not missing value then
			move record item to proposedGroup
			set processedCount to processedCount + 1
		end if
	end repeat
	return processedCount
end tell`

		out, err := osascript.Run(script)
		if err != nil {
			return fmt.Errorf("inbox process failed: %w", err)
		}

		fmt.Printf("Processed and classified %s items\n", out)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(inboxCmd)
	inboxCmd.AddCommand(inboxProcessCmd)
}
