package main

import (
	"fmt"
	"os"

	"flag"

	"path/filepath"

	"github.com/alistanis/cpr"
)

var options *cpr.Options

func main() {
	checkAndExit(run())
}

func run() error {

	r, err := cpr.Open(".")
	if err != nil {
		return err
	}

	url, err := cpr.GithubURL(r)
	if err != nil {
		return err
	}
	fmt.Println(options)
	if options.UserName == "" || options.Password == "" {
		home, err := cpr.HomeDir()
		if err != nil {
			return err
		}

		configFile := filepath.Join(home, cpr.ConfigFileName)

		c, err := cpr.LoadConfig(configFile)
		if err != nil {
			return err
		}
		options.UserName = c.User
		options.Password = string(c.Password)
	}

	pr, _, err := options.PullRequest(url)
	if err != nil {
		return err
	}
	fmt.Println(*pr.Title, " was created at ", *pr.CreatedAt, " by ", *pr.User.Name)
	if len(pr.Assignees) > 0 {
		fmt.Println("Assignees:")
		for _, a := range pr.Assignees {
			fmt.Println("\t", *a.Name)
		}
	}
	return nil
}

func init() {
	opts, err := cpr.ParseOptions(flag.CommandLine, os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
	options = opts

	if options.GenerateConfig {
		checkAndExit(cpr.GenerateConfig())
	}

	err = options.Validate()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}

}

func checkAndExit(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		options.FlagSet.Usage()
		os.Exit(-1)
	}
	os.Exit(0)
}
