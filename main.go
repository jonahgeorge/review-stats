package main

import (
	"context"
	"flag"
	"os"
	"time"

	"github.com/google/go-github/v57/github"
	"github.com/montanaflynn/stats"
	"github.com/olekukonko/tablewriter"
)

var (
	owner       string
	repo        string
	githubToken string
)

func init() {
	flag.StringVar(&owner, "owner", "", "")
	flag.StringVar(&repo, "repo", "", "")
	flag.StringVar(&githubToken, "github-token", "", "")
	flag.Parse()
}

func main() {
	client := github.NewClient(nil).WithAuthToken(githubToken)

	reviewStats := &Stats{
		Owner:  owner,
		Repo:   repo,
		Client: client,
	}

	times, err := reviewStats.GetReviewTimes(context.Background())
	if err != nil {
		panic(err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "p25", "p50 (Median)", "p75"})

	for user, durations := range times {
		data := stats.LoadRawData(durations)

		p25, _ := stats.Percentile(data, 25)
		p50, _ := stats.Percentile(data, 50)
		p75, _ := stats.Percentile(data, 75)

		table.Append([]string{
			user,
			(time.Duration(p25) * time.Second).String(),
			(time.Duration(p50) * time.Second).String(),
			(time.Duration(p75) * time.Second).String(),
		})
	}

	table.Render()
}

type Stats struct {
	Owner  string
	Repo   string
	Client *github.Client
}

func (s *Stats) GetPullRequests(ctx context.Context) ([]*github.PullRequest, *github.Response, error) {
	return s.Client.PullRequests.List(ctx, s.Owner, s.Repo, &github.PullRequestListOptions{
		ListOptions: github.ListOptions{PerPage: 1000},
	})
}

func (s *Stats) GetReviews(ctx context.Context, pull_number int) ([]*github.PullRequestReview, *github.Response, error) {
	return s.Client.PullRequests.ListReviews(ctx, s.Owner, s.Repo, pull_number, nil)
}

func (s *Stats) GetReviewTimes(ctx context.Context) (map[string][]float64, error) {
	reviewTimes := make(map[string][]float64)

	pulls, _, err := s.GetPullRequests(ctx)
	if err != nil {
		return nil, err
	}

	for _, pull := range pulls {
		reviews, _, err := s.GetReviews(ctx, *pull.Number)
		if err != nil {
			return nil, err
		}

		for _, review := range reviews {
			times, ok := reviewTimes[*review.User.Login]
			if !ok {
				times = make([]float64, 0)
			}
			reviewTimes[*review.User.Login] = append(times, review.SubmittedAt.Sub(pull.CreatedAt.Time).Seconds())
		}
	}

	return reviewTimes, nil
}
