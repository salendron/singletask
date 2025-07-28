package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var commandHandler CommandHandlerInterface = &CommandHandler{}

func main() {
	todoStorage, err := NewSQLiteTodoStorage("singletask.db")
	if err != nil {
		fmt.Printf("Error initializing storage: %v\n", err)
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Welcome to Singletask!")
	fmt.Println("To add a new todo, type 'a' followed by your todo title.")
	fmt.Println("To mark the oldest undone todo as done, type 'd'.")
	fmt.Println("To list all undone todos, type 'l'.")
	fmt.Println("To terminate the app, type 'q'.")

	for {
		commandHandler.ShowNextUndoneTodo(todoStorage)

		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		if input == "q" {
			break
		}

		args := strings.Fields(input)
		if len(args) == 0 {
			continue
		}

		if err := commandHandler.HandleCommand(args, todoStorage); err != nil {
			fmt.Printf("Error handling command: %v\n", err)
		}
	}
}
