package git

import (
	"context"
	"os"
	"os/exec"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCurrentBranch(t *testing.T) {
	result, err := CurrentBranch(context.Background())
	// TODO find a way to test this that does not rely on the current state of
	// the git repository.
	_ = err
	_ = result
}

func TestCurrentBranchSlash(t *testing.T) {
	dir := t.TempDir()
	cmd := exec.Command("git", "init")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init: %v\n%s", err, out)
	}
	cmd = exec.Command("git", "commit", "--allow-empty", "-m", "init")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git commit: %v\n%s", err, out)
	}
	cmd = exec.Command("git", "checkout", "-b", "feature/my-fix")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git checkout: %v\n%s", err, out)
	}

	// CurrentBranch reads from the cwd, so we need to chdir
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(orig) })

	branch, err := CurrentBranch(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if branch != "feature/my-fix" {
		t.Errorf("CurrentBranch: got %q, want %q", branch, "feature/my-fix")
	}
}

func TestRoot(t *testing.T) {
	result, err := Root("")
	// TODO figure out a way to test this as well - it depends on your current
	// working directory, and you can run "go test" from anywhere on the
	// filesystem.
	_ = err
	_ = result
}

func TestVersion(t *testing.T) {
	if len(version) < 1 {
		t.Errorf("git version is empty, expected a non-empty string")
	}
}

var remoteTests = []struct {
	remote   string
	expected RemoteURL
}{
	{
		"git@github.com:Kevinburke/shyp_api.git", RemoteURL{
			Host:     "github.com",
			Port:     22,
			Path:     "Kevinburke",
			RepoName: "shyp_api",
			Format:   SSHFormat,
			URL:      "git@github.com:Kevinburke/shyp_api.git",
			SSHUser:  "git",
		},
	}, {
		"git@github.com:Kevinburke/shyp_api", RemoteURL{
			Host:     "github.com",
			Port:     22,
			Path:     "Kevinburke",
			RepoName: "shyp_api",
			Format:   SSHFormat,
			URL:      "git@github.com:Kevinburke/shyp_api",
			SSHUser:  "git",
		},
	}, {
		"git@github.com:Kevinburke/shyp_api.git/", RemoteURL{
			Path:     "Kevinburke",
			Host:     "github.com",
			Port:     22,
			RepoName: "shyp_api",
			Format:   SSHFormat,
			URL:      "git@github.com:Kevinburke/shyp_api.git/",
			SSHUser:  "git",
		},
	}, {
		"git@github.com:path/to/Kevinburke/shyp_api.git/", RemoteURL{
			Path:     "path/to/Kevinburke",
			Host:     "github.com",
			Port:     22,
			RepoName: "shyp_api",
			Format:   SSHFormat,
			URL:      "git@github.com:path/to/Kevinburke/shyp_api.git/",
			SSHUser:  "git",
		},
	}, {
		"https://github.com/Kevinburke/shyp_api.git", RemoteURL{
			Path:     "Kevinburke",
			Host:     "github.com",
			Port:     443,
			RepoName: "shyp_api",
			Format:   HTTPSFormat,
			URL:      "https://github.com/Kevinburke/shyp_api.git",
			SSHUser:  "",
		},
	}, {
		"https://github.com/Kevinburke/shyp_api.git/", RemoteURL{
			Path:     "Kevinburke",
			Host:     "github.com",
			Port:     443,
			RepoName: "shyp_api",
			Format:   HTTPSFormat,
			URL:      "https://github.com/Kevinburke/shyp_api.git/",
			SSHUser:  "",
		},
	}, {
		"https://github.com:11443/Kevinburke/shyp_api.git", RemoteURL{
			Path:     "Kevinburke",
			Host:     "github.com",
			Port:     11443,
			RepoName: "shyp_api",
			Format:   HTTPSFormat,
			URL:      "https://github.com:11443/Kevinburke/shyp_api.git",
			SSHUser:  "",
		},
	}, {
		"https://github.com/Kevinburke/shyp_api", RemoteURL{
			Path:     "Kevinburke",
			Host:     "github.com",
			Port:     443,
			RepoName: "shyp_api",
			Format:   HTTPSFormat,
			URL:      "https://github.com/Kevinburke/shyp_api",
			SSHUser:  "",
		},
	}, {
		"https://github.com/Kevinburke/repo.name.with.periods.git", RemoteURL{
			Path:     "Kevinburke",
			Host:     "github.com",
			Port:     443,
			RepoName: "repo.name.with.periods",
			Format:   HTTPSFormat,
			URL:      "https://github.com/Kevinburke/repo.name.with.periods.git",
			SSHUser:  "",
		},
	}, {
		"ssh://git@github.com/Kevinburke/sshrepo.git", RemoteURL{
			Path:     "Kevinburke",
			Host:     "github.com",
			Port:     22,
			RepoName: "sshrepo",
			Format:   SSHFormat,
			URL:      "ssh://git@github.com/Kevinburke/sshrepo.git",
			SSHUser:  "git",
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
