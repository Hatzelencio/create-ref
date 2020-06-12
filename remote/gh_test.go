package remote

import (
	"bytes"
	"github.com/google/go-github/v32/github"
	"github.com/hatzelencio/crete-ref/utils/mocks"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
)

type envVariables struct {
	token      string
	repository string
	sha        string
	refs       string
}

func init() {
	cli = NewGithubClient(nil, &mocks.MockClient{})
}

func TestCreateGitRefSuccess(t *testing.T) {
	params := envVariables{
		token:      "secret-token",
		repository: "owner/owner-repo",
		sha:        "aae1a1fef...",
		refs:       "heads/branch,heads/branch-two",
	}
	setEnvVariables(params)

	mockGetRefNotFound()
	mockCreateRefOk()

	output := captureOutput(func() {
		CreateGitRef()
	})

	if !strings.Contains(output, "ref heads/branch-two was created") {
		t.Fatal("ref can't be created")
	}

	if !strings.Contains(output, "ref heads/branch was created") {
		t.Fatal("ref can't be created")
	}
}

func TestCreateGitRefAlreadyExists(t *testing.T) {
	setEnvVariables(envVariables{
		repository: "owner/owner-repo",
		refs:       "heads/branch",
	})

	mockGetRefFound()
	output := captureOutput(func() {
		CreateGitRef()
	})

	if !strings.Contains(output, "the reference heads/branch is already exists") {
		t.Fatal("ref can't be found")
	}
}

func captureOutput(f func()) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	f()
	log.SetOutput(os.Stderr)
	return buf.String()
}

func setEnvVariables(env envVariables) {
	_ = os.Setenv("GITHUB_TOKEN", env.token)
	_ = os.Setenv("GITHUB_REPOSITORY", env.repository)
	_ = os.Setenv("GITHUB_SHA", env.sha)
	_ = os.Setenv("INPUT_REFS", env.refs)
}

func mockCreateRefOk() {
	mocks.GetCreateRefFunc = func(ctx context.Context, owner string, repo string, ref *github.Reference) (*github.Reference, *github.Response, error) {
		res := github.Response{
			Response: &http.Response{StatusCode: 201},
		}
		return nil, &res, nil
	}
}

func mockGetRefNotFound() {
	mocks.GetGetRefFunc = func(ctx context.Context, owner string, repo string, ref string) (*github.Reference, *github.Response, error) {
		res := github.Response{
			Response: &http.Response{StatusCode: 404},
		}
		return nil, &res, nil
	}
}

func mockGetRefFound() {
	mocks.GetGetRefFunc = func(ctx context.Context, owner string, repo string, ref string) (*github.Reference, *github.Response, error) {
		res := github.Response{
			Response: &http.Response{StatusCode: 200},
		}
		return &github.Reference{}, &res, nil
	}
}
