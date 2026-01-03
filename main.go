package main

import (
	"fmt"
	"os"

	"github.com/wfcornelissen/blogag/internal/config"
	"github.com/wfcornelissen/blogag/internal/handling"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error fetching original cfg")
		fmt.Println(err)

		return
	}
	newState := config.State{State: &cfg}
	cmds := handling.Commands{
		Commands: make(map[string]func(*config.State, handling.Command) error),
	}
	cmds.Register("login", handling.HandlerLogin)
	var newCommand handling.Command
	input := os.Args
	switch len(input) {
	case 0:
		fmt.Println("Too few arguements")
		os.Exit(1)
		return
	case 1:
		fmt.Println("Too few arguements")
		os.Exit(1)
		return
	case 2:
		newCommand = handling.Command{
			Name: input[1],
		}
	default:
		newCommand = handling.Command{
			Name: input[1],
			Args: input[2:],
		}
	}
	err = cmds.Run(&newState, newCommand)
	if err != nil {
		fmt.Println("No username supplied")
		fmt.Println(err)
		os.Exit(1)
	}

}
