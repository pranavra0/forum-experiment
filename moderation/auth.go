package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func handleLogin(client *APIClient) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Username: ")
	username, _ := reader.ReadString('\n')
	fmt.Print("Password: ")
	password, _ := reader.ReadString('\n')

	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)

	if err := client.Login(username, password); err != nil {
		return fmt.Errorf("login failed: %w", err)
	}
	fmt.Println("Logged in successfully!")
	return nil
}
