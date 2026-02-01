package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/walkern/dt/internal/osascript"
)

var exportDest string

var exportCmd = &cobra.Command{
	Use:   "export <tag>",
	Short: "Export tagged items to filesystem",
	Long: `Export all items with a specific tag to a destination folder.

The destination defaults to the current directory if not specified.

Examples:
  dt export important
  dt export work --dest ~/exports/work`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		tag := args[0]

		dest := exportDest
		if dest == "" {
			var err error
			dest, err = os.Getwd()
			if err != nil {
				return err
			}
		}

		// Expand ~ in path
		if strings.HasPrefix(dest, "~") {
			home, _ := os.UserHomeDir()
			dest = filepath.Join(home, dest[1:])
		}

		// Ensure destination exists
		if err := os.MkdirAll(dest, 0755); err != nil {
			return fmt.Errorf("failed to create destination: %w", err)
		}

		absPath, err := filepath.Abs(dest)
		if err != nil {
			return err
		}

		script := fmt.Sprintf(`
tell application "DEVONthink"
	set results to search "tag:%s"
	set exportCount to 0
	repeat with r in results
		set filePath to path of r
		if filePath is not missing value then
			do shell script "cp " & quoted form of filePath & " " & quoted form of "%s/"
			set exportCount to exportCount + 1
		end if
	end repeat
	return exportCount
end tell`, strings.ReplaceAll(tag, `"`, `\"`), absPath)

		out, err := osascript.Run(script)
		if err != nil {
			return fmt.Errorf("export failed: %w", err)
		}

		fmt.Printf("Exported %s items to %s\n", out, absPath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().StringVar(&exportDest, "dest", "", "Destination directory for exported files")
}
