package main

import (
	"fmt"
	"strings"
	"time"
)

// PREvent simulates a GitHub pull request webhook payload
type PREvent struct {
	Action string      `json:"action"`
	Number int         `json:"number"`
	PR     PullRequest `json:"pull_request"`
}

type PullRequest struct {
	Title          string    `json:"title"`
	Body           string    `json:"body"`
	State          string    `json:"state"`
	Draft          bool      `json:"draft"`
	Merged         bool      `json:"merged"`
	MergedAt       *time.Time `json:"merged_at"`
	Commits        int       `json:"commits"`
	Additions      int       `json:"additions"`
	Deletions      int       `json:"deletions"`
	ChangedFiles   int       `json:"changed_files"`
	Labels         []Label   `json:"labels"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Head           Branch    `json:"head"`
	Base           Branch    `json:"base"`
	Author         User      `json:"user"`
	Reviewers      []User    `json:"requested_reviewers"`
}

type Branch struct {
	Ref  string `json:"ref"`
	SHA  string `json:"sha"`
}

type User struct {
	Login string `json:"login"`
}

type Label struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

func handlePREvent(event PREvent) {
	pr := event.PR

	fmt.Printf("=== PR #%d [%s] ===\n", event.Number, strings.ToUpper(event.Action))
	fmt.Printf("  title:   %s\n", pr.Title)
	fmt.Printf("  author:  @%s\n", pr.Author.Login)
	fmt.Printf("  branch:  %s → %s\n", pr.Head.Ref, pr.Base.Ref)
	fmt.Printf("  state:   %s", pr.State)
	if pr.Draft {
		fmt.Print(" (draft)")
	}
	fmt.Println()

	if pr.Merged {
		fmt.Printf("  merged:  %s\n", pr.MergedAt.Format(time.RFC3339))
	}

	fmt.Printf("  diff:    +%d -%d across %d file(s) in %d commit(s)\n",
		pr.Additions, pr.Deletions, pr.ChangedFiles, pr.Commits)

	if len(pr.Labels) > 0 {
		names := make([]string, len(pr.Labels))
		for i, l := range pr.Labels {
			names[i] = l.Name
		}
		fmt.Printf("  labels:  %s\n", strings.Join(names, ", "))
	}

	if len(pr.Reviewers) > 0 {
		logins := make([]string, len(pr.Reviewers))
		for i, r := range pr.Reviewers {
			logins[i] = "@" + r.Login
		}
		fmt.Printf("  reviewers: %s\n", strings.Join(logins, ", "))
	}

	fmt.Printf("  created: %s\n", pr.CreatedAt.Format(time.RFC3339))
	fmt.Printf("  updated: %s\n", pr.UpdatedAt.Format(time.RFC3339))
}

func main() {
	now := time.Now()

	event := PREvent{
		Action: "opened",
		Number: 42,
		PR: PullRequest{
			Title:        "feat: add awesome feature",
			Body:         "This PR adds an awesome feature.",
			State:        "open",
			Draft:        false,
			Merged:       false,
			Commits:      3,
			Additions:    120,
			Deletions:    40,
			ChangedFiles: 5,
			Labels:       []Label{{Name: "enhancement", Color: "84b6eb"}, {Name: "reviewed", Color: "00ff00"}},
			CreatedAt:    now,
			UpdatedAt:    now,
			Head:         Branch{Ref: "feat/awesome", SHA: "abc1234"},
			Base:         Branch{Ref: "main", SHA: "def5678"},
			Author:       User{Login: "johndoe"},
			Reviewers:    []User{{Login: "janedoe"}, {Login: "bobsmith"}},
		},
	}

	handlePREvent(event)
}
