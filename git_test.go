package git

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCurrentBranch(t *testing.T) {
	t.Skip("CI checks out a specific commit, not on a branch")
	result, err := CurrentBranch()
	// TODO find a way to test this that does not rely on the current state of
	// the git repository.
	if err != nil {
		t.Fatal(err)
	}
	if result != "go-git" {
		t.Errorf("wrong branch name, got %s", result)
	}
}

func TestRoot(t *testing.T) {
	result, err := Root("")
	if err != nil {
		t.Fatal(err)
	}
	_, end := filepath.Split(result)
	if end != "go-git" {
		t.Errorf("wrong directory path: %s", end)
	}
}

var remoteTests = []struct {
	remote   string
	expected RemoteURL
}{
	{
		"git@github.com:Shyp/shyp_api.git", RemoteURL{
			Host:     "github.com",
			Port:     22,
			Path:     "Shyp",
			RepoName: "shyp_api",
			Format:   SSHFormat,
			URL:      "git@github.com:Shyp/shyp_api.git",
			SSHUser:  "git",
		},
	}, {
		"git@github.com:Shyp/shyp_api", RemoteURL{
			Host:     "github.com",
			Port:     22,
			Path:     "Shyp",
			RepoName: "shyp_api",
			Format:   SSHFormat,
			URL:      "git@github.com:Shyp/shyp_api",
			SSHUser:  "git",
		},
	}, {
		"git@github.com:Shyp/shyp_api.git/", RemoteURL{
			Path:     "Shyp",
			Host:     "github.com",
			Port:     22,
			RepoName: "shyp_api",
			Format:   SSHFormat,
			URL:      "git@github.com:Shyp/shyp_api.git/",
			SSHUser:  "git",
		},
	}, {
		"git@github.com:path/to/Shyp/shyp_api.git/", RemoteURL{
			Path:     "path/to/Shyp",
			Host:     "github.com",
			Port:     22,
			RepoName: "shyp_api",
			Format:   SSHFormat,
			URL:      "git@github.com:path/to/Shyp/shyp_api.git/",
			SSHUser:  "git",
		},
	}, {
		"https://github.com/Shyp/shyp_api.git", RemoteURL{
			Path:     "Shyp",
			Host:     "github.com",
			Port:     443,
			RepoName: "shyp_api",
			Format:   HTTPSFormat,
			URL:      "https://github.com/Shyp/shyp_api.git",
			SSHUser:  "",
		},
	}, {
		"https://github.com/Shyp/shyp_api.git/", RemoteURL{
			Path:     "Shyp",
			Host:     "github.com",
			Port:     443,
			RepoName: "shyp_api",
			Format:   HTTPSFormat,
			URL:      "https://github.com/Shyp/shyp_api.git/",
			SSHUser:  "",
		},
	}, {
		"https://github.com:11443/Shyp/shyp_api.git", RemoteURL{
			Path:     "Shyp",
			Host:     "github.com",
			Port:     11443,
			RepoName: "shyp_api",
			Format:   HTTPSFormat,
			URL:      "https://github.com:11443/Shyp/shyp_api.git",
			SSHUser:  "",
		},
	}, {
		"https://github.com/Shyp/shyp_api", RemoteURL{
			Path:     "Shyp",
			Host:     "github.com",
			Port:     443,
			RepoName: "shyp_api",
			Format:   HTTPSFormat,
			URL:      "https://github.com/Shyp/shyp_api",
			SSHUser:  "",
		},
	}, {
		"https://github.com/Shyp/repo.name.with.periods.git", RemoteURL{
			Path:     "Shyp",
			Host:     "github.com",
			Port:     443,
			RepoName: "repo.name.with.periods",
			Format:   HTTPSFormat,
			URL:      "https://github.com/Shyp/repo.name.with.periods.git",
			SSHUser:  "",
		},
	},
}

func TestParseRemoteURL(t *testing.T) {
	for _, tt := range remoteTests {
		remote, err := ParseRemoteURL(tt.remote)
		if err != nil {
			t.Fatal(err)
		}
		if remote == nil {
			t.Fatalf("expected ParseRemoteURL(%s) to be %v, was nil", tt.remote, tt.expected)
		}
		if diff := cmp.Diff(*remote, tt.expected); diff != "" {
			t.Errorf("ParseRemoteURL(%q): (-got +want)\n%s", tt.remote, diff)
		}
	}
}

func TestGetRemoteURL(t *testing.T) {
	r, err := GetRemoteURL("origin")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(r.RepoName)
}

func TestTip(t *testing.T) {
	s, err := Tip("")
	fmt.Println("s", s)
	fmt.Println("err", err)
	s, err = Tip("master")
	fmt.Println("s", s)
	fmt.Println("err", err)
	//t.Fail()
}
