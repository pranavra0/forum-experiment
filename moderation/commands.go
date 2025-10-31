package main

import (
	"fmt"
	"strings"
)

func HandleCommand(input string, client *APIClient) error {
	args := strings.Fields(input)
	if len(args) == 0 {
		return nil
	}

	switch args[0] {
	case "help":
		fmt.Println(`
Available commands:
  login                   - Log in as admin
  get threads             - List all threads
  delete thread <id>      - Delete a thread by ID
  exit                    - Exit the console
`)
	case "login":
		return handleLogin(client)
	case "get":
		if len(args) >= 2 && args[1] == "threads" {
			return client.GetThreads()
		}
	case "delete":
		if len(args) == 3 && args[1] == "thread" {
			return client.DeleteThread(args[2])
		}
	default:
		fmt.Println("Unknown command. Type 'help' for options.")
	}
	return nil
}
