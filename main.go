package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, err := reader.ReadString('\n')
		if text = strings.Trim(text, "\n"); text != "" {
			err = try(3, time.Second, func() error {
				name, ok, err := query(text)
				if err != nil {
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
				fmt.Fprintf(os.Stderr, "%s %s", text, err)
			}
		}
		if err != nil {
			if err == io.EOF {
				return
			}
			fmt.Fprintf(os.Stderr, "%s %s", strings.Trim(text, "\n"), err)
		}
	}
}
