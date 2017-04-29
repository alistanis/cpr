package cpr

import (
	"fmt"
	"os"

	"github.com/alistanis/cpr"
)

func main() {
	r, err := cpr.Open(".")
	exitError(err)
	c, err := r.Config()
	exitError(err)
	b, err := c.Marshal()
	exitError(err)
	fmt.Println(string(b))
}

func exitError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}
