# Review Stats

Calculates Pull Request review stats for a given Github repository.

## Quick Start

```sh
go install github.com/jonahgeorge/review-stats
```

For a repository with ~2,000 PRs, the script takes roughly 10 minutes to complete.

```sh
review-stats -owner jonahgeorge -repo review-stats -github-token $GITHUB_TOKEN

+----------------------+------------+--------------+------------+
|         NAME         |    P25     | P50 (MEDIAN) |    P75     |
+----------------------+------------+--------------+------------+
| jonahgeorge          | 1h56m29s   | 15h56m29s    | 15h56m29s  |
```
