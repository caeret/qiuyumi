package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

var (
	verbose bool
)

func init() {
	flag.BoolVar(&verbose, "v", false, "verbose output.")
	flag.Parse()
}

func main() {
	if len(flag.Args()) > 0 {
		for _, name := range flag.Args() {
			tryQuery(name)
		}
		return
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		text, err := reader.ReadString('\n')
		if text = strings.Trim(text, "\n"); text != "" {
			tryQuery(text)
		}
		if err != nil {
			if err == io.EOF {
				return
			}
			log("%s %s", strings.Trim(text, "\n"), err)
		}
	}
}

func tryQuery(name string) {
	err := try(3, time.Second, func() error {
		debug("try to query %s", name)
		name, ok, err := query(name)
		if err != nil {
			debug("fail to query: %s %s", name, err.Error())
			return err
		}
		if ok {
			fmt.Printf("%-16s available\n", name)
		} else {
			fmt.Printf("%-16s invalid\n", name)
		}
		return nil
	})
	if err != nil {
		log("%s %s", name, err)
	}
}

func debug(format string, args ...interface{}) {
	if verbose {
		if !strings.HasSuffix(format, "\n") {
			format += "\n"
		}
		fmt.Fprintf(os.Stderr, format, args...)
	}
}

func log(format string, args ...interface{}) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	fmt.Fprintf(os.Stderr, format, args...)
}
