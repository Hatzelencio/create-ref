package mocks

import (
	"github.com/google/go-github/v32/github"
	"golang.org/x/net/context"
)

// MockClient mock of github.Client
type MockClient struct {
	GetRefFunc func(ctx context.Context, owner string, repo string, ref string) (*github.Reference, *github.Response, error)
	CreateRefFunc func(ctx context.Context, owner string, repo string, ref *github.Reference) (*github.Reference, *github.Response, error)
}

var (
	// GetGetRefFunc mock variable to retrieve a git ref from github repository
	GetGetRefFunc func(ctx context.Context, owner string, repo string, ref string) (*github.Reference, *github.Response, error)
	// GetCreateRefFunc mock variable to create a git ref to github repository
	GetCreateRefFunc func(ctx context.Context, owner string, repo string, ref *github.Reference) (*github.Reference, *github.Response, error)
)

// GetRef mock function to retrieve a git ref from github repository
func (m *MockClient) GetRef(ctx context.Context, owner string, repo string, ref string) (*github.Reference, *github.Response, error) {
	return GetGetRefFunc(ctx, owner, repo, ref)
}

// CreateRef mock function to create a git ref to github repository
func (m *MockClient) CreateRef(ctx context.Context, owner string, repo string, ref *github.Reference) (*github.Reference, *github.Response, error){
	return GetCreateRefFunc(ctx, owner, repo, ref)
}