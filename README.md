# Review Stats

Calculates Pull Request review stats for a given Github repository.

## Quick Start

```sh
go install github.com/jonahgeorge/review-stats
```

```sh
review-stats -owner jonahgeorge -repo review-stats -github-token $GITHUB_TOKEN

+----------------------+------------+--------------+------------+
|         NAME         |    P25     | P50 (MEDIAN) |    P75     |
+----------------------+------------+--------------+------------+
| jonahgeorge          | 1h56m29s   | 15h56m29s    | 15h56m29s  |
```
