package main

import (
	"fmt"
	"os"

	"github.com/bdpriyambodo/blog-aggregator/internal/config"
)

func main() {
	// c := config.Read()
	// c.SetUser("prymbd")

	// //check
	// cNew := config.Read()
	// fmt.Println(cNew.DbURL)
	// fmt.Println(cNew.CurrentUserName)

	var s config.State
	s.ConfigPointer = config.Read()

	var c config.Commands
	c.Handlers = make(map[string]func(*config.State, config.Command) error)

	c.Register("login", config.HandlerLogin)

	userArgs := os.Args
	// for i, arg := range userArgs {
	// 	fmt.Println(i, arg)
	// }

	if len(userArgs) < 2 {
		fmt.Println("Not enough argument")
		os.Exit(1)
	}

	if len(userArgs) < 3 {
		fmt.Println("Username required")
		os.Exit(1)
	}

	cmd := config.Command{
		Name: userArgs[1],
		Args: userArgs[2:],
	}

	c.Run(&s, cmd)
}
