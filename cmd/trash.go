package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/nwalker85/dt/internal/osascript"
)

var trashForce bool

var trashCmd = &cobra.Command{
	Use:   "trash <query>",
	Short: "Move items to trash",
	Long: `Move items matching a query to the trash.

Examples:
  dt trash "tag:delete"
  dt trash "name:temp*" --force`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]

		// First, count items
		countScript := fmt.Sprintf(`
tell application "DEVONthink"
	return count of (search "%s")
end tell`, strings.ReplaceAll(query, `"`, `\"`))

		countOut, err := osascript.Run(countScript)
		if err != nil {
			return fmt.Errorf("trash failed: %w", err)
		}

		if countOut == "0" {
			fmt.Println("No items match the query")
			return nil
		}

		if !trashForce {
			fmt.Printf("About to trash %s items. Use --force to confirm.\n", countOut)
			return nil
		}

		script := fmt.Sprintf(`
tell application "DEVONthink"
	set results to search "%s"
	set trashedCount to 0
	repeat with r in results
		delete record r
		set trashedCount to trashedCount + 1
	end repeat
	return trashedCount
end tell`, strings.ReplaceAll(query, `"`, `\"`))

		out, err := osascript.Run(script)
		if err != nil {
			return fmt.Errorf("trash failed: %w", err)
		}

		fmt.Printf("Trashed %s items\n", out)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(trashCmd)
	trashCmd.Flags().BoolVar(&trashForce, "force", false, "Actually delete items (required for safety)")
}
