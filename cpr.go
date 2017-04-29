package cpr

import (
	"errors"
	"os"

	git "gopkg.in/src-d/go-git.v4"
)

var (
	ErrNoGitParent = errors.New("Could not find a parent git directory.")
)

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
