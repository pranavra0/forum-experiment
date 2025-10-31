package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("🧭 Forum Moderator Console")
	fmt.Println("Type 'help' for available commands.\n")

	reader := bufio.NewReader(os.Stdin)
	client := NewAPIClient("http://localhost:8080")

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		cmd := strings.TrimSpace(input)

		if cmd == "exit" {
			fmt.Println("Goodbye!")
			break
		}

		if err := HandleCommand(cmd, client); err != nil {
			fmt.Println("Error:", err)
		}
	}
}
