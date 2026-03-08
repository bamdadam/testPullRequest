
package main

import (
	"fmt"
	"time"
)

type PREvent struct {
	Action string
	Number int
	PR     PullRequest
}

type PullRequest struct {
	Title     string
	State     string
	CreatedAt time.Time
	Head      Branch
	Base      Branch
}

type Branch struct {
	Ref string
	SHA string
}

func handlePREvent(event PREvent) {
	fmt.Printf("[PR #%d] action=%q title=%q\n", event.Number, event.Action, event.PR.Title)
	fmt.Printf("  %s → %s\n", event.PR.Head.Ref, event.PR.Base.Ref)
}

func main() {
	event := PREvent{
		Action: "opened",
		Number: 42,
		PR: PullRequest{
			Title:     "feat: add awesome feature",
			State:     "open",
			CreatedAt: time.Now(),
			Head:      Branch{Ref: "feat/awesome", SHA: "abc1234"},
			Base:      Branch{Ref: "main", SHA: "def5678"},
		},
	}
	handlePREvent(event)
}
