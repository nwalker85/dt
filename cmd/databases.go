package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/nwalker85/dt/internal/osascript"
)

var databasesCmd = &cobra.Command{
	Use:   "databases",
	Short: "List database names",
	Long:  `List the names of all open DEVONthink databases.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		script := `
tell application "DEVONthink"
	set names to {}
	repeat with db in databases
		set end of names to name of db
	end repeat
	return names
end tell`

		out, err := osascript.Run(script)
		if err != nil {
			return fmt.Errorf("databases failed: %w", err)
		}

		if out == "" {
			return nil
		}

		names := strings.Split(out, ", ")

		if jsonOutput {
			jsonBytes, err := json.MarshalIndent(names, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(jsonBytes))
		} else {
			for _, name := range names {
				fmt.Println(name)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(databasesCmd)
}
