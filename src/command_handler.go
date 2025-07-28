package main

import "fmt"

type CommandHandlerInterface interface {
	HandleCommand(args []string, storage TodoStorageInterface) error
	AddTodo(args []string, storage TodoStorageInterface) error
	ShowNextUndoneTodo(storage TodoStorageInterface) error
}

type CommandHandler struct{}

func (ch *CommandHandler) HandleCommand(args []string, storage TodoStorageInterface) error {
	if len(args) == 0 {
		return nil
	}

	switch args[0] {
	case "a":
		ch.AddTodo(args, storage)

	case "d":
		ch.setNextUndoneTodoDone(storage)

	default:
		fmt.Printf("Unknown command: %s\n", args[0])
	}

	return nil
}

func (ch *CommandHandler) AddTodo(args []string, storage TodoStorageInterface) error {
	text := ""
	if len(args) > 1 {
		text = args[1]
		for i := 2; i < len(args); i++ {
			text += " " + args[i]
		}
	}

	todo := &Todo{
		Title: text,
	}
	if err := storage.Save(todo); err != nil {
		return fmt.Errorf("error saving todo: %w", err)
	}
	fmt.Printf("Added todo: %s\n", text)

	return nil
}

func (ch *CommandHandler) ShowNextUndoneTodo(storage TodoStorageInterface) error {
	todo, _ := storage.GetOldestUndone()

	if todo == nil {
		fmt.Println("\nNo undone todos found.")
		return nil
	}

	fmt.Printf("\nNext undone todo: \"%s\"\n", todo.Title)
	return nil
}

func (ch *CommandHandler) setNextUndoneTodoDone(storage TodoStorageInterface) error {
	todo, err := storage.GetOldestUndone()
	if err != nil {
		return fmt.Errorf("error getting oldest undone todo: %w", err)
	}
	if todo == nil {
		fmt.Println("No undone todos found.")
		return nil
	}

	todo.Done = true
	if err := storage.Update(todo); err != nil {
		return fmt.Errorf("error updating todo: %w", err)
	}
	fmt.Printf("Marked todo as done: %s\n", todo.Title)
	return nil
}
