package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/mikenitres/github-pr-fetcher/internal/config"
	"github.com/mikenitres/github-pr-fetcher/internal/github"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "github-pr-fetcher",
		Short: "Fetch merged pull requests from GitHub repositories",
		Long:  "Fetch merged pull requests from GitHub repositories since a specific date or number of weeks ago",
		Run:   run,
	}

	rootCmd.Flags().StringSliceP("repos", "r", []string{}, "GitHub repositories (comma-separated, e.g., owner/repo1,owner/repo2)")
	rootCmd.Flags().StringP("branch", "b", "main", "Base branch to check PRs against")
	rootCmd.Flags().StringP("config", "c", "config.json", "Path to config file")
	rootCmd.Flags().StringP("since", "s", "", "Date since when to fetch PRs (YYYY-MM-DD format)")
	rootCmd.Flags().IntP("weeks", "w", 0, "Number of weeks ago to fetch PRs from (overrides --since)")
	rootCmd.MarkFlagRequired("repos")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) {
	repos, _ := cmd.Flags().GetStringSlice("repos")
	branch, _ := cmd.Flags().GetString("branch")
	configPath, _ := cmd.Flags().GetString("config")
	sinceStr, _ := cmd.Flags().GetString("since")
	weeks, _ := cmd.Flags().GetInt("weeks")

	// Determine the since date
	var since time.Time
	if weeks > 0 {
		since = time.Now().AddDate(0, 0, -7*weeks)
	} else if sinceStr != "" {
		parsedDate, err := time.Parse("2006-01-02", sinceStr)
		if err != nil {
			log.Fatalf("Invalid date format. Use YYYY-MM-DD: %v", err)
		}
		since = parsedDate
	} else {
		// Default to 1 week ago
		since = time.Now().AddDate(0, 0, -7)
	}

	// Load configuration
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if cfg.GitHubToken == "" {
		log.Fatalf("GitHub token is missing in configuration")
	}

	// Create GitHub client
	client := github.NewClient(cfg.GitHubToken)

	fmt.Printf("Fetching PRs merged since: %s\n\n", since.Format("2006-01-02 15:04:05"))

	// Fetch merged PRs for all repositories
	results := client.FetchMergedPRsForRepos(repos, branch, since)

	totalPRs := 0

	// Display results for each repository
	for _, result := range results {
		if result.Error != nil {
			fmt.Printf("❌ Error fetching PRs for %s: %v\n\n", result.Repository, result.Error)
			continue
		}

		fmt.Printf("📂 Repository: %s\n", result.Repository)
		fmt.Printf("Found %d merged PRs since %s\n\n", len(result.PRs), since.Format("2006-01-02"))

		totalPRs += len(result.PRs)

		for i, pr := range result.PRs {
			fmt.Printf("%d. PR #%d: %s\n", i+1, pr.Number, pr.Title)
			fmt.Printf("   URL: %s\n", pr.URL)
			fmt.Printf("   Merged at: %s\n", pr.MergedAt.Format("2006-01-02 15:04:05"))
			fmt.Println()
		}

		fmt.Println(strings.Repeat("-", 50))
	}

	fmt.Printf("\n📊 Total PRs across all repositories: %d\n", totalPRs)
}
