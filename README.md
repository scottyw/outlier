# Outlier: Find outlier PRs in a GitHub repo

**Note: This is example code for pulling PRs from the GitHub API. I can't recommend using the output to make real decisions.**

This is a proof-of-concept that will fetch the last 100 pull requests from a GitHub repo then find outliers.

Statistics are based on the time from creation to merge:
* Median
* 75th percentile (i.e. 75% of PRs have merged after this time)
* 90th percentile (i.e. 90% of PRs have merged after this time)
* Mild and extreme outliers

## Caveats

The difference between creation time and merge time as choice of metric here is arbitrary and may not have meaning for you, depending on context. In particular, if you use GitHub's brilliant "draft" PR functionality to show work-in progress - which is a great way to work collaboratively - then this tool is pulling meaningless metrics.

On the philosopical side, pull requests do not reflect business value. You cannot judge the performance of an individual, or even really a team, based on data pulled from GitHub. A team that wants to improve its own performance might benefit from aggregate PR data and this code is nothing more than a skeletal framework to get that exercise started.

## Running the tool

You need to run from source like this:

    go run cmd/cli/main.go <repo-owner> <repo-name>

For example:

    go run cmd/cli/main.go puppetlabs go-pe-client

You will need to specify an environment variable called `GITHUB_TOKEN` containing an API token that allows access to GitHub's GraphQL API:

    export GITHUB_TOKEN=xxxx

If you want to try it out without providing an API token then uncomment this line in `main.go` to use the included mock data:

	// outlier.EnableHttpmock()

## Example Output
```
~/github/outlier (master ✘)✖✹ ᐅ go run cmd/cli/main.go puppetlabs go-pe-client
----------------------------------------------------------------
PR Stats for puppetlabs / go-pe-client
----------------------------------------------------------------
Total PRs                : 48
Duration (median)        : 2 hours
Duration (75%-ile)       : 41 hours
Duration (90%-ile)       : 156 hours
Total Outliers (mild)    : 2 (between 145 and 191 hours)
Total Outliers (extreme) : 4 (over 191 hours)
----------------------------------------------------------------
```
