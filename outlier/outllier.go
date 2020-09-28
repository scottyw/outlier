package outlier

import (
	"context"
	"fmt"
	"time"

	"github.com/montanaflynn/stats"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

// RepoStats contains statistics generated from the most recent 100 PRs specifially:
// * median i.e. the median duration of merged PRs
// * the 75th percentile i.e. the duration by which 75% of PRs had merged
// * the 90th percentile i.e. the duration by which 90% of PRs had merged
// * mild outlliers i.e. the duration of any PRs deemed to be somewhat longer than normal
// * extreme outlliers i.e. the duration of any PRs deemed to be much longer than normal
type RepoStats struct {
	Total               int
	Median              int
	Percentile75        int
	Percentile90        int
	MildOutlierTotal    int
	MildOutlierMin      int
	ExtremeOutlierTotal int
	ExtremeOutlierMin   int
}

// StatisticalOutliers summarises a list of outlier PRs
type StatisticalOutliers struct {
	Total int
	Min   int
}

// GenerateStats takes the owner and name of a GitHub repo along with a token that allows access and returns repo statistics
func GenerateStats(owner, name, token string) (*RepoStats, error) {
	prs, err := fetchGitHubPRs(owner, name, token)
	if err != nil {
		return nil, err
	}
	return calculateStats(prs)
}

func calculateStats(prs *githubPRs) (*RepoStats, error) {

	repoStats := &RepoStats{}
	hours := []float64{}
	for _, node := range prs.Repository.PullRequests.Nodes {
		duration := node.MergedAt.Sub(node.CreatedAt)
		duration = duration.Round(time.Hour)
		hours = append(hours, duration.Hours())
	}

	repoStats.Total = len(hours)

	median, err := stats.Median(hours)
	if err != nil {
		return nil, fmt.Errorf("error calculating median: %w", err)
	}
	repoStats.Median = int(median)

	percentile75, err := stats.Percentile(hours, 75)
	if err != nil {
		return nil, fmt.Errorf("error calculating 75th percentile: %w", err)
	}
	repoStats.Percentile75 = int(percentile75)

	percentile90, err := stats.Percentile(hours, 90)
	if err != nil {
		return nil, fmt.Errorf("error calculating 90th percentile: %w", err)
	}
	repoStats.Percentile90 = int(percentile90)

	outliers, err := stats.QuartileOutliers(hours)
	if err != nil {
		return nil, fmt.Errorf("error calculating quartile outliers: %w", err)
	}
	repoStats.MildOutlierTotal, repoStats.MildOutlierMin = findOutliers(outliers.Mild)
	repoStats.ExtremeOutlierTotal, repoStats.ExtremeOutlierMin = findOutliers(outliers.Extreme)

	return repoStats, nil

}

func findOutliers(hours []float64) (int, int) {
	if len(hours) == 0 {
		return 0, 0
	}
	min := int(hours[0])
	for _, hour := range hours {
		h := int(hour)
		if h < min {
			min = h
		}
	}
	return len(hours), min
}

func fetchGitHubPRs(owner, name, token string) (*githubPRs, error) {

	// Set up GraphQL client
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	client := githubv4.NewClient(httpClient)

	// Fetch PRs from GitHub
	githubPRs := &githubPRs{}
	variables := map[string]interface{}{
		"repositoryOwner":    githubv4.String(owner),
		"repositoryName":     githubv4.String(name),
		"pullRequestsStates": []githubv4.PullRequestState{githubv4.PullRequestStateMerged},
		"pullRequestsLast":   githubv4.Int(100),
	}
	err := client.Query(context.Background(), githubPRs, variables)
	if err != nil {
		return nil, err
	}
	return githubPRs, nil

}

type githubPRs struct {
	Repository struct {
		PullRequests struct {
			Nodes []struct {
				Number    int
				CreatedAt time.Time
				MergedAt  time.Time
				// Reviews   struct {
				// 	Nodes []struct {
				// 		CreatedAt time.Time
				// 		UpdatedAt time.Time
				// 		State     string
				// 	}
				// } `graphql:"reviews(last:$reviewsLast)"`
			}
		} `graphql:"pullRequests(states:$pullRequestsStates,last:$pullRequestsLast)"`
	} `graphql:"repository(owner:$repositoryOwner,name:$repositoryName)"`
}
