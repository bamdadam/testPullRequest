package main_test

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPREvent(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PREvent Suite")
}

var _ = Describe("PREvent", func() {

	Describe("action types", func() {
		It("should handle opened action", func() {
			event := PREvent{Action: "opened", Number: 1, PR: PullRequest{
				Title:  "feat: something",
				State:  "open",
				Author: User{Login: "johndoe"},
			}}
			Expect(event.Action).To(Equal("opened"))
			Expect(event.PR.State).To(Equal("open"))
		})

		It("should handle closed action", func() {
			event := PREvent{Action: "closed", Number: 2, PR: PullRequest{
				Title:  "fix: something",
				State:  "closed",
				Author: User{Login: "johndoe"},
			}}
			Expect(event.Action).To(Equal("closed"))
			Expect(event.PR.Merged).To(BeFalse())
		})

		It("should handle synchronize action", func() {
			event := PREvent{Action: "synchronize", Number: 3, PR: PullRequest{
				Title:  "chore: update",
				State:  "open",
				Author: User{Login: "johndoe"},
			}}
			Expect(event.Action).To(Equal("synchronize"))
		})
	})

	Describe("PR draft state", func() {
		It("should not be draft by default", func() {
			pr := PullRequest{Title: "feat: ready", State: "open"}
			Expect(pr.Draft).To(BeFalse())
		})

		It("should correctly flag a draft PR", func() {
			pr := PullRequest{Title: "wip: not ready", State: "open", Draft: true}
			Expect(pr.Draft).To(BeTrue())
		})
	})

	Describe("PR merge state", func() {
		Context("when PR is merged", func() {
			var (
				mergedAt time.Time
				pr       PullRequest
			)

			BeforeEach(func() {
				mergedAt = time.Now()
				pr = PullRequest{
					Title:    "feat: merged",
					State:    "closed",
					Merged:   true,
					MergedAt: &mergedAt,
				}
			})

			It("should be marked as merged", func() {
				Expect(pr.Merged).To(BeTrue())
			})

			It("should have a non-nil MergedAt", func() {
				Expect(pr.MergedAt).NotTo(BeNil())
			})

			It("should have closed state", func() {
				Expect(pr.State).To(Equal("closed"))
			})
		})

		Context("when PR is not merged", func() {
			It("should have nil MergedAt", func() {
				pr := PullRequest{Title: "feat: open", State: "open", Merged: false}
				Expect(pr.MergedAt).To(BeNil())
			})
		})
	})

	Describe("labels", func() {
		It("should support multiple labels", func() {
			pr := PullRequest{
				Labels: []Label{
					{Name: "bug", Color: "ff0000"},
					{Name: "urgent", Color: "ff6600"},
				},
			}
			Expect(pr.Labels).To(HaveLen(2))
			Expect(pr.Labels[0].Name).To(Equal("bug"))
			Expect(pr.Labels[1].Name).To(Equal("urgent"))
		})

		It("should support no labels", func() {
			pr := PullRequest{}
			Expect(pr.Labels).To(BeEmpty())
		})
	})

	Describe("reviewers", func() {
		It("should list requested reviewers", func() {
			pr := PullRequest{
				Reviewers: []User{
					{Login: "reviewer1"},
					{Login: "reviewer2"},
				},
			}
			Expect(pr.Reviewers).To(HaveLen(2))
			Expect(pr.Reviewers).To(ContainElement(User{Login: "reviewer1"}))
		})
	})

	Describe("diff stats", func() {
		It("should track additions and deletions", func() {
			pr := PullRequest{
				Commits:      4,
				Additions:    200,
				Deletions:    50,
				ChangedFiles: 8,
			}
			Expect(pr.Additions).To(BeNumerically(">", pr.Deletions))
			Expect(pr.ChangedFiles).To(Equal(8))
			Expect(pr.Commits).To(Equal(4))
		})
	})

	Describe("branch info", func() {
		It("should have head and base branches", func() {
			pr := PullRequest{
				Head: Branch{Ref: "feat/my-feature", SHA: "abc123"},
				Base: Branch{Ref: "main", SHA: "def456"},
			}
			Expect(pr.Head.Ref).To(Equal("feat/my-feature"))
			Expect(pr.Base.Ref).To(Equal("main"))
			Expect(pr.Head.SHA).To(HaveLen(6))
		})
	})
})
