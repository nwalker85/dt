package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/walkern/dt/internal/osascript"
)

var importDB string

var importCmd = &cobra.Command{
	Use:   "import <source>",
	Short: "Import folder into DEVONthink database",
	Long: `Import a folder and its contents into a DEVONthink database.

If no database is specified, imports into the first open database.

Examples:
  dt import ~/Documents/reports
  dt import ~/Downloads/papers --db "Research"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		source := args[0]

		// Expand ~ in path
		if strings.HasPrefix(source, "~") {
			home, _ := os.UserHomeDir()
			source = filepath.Join(home, source[1:])
		}

		absPath, err := filepath.Abs(source)
		if err != nil {
			return err
		}

		// Check source exists
		if _, err := os.Stat(absPath); os.IsNotExist(err) {
			return fmt.Errorf("source path does not exist: %s", absPath)
		}

		var script string
		if importDB != "" {
			script = fmt.Sprintf(`
tell application "DEVONthink"
	set targetDB to database "%s"
	if targetDB is missing value then
		error "Database not found: %s"
	end if
	set importedItems to import POSIX path "%s" to root of targetDB
	return count of importedItems
end tell`, importDB, importDB, absPath)
		} else {
			script = fmt.Sprintf(`
tell application "DEVONthink"
	set targetDB to first database
	set importedItems to import POSIX path "%s" to root of targetDB
	return count of importedItems
end tell`, absPath)
		}

		out, err := osascript.Run(script)
		if err != nil {
			return fmt.Errorf("import failed: %w", err)
		}

		dbName := importDB
		if dbName == "" {
			dbName = "default database"
		}
		fmt.Printf("Imported %s items from %s into %s\n", out, absPath, dbName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.Flags().StringVar(&importDB, "db", "", "Target database name")
}
