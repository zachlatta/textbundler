package util

import (
	"errors"
	"net/url"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/djherbis/times"
)

// IsValidURL attempts to parse the given text as a URL, returning true on
// success and false on failure.
func IsValidURL(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}

	return true
}

// GetBirthTime returns the given filepath's birth time, returning an error if
// the OS does not support file birth times.
func GetBirthTime(filepath string) (time.Time, error) {
	t, err := times.Stat(filepath)
	if err != nil {
		return time.Time{}, err
	}

	if !t.HasBirthTime() {
		return time.Time{}, errors.New("current OS does not support getting creation time of files")
	}

	return t.BirthTime(), nil
}

func GetModTime(filepath string) (time.Time, error) {
	t, err := times.Stat(filepath)
	if err != nil {
		return time.Time{}, err
	}

	if !t.HasChangeTime() {
		return time.Time{}, errors.New("current OS does not support getting change time of files")
	}

	return t.ChangeTime(), nil
}

// GetGitBirthTime returns the time the given filepath was first committed in a
// git repo. Assumes the given filepath is in a git repo and exists in history.
func GetGitBirthTime(path string) (time.Time, error) {
	_, err := exec.LookPath("git")
	if err != nil {
		return time.Time{}, errors.New("git CLI not found in PATH (you should install git)")
	}

	dir := filepath.Dir(path)
	filename := filepath.Base(path)

	cmd := exec.Command("git", "log", "--diff-filter=A", "--follow", "--format=%aD", "-1", "--", filename)
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		return time.Time{}, err
	}

	rawTime := strings.TrimSpace(string(output))

	creationTime, err := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", string(rawTime))
	if err != nil {
		return time.Time{}, err
	}

	return creationTime, nil
}

func GetGitModTime(path string) (time.Time, error) {
	_, err := exec.LookPath("git")
	if err != nil {
		return time.Time{}, errors.New("git CLI not found in PATH (you should install git)")
	}

	dir := filepath.Dir(path)
	filename := filepath.Base(path)

	cmd := exec.Command("git", "log", "--format=%aD", "-1", "--", filename)
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		return time.Time{}, err
	}

	rawTime := strings.TrimSpace(string(output))

	changeTime, err := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", string(rawTime))
	if err != nil {
		return time.Time{}, err
	}

	return changeTime, nil
}

// SetBirthTime sets the given filepath's birth time, returning an error if a
// strategy is not implemented for setting file birthtimes in the current OS.
func SetBirthTime(filepath string, time time.Time) error {
	_, err := exec.LookPath("SetFile")
	if err != nil {
		return errors.New("current OS / textbundler implementation does not support setting file birth times")
	}

	cmd := exec.Command("SetFile", "-d", time.Format("01/02/2006 03:04:05 PM"), filepath)
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// SetModTime sets the given filepath's modification time, returning an error if
// a strategy is not implemented for setting file modification times in the
// current OS.
func SetModTime(filepath string, time time.Time) error {
	_, err := exec.LookPath("SetFile")
	if err != nil {
		return errors.New("current OS / textbundler implementation does not support setting file modification times")
	}

	cmd := exec.Command("SetFile", "-m", time.Format("01/02/2006 03:04:05 PM"), filepath)
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
