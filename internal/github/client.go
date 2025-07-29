package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	token      string
	httpClient *http.Client
}

// RepoResult stores the results for a single repository
type RepoResult struct {
	Repository string
	PRs        []PullRequest
	Error      error
}

func NewClient(token string) *Client {
	return &Client{
		token: token,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// FetchMergedPRsForRepos fetches merged PRs for multiple repositories since a given date
func (c *Client) FetchMergedPRsForRepos(repos []string, branch string, since time.Time) []RepoResult {
	results := make([]RepoResult, 0, len(repos))

	for _, repo := range repos {
		prs, err := c.FetchMergedPRsSince(repo, branch, since)
		results = append(results, RepoResult{
			Repository: repo,
			PRs:        prs,
			Error:      err,
		})
	}

	return results
}

func (c *Client) FetchMergedPRsSince(repo, branch string, since time.Time) ([]PullRequest, error) {
	var mergedPRs []PullRequest
	nextPageURL := fmt.Sprintf("https://api.github.com/repos/%s/pulls?state=closed&base=%s&per_page=100&sort=updated&direction=desc", repo, branch)

	// Continue fetching until we find a PR older than our since date
	for nextPageURL != "" {
		fmt.Printf("Fetching: %s\n", nextPageURL)

		req, err := http.NewRequest("GET", nextPageURL, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Authorization", "Bearer "+c.token)
		req.Header.Set("Accept", "application/vnd.github.v3+json")
		req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("GitHub API returned status: %d", resp.StatusCode)
		}

		var pagePRs []PullRequest
		if err := json.NewDecoder(resp.Body).Decode(&pagePRs); err != nil {
			resp.Body.Close()
			return nil, err
		}

		// Extract the Link header for pagination
		nextPageURL = ""
		if linkHeader := resp.Header.Get("Link"); linkHeader != "" {
			for _, link := range strings.Split(linkHeader, ",") {
				parts := strings.Split(strings.TrimSpace(link), ";")
				if len(parts) != 2 {
					continue
				}

				// Check if this is the 'next' link
				if strings.Contains(parts[1], `rel="next"`) {
					// Extract URL from the angle brackets
					url := strings.TrimSpace(parts[0])
					nextPageURL = strings.Trim(url, "<>")
					break
				}
			}
		}

		resp.Body.Close()

		// Process PRs and check if they're within our time range
		foundOlder := false
		for _, pr := range pagePRs {
			if !pr.MergedAt.IsZero() {
				if pr.MergedAt.After(since) || pr.MergedAt.Equal(since) {
					mergedPRs = append(mergedPRs, pr)
				} else {
					// Found a PR older than our since date, we can stop after this page
					foundOlder = true
				}
			}
		}

		// If we found PRs older than our date, stop pagination
		if foundOlder {
			nextPageURL = ""
		}

		// If we got no results on this page and there's no next page, break
		if len(pagePRs) == 0 && nextPageURL == "" {
			break
		}
	}

	return mergedPRs, nil
}
