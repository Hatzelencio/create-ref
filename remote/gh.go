package remote

import (
	"fmt"
	"github.com/google/go-github/v32/github"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

const (
	token         = "GITHUB_TOKEN"
	ghrepo        = "GITHUB_REPOSITORY"
	ghsha         = "GITHUB_SHA"
	refs          = "INPUT_REFS"
	sha           = "INPUT_SHA"
	repository    = "INPUT_REPOSITORY"
	failRefExists = "INPUT_FAIL_IF_REF_EXISTS"
)

var (
	ctx context.Context
	cli GithubClient
)

func init() {
	ctx = context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv(token)},
	)
	tc := oauth2.NewClient(ctx, ts)
	cli = NewGithubClient(tc, nil)
}

// GithubGitService interface
type GithubGitService interface {
	GetRef(ctx context.Context, owner string, repo string, ref string) (*github.Reference, *github.Response, error)
	CreateRef(ctx context.Context, owner string, repo string, ref *github.Reference) (*github.Reference, *github.Response, error)
}

// GithubClient is a wrapper of github.Client
type GithubClient struct {
	Git GithubGitService
	*github.Client
}

// NewGithubClient Create a new github client
func NewGithubClient(client *http.Client, repoMock GithubGitService) GithubClient {
	if repoMock != nil {
		return GithubClient{
			Git: repoMock,
		}
	}

	cli := github.NewClient(client)
	return GithubClient{
		Git: cli.Git,
	}
}

// ValidateInputs validate if GITHUB_TOKEN and INPUT_REFS are present like environment variables
func ValidateInputs() error {
	if len(os.Getenv(token)) == 0 {
		return fmt.Errorf("%v is required env variable to trigger this action", token)
	}

	if len(os.Getenv(refs)) == 0 {
		return fmt.Errorf("%v is required input to trigger this action", "refs")
	}

	return nil
}

func getRefs() []string {
	return strings.Split(strings.Replace(
		os.Getenv(refs), " ", "", -1), ",")
}

func getOwnerRepo() (string, string) {
	var ownerRepo []string

	ownerRepo = strings.Split(strings.Replace(
		os.Getenv(repository), " ", "", -1), "/")

	if len(ownerRepo) != 2 {
		ownerRepo = strings.Split(os.Getenv(ghrepo), "/")
	}

	return ownerRepo[0], ownerRepo[1]
}

func getSHA() string {
	var commit string
	commit = strings.Replace(os.Getenv(sha), " ", "", -1)

	if len(commit) == 0 {
		commit = strings.Replace(os.Getenv(ghsha), " ", "", -1)
	}

	return commit
}

// CreateGitRef create a reference over github. It can be branch, tag
func CreateGitRef() {
	var wg sync.WaitGroup

	refs := getRefs()
	owner, repo := getOwnerRepo()
	sha := getSHA()

	wg.Add(len(refs))

	for _, ref := range refs {
		go func(client *GithubClient, owner, repo, sha, ref string) {
			var reference *github.Reference
			var err error

			defer wg.Done()

			reference, res, err := cli.Git.GetRef(ctx, owner, repo, ref)

			if res == nil {
				log.Fatal(err)
			}

			if err != nil && res.StatusCode != 404 {
				log.Fatal(err)
			}

			if reference != nil {
				message := fmt.Sprintf("the reference %v is already exists", ref)
				if os.Getenv(failRefExists) == "FORCE" {
					log.Fatal(message)
				}
				log.Println(message)
				return
			}

			gitRef := &github.Reference{
				Ref: &ref,
				Object: &github.GitObject{
					SHA: &sha,
				},
			}
			_, _, err = cli.Git.CreateRef(ctx, owner, repo, gitRef)

			if err != nil {
				log.Fatal(err)
			}
			log.Println(fmt.Sprintf("ref %v was created", ref))
		}(&cli, owner, repo, sha, ref)
	}
	wg.Wait()
}
