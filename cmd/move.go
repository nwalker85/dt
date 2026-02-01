package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/nwalker85/dt/internal/osascript"
)

var moveToDb string
var moveToGroup string

var moveCmd = &cobra.Command{
	Use:   "move <query>",
	Short: "Move items to another database or group",
	Long: `Move items matching a query to a different database or group.

Examples:
  dt move "tag:archive" --to "Archive"
  dt move "kind:pdf" --to "Work" --group "Reports"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]

		if moveToDb == "" {
			return fmt.Errorf("--to flag is required (target database name)")
		}

		var script string
		if moveToGroup != "" {
			script = fmt.Sprintf(`
tell application "DEVONthink"
	set targetDB to database "%s"
	if targetDB is missing value then
		error "Database not found: %s"
	end if
	set targetGroup to get record at "%s" in targetDB
	if targetGroup is missing value then
		set targetGroup to create location "%s" in targetDB
	end if
	set results to search "%s"
	set movedCount to 0
	repeat with r in results
		move record r to targetGroup
		set movedCount to movedCount + 1
	end repeat
	return movedCount
end tell`, moveToDb, moveToDb, moveToGroup, moveToGroup, strings.ReplaceAll(query, `"`, `\"`))
		} else {
			script = fmt.Sprintf(`
tell application "DEVONthink"
	set targetDB to database "%s"
	if targetDB is missing value then
		error "Database not found: %s"
	end if
	set results to search "%s"
	set movedCount to 0
	repeat with r in results
		move record r to root of targetDB
		set movedCount to movedCount + 1
	end repeat
	return movedCount
end tell`, moveToDb, moveToDb, strings.ReplaceAll(query, `"`, `\"`))
		}

		out, err := osascript.Run(script)
		if err != nil {
			return fmt.Errorf("move failed: %w", err)
		}

		dest := moveToDb
		if moveToGroup != "" {
			dest = moveToDb + "/" + moveToGroup
		}
		fmt.Printf("Moved %s items to %s\n", out, dest)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(moveCmd)
	moveCmd.Flags().StringVar(&moveToDb, "to", "", "Target database name (required)")
	moveCmd.Flags().StringVar(&moveToGroup, "group", "", "Target group/folder within the database")
	moveCmd.MarkFlagRequired("to")
}
