package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/walkern/dt/internal/osascript"
)

var weeklyReportCmd = &cobra.Command{
	Use:   "weekly-report",
	Short: "Generate a summary of items added this week",
	Long:  `Generate a report showing items added to DEVONthink in the past week, grouped by database.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		script := `
var app = Application("DEVONthink");
var dbs = app.databases();
var cutoff = new Date();
cutoff.setDate(cutoff.getDate() - 7);

var report = {
	generatedAt: new Date().toISOString(),
	period: "last 7 days",
	databases: []
};

var totalItems = 0;

for (var i = 0; i < dbs.length; i++) {
	var db = dbs[i];
	var contents = db.contents();
	var dbReport = {
		name: db.name(),
		newItems: [],
		count: 0
	};

	for (var j = 0; j < contents.length; j++) {
		var item = contents[j];
		var addDate = item.additionDate();
		if (addDate && addDate > cutoff) {
			dbReport.newItems.push({
				name: item.name(),
				kind: item.kind(),
				additionDate: addDate.toISOString(),
				tags: item.tags()
			});
			dbReport.count++;
			totalItems++;
		}
	}

	if (dbReport.count > 0) {
		report.databases.push(dbReport);
	}
}

report.totalNewItems = totalItems;

JSON.stringify(report, null, 2);`

		out, err := osascript.RunJS(script)
		if err != nil {
			return fmt.Errorf("weekly-report failed: %w", err)
		}

		fmt.Println(out)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(weeklyReportCmd)
}
