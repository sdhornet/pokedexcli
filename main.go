package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("pokedex > ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		command := cleanInput(line)
		if len(command) == 0 {
			continue
		}

		fmt.Printf("Your command was: %s\n", command[0])
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
		os.Exit(1)
	}
}
