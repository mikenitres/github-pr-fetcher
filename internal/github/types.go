package github

import "time"

type PullRequest struct {
    ID       int       `json:"id"`
    Number   int       `json:"number"`
    Title    string    `json:"title"`
    URL      string    `json:"html_url"`
    State    string    `json:"state"`
    MergedAt time.Time `json:"merged_at"`
    User     User      `json:"user"`
}

type User struct {
    Login string `json:"login"`
    
}