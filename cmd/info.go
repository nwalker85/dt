package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/walkern/dt/internal/osascript"
)

var infoCmd = &cobra.Command{
	Use:   "info <uuid>",
	Short: "Show detailed item information",
	Long:  `Display detailed information about a specific item by its UUID.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		uuid := args[0]

		script := fmt.Sprintf(`
var app = Application("DEVONthink");
var record = app.getRecordWithUuid("%s");
if (record) {
	var info = {
		name: record.name(),
		uuid: record.uuid(),
		path: record.path(),
		kind: record.kind(),
		size: record.size(),
		tags: record.tags(),
		rating: record.rating(),
		label: record.label(),
		comment: record.comment(),
		url: record.url(),
		creationDate: record.creationDate() ? record.creationDate().toISOString() : null,
		modificationDate: record.modificationDate() ? record.modificationDate().toISOString() : null,
		additionDate: record.additionDate() ? record.additionDate().toISOString() : null,
		wordCount: record.wordCount(),
		characterCount: record.characterCount(),
		indexed: record.indexed(),
		duplicates: record.duplicates().length
	};
	JSON.stringify(info, null, 2);
} else {
	"null";
}`, uuid)

		out, err := osascript.RunJS(script)
		if err != nil {
			return fmt.Errorf("info failed: %w", err)
		}

		if out == "null" {
			return fmt.Errorf("item not found: %s", uuid)
		}

		fmt.Println(out)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
