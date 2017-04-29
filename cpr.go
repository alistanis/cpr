package cpr

import (
	"context"
	"errors"
	"io/ioutil"
	"os"

	"strings"

	"encoding/json"

	"syscall"

	"flag"

	"time"

	"github.com/google/go-github/github"
	"golang.org/x/crypto/ssh/terminal"
	git "gopkg.in/src-d/go-git.v4"
)

var (
	ErrNoGitParent            = errors.New("Could not find a parent git directory.")
	ErrNoGithubRemote         = errors.New("No github remote could be found for the current git repository.")
	ErrMalformedRepositoryUrl = errors.New("Could not parse the repository url.")
)

// Open recursively searches backwards for a git repository root. If it finds one, it returns that repository object.
// If it does not, returns the error ErrNotGitParent
func Open(path string) (*git.Repository, error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		err = os.Chdir("..")
		if err != nil {
			return nil, err
		}
		d, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		if d == "/" {
			return nil, ErrNoGitParent
		}
		return Open(d)
	}
	return r, nil
}

// GithubURL returns the first github URL found in a remote, or an error if found
func GithubURL(r *git.Repository) (string, error) {
	c, err := r.Config()
	if err != nil {
		return "", err
	}
	// We will return the first remote found that includes the string github
	for _, v := range c.Remotes {
		if strings.Contains(v.URL, "github") {
			// I am slightly not ok with this method chaining but it's the best way to accomplish it
			return strings.Replace(
					strings.TrimSuffix(
						strings.TrimPrefix(
							v.URL, "git@"),
						".git"),
					":", "/", -1),
				nil
		}
	}
	return "", ErrNoGithubRemote
}

type RepoInfo struct {
	Owner      string
	Repository string
}

func GetRepoInfo(s string) (*RepoInfo, error) {
	if !strings.Contains(s, "github") {
		return nil, ErrNoGithubRemote
	}

	s = strings.TrimPrefix(s, "github.com/")

	r := &RepoInfo{}
	split := strings.Split(s, "/")
	if len(split) < 2 {
		return nil, ErrMalformedRepositoryUrl
	}
	r.Owner = split[0]
	r.Repository = split[1]
	return r, nil
}

type Config struct {
	User     string
	Password string
}

func DefaultConfig() (*Config, error) {
	return ReadConfigFromFile("cpr.json")
}

func (c *Config) Save(path string) error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 0644)
}

func ReadConfigFromFile(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	c := &Config{}
	return c, json.Unmarshal(data, c)
}

func (o *Options) Transport() *github.BasicAuthTransport {
	return &github.BasicAuthTransport{Username: o.UserName, Password: o.Password}
}

func GetPasswd() (string, error) {
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	password := string(bytePassword)
	return strings.TrimSpace(password), nil
}

type Options struct {
	BaseBranch    string
	CompareBranch string
	Reviewers     []string
	Assignees     []string
	Comment       string
	UserName      string
	Password      string
	ConfigFile    string
	Title         string
	Body          string
	FlagSet       *flag.FlagSet
}

var (
	ErrNoBaseBranch    = errors.New("No base-branch was given, base-branch is required")
	ErrNoCompareBranch = errors.New("No compare-branch was given, compare-branch is required")
)

func (o *Options) Validate() error {
	if o.BaseBranch == "" {
		return ErrNoBaseBranch
	}
	if o.CompareBranch == "" {
		return ErrNoCompareBranch
	}
	return nil
}

func (o *Options) PullRequest(url string) (*github.PullRequest, *github.Response, error) {
	transport := o.Transport()
	client := github.NewClient(transport.Client())

	service := client.PullRequests
	ctx, cancel := context.WithTimeout(nil, time.Second*15)
	defer cancel()
	pr := &github.NewPullRequest{}
	pr.Base = &o.BaseBranch
	pr.Head = &o.CompareBranch
	pr.Title = &o.Title
	*pr.MaintainerCanModify = true
	if o.Body != "" {
		pr.Body = &o.Body
	}
	info, err := GetRepoInfo(url)
	if err != nil {
		return nil, nil, err
	}
	return service.Create(ctx, info.Owner, info.Repository, pr)
}

func ParseOptions(f *flag.FlagSet, args []string) (*Options, error) {
	o := &Options{}
	o.FlagSet = f
	f.StringVar(&o.BaseBranch, "base-branch", "", "The base branch to merge into (master|develop|release|staging) (Required)")
	f.StringVar(&o.CompareBranch, "compare-branch", "", "The branch you are attempting to merge (feature|bugfix) (Required)")

	var reviewersString string
	f.StringVar(&reviewersString, "reviewers", "", "A comma separated list of reviewers (Chris,Paul) (Optional)")
	f.StringVar(&reviewersString, "r", "", "A comma separated list of reviewers (Chris,Paul) (Optional)")

	var assigneesString string
	f.StringVar(&assigneesString, "assignees", "", "A comma separated list of assignees (Chris,Dan) (Optional)")
	f.StringVar(&assigneesString, "a", "", "A comma separated list of assignees (Chris,Dan) (Optional)")

	f.StringVar(&o.UserName, "user", "", "Github username (alistanis) (Optional)")
	f.StringVar(&o.Password, "pass", "", "Github password (asckoq14rf0n!@$) (Optional)")
	f.StringVar(&o.ConfigFile, "config", "cpr.json", "Config file location for this repository (Optional)")

	f.StringVar(&o.Title, "title", "", "The title of this pull request (Required)")
	f.StringVar(&o.Body, "body", "", "The description of this pull request (Optional)")

	err := f.Parse(args)
	if err != nil {
		return nil, err
	}

	o.Reviewers = strings.Split(reviewersString, ",")
	o.Assignees = strings.Split(assigneesString, ",")
	return o, nil
}
