package cmd

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/walkern/dt/internal/osascript"
)

var tagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "Tag-related commands",
	Long:  `Commands for working with DEVONthink tags.`,
}

var tagsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all unique tags",
	Long:  `List all unique tags across all DEVONthink databases.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		script := `
var app = Application("DEVONthink");
var dbs = app.databases();
var tagSet = {};
for (var i = 0; i < dbs.length; i++) {
	var contents = dbs[i].contents();
	for (var j = 0; j < contents.length; j++) {
		var tags = contents[j].tags();
		if (tags) {
			for (var k = 0; k < tags.length; k++) {
				tagSet[tags[k]] = true;
			}
		}
	}
}
Object.keys(tagSet).join("\n");`

		out, err := osascript.RunJS(script)
		if err != nil {
			return fmt.Errorf("tags list failed: %w", err)
		}

		if out == "" {
			return nil
		}

		tags := strings.Split(out, "\n")
		sort.Strings(tags)

		if jsonOutput {
			jsonBytes, err := json.MarshalIndent(tags, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(jsonBytes))
		} else {
			for _, tag := range tags {
				if tag != "" {
					fmt.Println(tag)
				}
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(tagsCmd)
	tagsCmd.AddCommand(tagsListCmd)
}
