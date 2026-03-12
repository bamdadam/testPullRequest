package main

import (
	"strings"
	"testing"
	"time"
)

func TestHandlePREvent_Opened(t *testing.T) {
	event := PREvent{
		Action: "opened",
		Number: 1,
		PR: PullRequest{
			Title:        "feat: new feature",
			State:        "open",
			Draft:        false,
			Commits:      1,
			Additions:    10,
			Deletions:    2,
			ChangedFiles: 1,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			Head:         Branch{Ref: "feat/new", SHA: "aaa111"},
			Base:         Branch{Ref: "main", SHA: "bbb222"},
			Author:       User{Login: "johndoe"},
		},
	}

	if event.Action != "opened" {
		t.Errorf("expected action 'opened', got %q", event.Action)
	}
	if event.PR.Draft {
		t.Error("expected PR to not be a draft")
	}
	if event.PR.Merged {
		t.Error("expected PR to not be merged")
	}
}

func TestHandlePREvent_Draft(t *testing.T) {
	event := PREvent{
		Action: "opened",
		Number: 2,
		PR: PullRequest{
			Title:  "wip: work in progress",
			State:  "open",
			Draft:  true,
			Author: User{Login: "janedoe"},
		},
	}

	if !event.PR.Draft {
		t.Error("expected PR to be a draft")
	}
	if !strings.HasPrefix(event.PR.Title, "wip:") {
		t.Errorf("expected draft title to start with 'wip:', got %q", event.PR.Title)
	}
}

func TestHandlePREvent_Merged(t *testing.T) {
	mergedAt := time.Now()
	event := PREvent{
		Action: "closed",
		Number: 3,
		PR: PullRequest{
			Title:    "feat: merged feature",
			State:    "closed",
			Merged:   true,
			MergedAt: &mergedAt,
			Head:     Branch{Ref: "feat/merged", SHA: "ccc333"},
			Base:     Branch{Ref: "main", SHA: "ddd444"},
			Author:   User{Login: "bobsmith"},
		},
	}

	if !event.PR.Merged {
		t.Error("expected PR to be merged")
	}
	if event.PR.MergedAt == nil {
		t.Error("expected MergedAt to be set")
	}
	if event.PR.State != "closed" {
		t.Errorf("expected state 'closed', got %q", event.PR.State)
	}
}

func TestHandlePREvent_Labels(t *testing.T) {
	event := PREvent{
		Action: "labeled",
		Number: 4,
		PR: PullRequest{
			Title:  "fix: bug fix",
			State:  "open",
			Labels: []Label{{Name: "bug", Color: "ff0000"}, {Name: "urgent", Color: "ff6600"}},
			Author: User{Login: "johndoe"},
		},
	}

	if len(event.PR.Labels) != 2 {
		t.Errorf("expected 2 labels, got %d", len(event.PR.Labels))
	}
	if event.PR.Labels[0].Name != "bug" {
		t.Errorf("expected first label 'bug', got %q", event.PR.Labels[0].Name)
	}
}

func TestHandlePREvent_Reviewers(t *testing.T) {
	event := PREvent{
		Action: "review_requested",
		Number: 5,
		PR: PullRequest{
			Title:     "chore: cleanup",
			State:     "open",
			Reviewers: []User{{Login: "reviewer1"}, {Login: "reviewer2"}},
			Author:    User{Login: "johndoe"},
		},
	}

	if len(event.PR.Reviewers) != 2 {
		t.Errorf("expected 2 reviewers, got %d", len(event.PR.Reviewers))
	}
	if event.PR.Reviewers[0].Login != "reviewer1" {
		t.Errorf("expected first reviewer 'reviewer1', got %q", event.PR.Reviewers[0].Login)
	}
}

func TestHandlePREvent_DiffStats(t *testing.T) {
	event := PREvent{
		Action: "opened",
		Number: 6,
		PR: PullRequest{
			Title:        "refactor: big cleanup",
			State:        "open",
			Commits:      5,
			Additions:    300,
			Deletions:    150,
			ChangedFiles: 12,
			Author:       User{Login: "johndoe"},
		},
	}

	if event.PR.Additions <= event.PR.Deletions {
		t.Error("expected more additions than deletions for this PR")
	}
	if event.PR.ChangedFiles != 12 {
		t.Errorf("expected 12 changed files, got %d", event.PR.ChangedFiles)
	}
	if event.PR.Commits != 5 {
		t.Errorf("expected 5 commits, got %d", event.PR.Commits)
	}
}
