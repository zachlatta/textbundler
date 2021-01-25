package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/zachlatta/textbundler"
	"github.com/zachlatta/textbundler/util"
)

var processAttachments, useGitDates bool
var toAppend string

func init() {
	RootCmd.PersistentFlags().BoolVarP(
		&processAttachments,
		"process-attachments",
		"p",
		false,
		"Replace links to local files with Bear-compatible tags to ease processing",
	)
	RootCmd.PersistentFlags().BoolVarP(
		&useGitDates,
		"git-dates",
		"g",
		false,
		"Instead of using OS creation / modification dates of Markdown file, use the dates from git commit history (must be in a git repo & have git CLI)",
	)
	RootCmd.PersistentFlags().StringVarP(
		&toAppend,
		"append",
		"a",
		"",
		"Text to append to end of Markdown file. Use %f to template the original filename.",
	)
}

// RootCmd handles the base case for textbundler: processing Markdown files.
var RootCmd = &cobra.Command{
	Use:   "textbundler [file] [file2] [file3]...",
	Short: "Convert markdown files into textbundles",
	Run: func(md *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "Please pass at least one argument.")
			os.Exit(1)
		}

		for _, mdPath := range args {
			if err := process(mdPath); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}
	},
}

func process(mdPath string) error {
	contents, err := ioutil.ReadFile(mdPath)
	if err != nil {
		return err
	}

	absMdPath, err := filepath.Abs(mdPath)
	if err != nil {
		return err
	}

	var creation, change time.Time

	if useGitDates {
		creation, err = util.GetGitBirthTime(absMdPath)
		if err != nil {
			return err
		}

		change, err = util.GetGitModTime(absMdPath)
		if err != nil {
			return err
		}
	} else {
		creation, err = util.GetBirthTime(absMdPath)
		if err != nil {
			return err
		}

		change, err = util.GetModTime(absMdPath)
		if err != nil {
			return err
		}
	}

	err = Textbundler.GenerateBundle(
		contents,
		absMdPath,
		creation,
		change,
		filepath.Dir(absMdPath)+"/",
		processAttachments,
		strings.Replace(toAppend, `\n`, "\n", -1),
	)
	if err != nil {
		return err
	}

	return nil
}

// Execute begins the CLI processing flow
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
