package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/bdpriyambodo/blog-aggregator/internal/config"
	"github.com/bdpriyambodo/blog-aggregator/internal/database"
	_ "github.com/lib/pq"
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

	dbUrl := s.ConfigPointer.DbURL
	db, _ := sql.Open("postgres", dbUrl)

	dbQueries := database.New(db)

	s.DataBase = dbQueries

	// register commands
	var c config.Commands
	c.Handlers = make(map[string]func(*config.State, config.Command) error)

	c.Register("login", config.HandlerLogin)
	c.Register("register", config.HandlerRegister)
	c.Register("reset", config.HandlerReset)
	c.Register("users", config.HandlerGetUsers)
	c.Register("agg", config.HandlerAgg)

	// ACTUAL RUN
	userArgs := os.Args
	for i, arg := range userArgs {
		fmt.Println(i, arg)
	}

	arg1 := os.Args[1]
	if arg1 != "reset" && arg1 != "users" && arg1 != "agg" {
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

		fmt.Printf("Command name: %s\n", cmd.Name)
		fmt.Printf("Command argument: %s\n", cmd.Args)
		c.Run(&s, cmd)
	} else {
		c.Run(&s, config.Command{Name: arg1})
	}

}
