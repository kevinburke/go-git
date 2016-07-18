package pr

import (
	"os/exec"
	"strings"

	"github.com/github/hub/commands"
)

func getLastCommitMessage() (string, error) {
	result, err := exec.Command("git", "log", "-1", "--pretty=%B").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(result)), nil
}

// CreateAndOpenPullRequest creates a PR on Github and opens it up in your
// browser. Returns an error if you are not on a branch, or if you are not in a
// git repository.
func CreateAndOpenPullRequest() error {
	message, err := getLastCommitMessage()
	if err != nil {
		return err
	}
	cmd := commands.CmdRunner.Lookup("pull-request")
	args := commands.NewArgs([]string{"pull-request", "-m", message, "-o"})
	execError := commands.CmdRunner.Call(cmd, args)
	return execError.Err
}
