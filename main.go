package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/google/go-github/v57/github"
	"github.com/montanaflynn/stats"
	"github.com/olekukonko/tablewriter"
	"github.com/schollz/progressbar/v3"
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
		log.Fatal(err)
	}

	reviewStats.PrintTable(times)
}

type Stats struct {
	Owner  string
	Repo   string
	Client *github.Client
}

func (s *Stats) getPullRequests(ctx context.Context) ([]*github.PullRequest, error) {
	var allPulls []*github.PullRequest

	opt := &github.PullRequestListOptions{
		State:       "all",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	bar := progressbar.NewOptions(
		-1,
		progressbar.OptionSetDescription("Fetching pull requests..."),
	)
	for {
		bar.Add(1)

		pulls, resp, err := s.Client.PullRequests.List(ctx, s.Owner, s.Repo, opt)
		if err != nil {
			return nil, err
		}

		allPulls = append(allPulls, pulls...)
		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return allPulls, nil
}

func (s *Stats) getReviews(ctx context.Context, pull_number int) ([]*github.PullRequestReview, *github.Response, error) {
	return s.Client.PullRequests.ListReviews(ctx, s.Owner, s.Repo, pull_number, nil)
}

func (s *Stats) GetReviewTimes(ctx context.Context) (map[string][]float64, error) {
	reviewTimes := make(map[string][]float64)

	pulls, err := s.getPullRequests(ctx)
	if err != nil {
		return nil, err
	}

	bar := progressbar.NewOptions(
		len(pulls),
		progressbar.OptionSetDescription("Fetching reviews..."),
	)
	for _, pull := range pulls {
		bar.Add(1)

		reviews, _, err := s.getReviews(ctx, *pull.Number)
		if err != nil {
			return nil, err
		}

		for _, review := range reviews {
			times, ok := reviewTimes[*review.User.Login]
			if !ok {
				times = make([]float64, 0)
			}

			if review.SubmittedAt != nil {
				reviewTimes[*review.User.Login] = append(
					times,
					review.SubmittedAt.Sub(pull.CreatedAt.Time).Seconds(),
				)
			}
		}
	}

	return reviewTimes, nil
}

func (s *Stats) PrintTable(timing map[string][]float64) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "p25", "p50 (Median)", "p75", "Count"})

	for user, durations := range timing {
		data := stats.LoadRawData(durations)

		p25, _ := stats.Percentile(data, 25)
		p50, _ := stats.Percentile(data, 50)
		p75, _ := stats.Percentile(data, 75)

		table.Append([]string{
			user,
			(time.Duration(p25) * time.Second).String(),
			(time.Duration(p50) * time.Second).String(),
			(time.Duration(p75) * time.Second).String(),
			strconv.Itoa(len(data)),
		})
	}

	table.Render()
}
