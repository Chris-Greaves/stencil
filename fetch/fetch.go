package fetch

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

// IsPath is used to determin if the input string is a valid path on the local system
func IsPath(input string) bool {
	info, err := os.Stat(input)

	if err != nil || !info.IsDir() {
		return false
	}

	return true
}

// IsGitInstalled checks if Git is installed on the PC
func IsGitInstalled() bool {
	_, err := exec.LookPath("git")
	if err != nil {
		return false
	}
	return true
}

// IsGitURL is used to determin if the input string is a valid Git Url
func IsGitURL(input string) bool {
	cmd := "git"
	args := []string{"ls-remote", input}
	if err := exec.Command(cmd, args...).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return false
	}

	return true
}

// PullTemplate clones the template from its git repo
func PullTemplate(repo string) (string, error) {
	dir, err := ioutil.TempDir("", "template-")
	if err != nil {
		return "", err
	}

	cmd := exec.Command("git", "clone", repo, dir)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return "", err
	}
	return dir, nil
}
