package main

import (
	"fmt"
	"os"

	"github.com/scottyw/outlier/outlier"
)

func main() {

	// Enable httpmock so that local test data is used rather than the real API
	// outlier.EnableHttpmock()

	if len(os.Args) != 3 {
		fmt.Println("Usage: go run cmd/cli/main.go <repo-owner> <repo-name> e.g. go run cmd/cli/main.go puppetlabs go-pe-client")
		os.Exit(1)
	}

	owner := os.Args[1]
	name := os.Args[2]
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		fmt.Println("GITHUB_TOKEN environment variable must be defined")
		os.Exit(1)
	}

	repoStats, err := outlier.GenerateStats(owner, name, token)
	if err != nil {
		panic(err)
	}

	fmt.Println("----------------------------------------------------------------")
	fmt.Printf("PR Stats for %s / %s\n", owner, name)
	fmt.Println("----------------------------------------------------------------")
	fmt.Printf("Total PRs                : %d\n", repoStats.Total)
	fmt.Printf("Duration (median)        : %d hours\n", repoStats.Median)
	fmt.Printf("Duration (75%%-ile)       : %d hours\n", repoStats.Percentile75)
	fmt.Printf("Duration (90%%-ile)       : %d hours\n", repoStats.Percentile90)
	fmt.Printf("Total Outliers (mild)    : %s\n", renderRange(repoStats.MildOutlierTotal, repoStats.MildOutlierMin, repoStats.ExtremeOutlierMin))
	fmt.Printf("Total Outliers (extreme) : %s\n", renderRange(repoStats.ExtremeOutlierTotal, repoStats.ExtremeOutlierMin, repoStats.ExtremeOutlierMin))
	fmt.Println("----------------------------------------------------------------")

}

func renderRange(total, min, max int) string {
	if total == 0 {
		return "None"
	}
	if min == max {
		return fmt.Sprintf("%d (over %d hours)", total, min)
	}
	return fmt.Sprintf("%d (between %d and %d hours)", total, min, max)
}
