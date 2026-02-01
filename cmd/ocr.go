package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/walkern/dt/internal/osascript"
)

var ocrCmd = &cobra.Command{
	Use:   "ocr <query>",
	Short: "OCR matching documents",
	Long: `Apply OCR to documents matching a query.

This converts scanned PDFs and images to searchable text.

Examples:
  dt ocr "kind:pdf tag:scanned"
  dt ocr "name:receipt*"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]

		script := fmt.Sprintf(`
tell application "DEVONthink"
	set results to search "%s"
	set ocrCount to 0
	repeat with r in results
		try
			ocr record r
			set ocrCount to ocrCount + 1
		end try
	end repeat
	return ocrCount
end tell`, strings.ReplaceAll(query, `"`, `\"`))

		out, err := osascript.Run(script)
		if err != nil {
			return fmt.Errorf("ocr failed: %w", err)
		}

		fmt.Printf("OCR applied to %s items\n", out)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(ocrCmd)
}
