package main

import (
	"fmt"

	"github.com/bdpriyambodo/blog-aggregator/internal/config"
)

func main() {
	c := config.Read()
	c.SetUser()

	//check
	cNew := config.Read()
	fmt.Println(cNew.DbURL)
	fmt.Println(cNew.CurrentUserName)

}
