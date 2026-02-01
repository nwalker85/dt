package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/nwalker85/dt/internal/osascript"
)

var archiveTag string
var archiveGroup string

var archiveCmd = &cobra.Command{
	Use:   "archive <query>",
	Short: "Archive items matching a query",
	Long: `Archive items by adding an archive tag and optionally moving to an archive group.

Examples:
  dt archive "date:<lastmonth tag:processed"
  dt archive "tag:old" --group "Archive/2024"
  dt archive "kind:pdf date:<lastyear" --tag archived`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]

		var script string
		if archiveGroup != "" {
			script = fmt.Sprintf(`
tell application "DEVONthink"
	set results to search "%s"
	set archivedCount to 0
	repeat with r in results
		-- Add archive tag
		set currentTags to tags of r
		if currentTags does not contain "%s" then
			set end of currentTags to "%s"
			set tags of r to currentTags
		end if

		-- Move to archive group (create if needed)
		set parentDB to database of r
		set archiveLocation to get record at "%s" in parentDB
		if archiveLocation is missing value then
			set archiveLocation to create location "%s" in parentDB
		end if
		move record r to archiveLocation
		set archivedCount to archivedCount + 1
	end repeat
	return archivedCount
end tell`, strings.ReplaceAll(query, `"`, `\"`), archiveTag, archiveTag, archiveGroup, archiveGroup)
		} else {
			script = fmt.Sprintf(`
tell application "DEVONthink"
	set results to search "%s"
	set archivedCount to 0
	repeat with r in results
		-- Add archive tag
		set currentTags to tags of r
		if currentTags does not contain "%s" then
			set end of currentTags to "%s"
			set tags of r to currentTags
		end if
		set archivedCount to archivedCount + 1
	end repeat
	return archivedCount
end tell`, strings.ReplaceAll(query, `"`, `\"`), archiveTag, archiveTag)
		}

		out, err := osascript.Run(script)
		if err != nil {
			return fmt.Errorf("archive failed: %w", err)
		}

		if archiveGroup != "" {
			fmt.Printf("Archived %s items (tagged '%s', moved to '%s')\n", out, archiveTag, archiveGroup)
		} else {
			fmt.Printf("Archived %s items (tagged '%s')\n", out, archiveTag)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(archiveCmd)
	archiveCmd.Flags().StringVar(&archiveTag, "tag", "archived", "Tag to apply to archived items")
	archiveCmd.Flags().StringVar(&archiveGroup, "group", "", "Group/folder to move archived items to")
}
