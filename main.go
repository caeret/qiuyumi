package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	verbose bool
	n       int
	queue   = make(chan string, 1024)
)

func init() {
	flag.BoolVar(&verbose, "v", false, "verbose output.")
	flag.IntVar(&n, "n", 10, "count of thread to execute the queries.")
	flag.Parse()
	if n < 1 {
		n = 1
	}
}

func main() {
	if len(flag.Args()) > 0 {
		go readFromCommand()
	} else {
		go readFromStdin()
	}

	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for name := range queue {
				name = strings.TrimSpace(name)
				if name != "" {
					tryQuery(name)
				}
			}
		}()
	}

	wg.Wait()
}

func readFromCommand() {
	defer close(queue)
	for _, name := range flag.Args() {
		queue <- name
	}
}

func readFromStdin() {
	defer close(queue)
	reader := bufio.NewReader(os.Stdin)
	for {
		text, err := reader.ReadString('\n')
		if text != "" {
			queue <- text
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
