package main

import (
	"fmt"
	"os"

	"flag"

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
	info, err := cpr.GetRepoInfo(url)
	if err != nil {
		return err
	}
	fmt.Println(info)
	err = options.Validate()
	if err != nil {
		return err
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
}

func checkAndExit(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		options.FlagSet.Usage()
		os.Exit(-1)
	}
	os.Exit(0)
}
