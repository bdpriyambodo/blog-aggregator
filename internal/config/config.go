package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"database/sql"

	"github.com/bdpriyambodo/blog-aggregator/internal/database"
	"github.com/bdpriyambodo/blog-aggregator/internal/xmlfetcher"

	"github.com/google/uuid"
)

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

type State struct {
	ConfigPointer *Config
	DataBase      *database.Queries
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Handlers map[string]func(*State, Command) error
}

func (c Commands) Run(s *State, cmd Command) error {
	handler, exists := c.Handlers[cmd.Name]
	if !exists {
		return fmt.Errorf("Unknown command: %s", cmd.Name)
	}
	return handler(s, cmd)
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.Handlers[name] = f
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("Empty handlers/argument")
	}

	inputName := cmd.Args[0]

	_, err := s.DataBase.GetUser(context.Background(), inputName)
	if err != nil {
		fmt.Println("User not exist")
		os.Exit(1)
	}

	(s.ConfigPointer).SetUser(cmd.Args[0])

	fmt.Println("The user has been set!")

	return nil

}

func HandlerRegister(s *State, cmd Command) error {
	// fmt.Println("1")

	if len(cmd.Args) == 0 {
		return fmt.Errorf("Empty handlers/argument")
	}

	// fmt.Println("2")

	inputName := cmd.Args[0]

	var userParams database.CreateUserParams
	userParams.ID = uuid.New()
	userParams.CreatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	userParams.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	userParams.Name = inputName

	_, err := s.DataBase.GetUser(context.Background(), inputName)
	if err == nil {
		fmt.Println("User already exist")
		os.Exit(1)
	}

	_, err = s.DataBase.CreateUser(context.Background(), userParams)
	if err != nil {
		fmt.Printf("Error in creating user: %s\n", err)
		return err
	}

	s.ConfigPointer.SetUser(inputName)

	fmt.Println("User has been created")
	user, err := s.DataBase.GetUser(context.Background(), inputName)
	fmt.Printf("UUID: %s\n", user.ID.String())
	fmt.Printf("Created at: %v\n", user.CreatedAt.Time)
	fmt.Printf("Updated at: %v\n", user.UpdatedAt.Time)
	fmt.Printf("Name: %s\n", user.Name)

	return nil
}

func HandlerReset(s *State, cmd Command) error {
	err := s.DataBase.DeleteAllUsers(context.Background())
	if err != nil {
		fmt.Printf("Error in deleting all users from table: %s\n", err)
	}

	return nil
}

func HandlerGetUsers(s *State, cmd Command) error {
	users, err := s.DataBase.GetUsers(context.Background())
	if err != nil {
		fmt.Printf("Error in getting all users from table: %s\n", err)
	}

	// fmt.Printf("current username: %s\n", s.ConfigPointer.CurrentUserName)

	for _, user := range users {
		if user != s.ConfigPointer.CurrentUserName {
			fmt.Printf(" * %s\n", user)
		} else {
			fmt.Printf(" * %s (current)\n", user)
		}
	}

	return nil
}

func HandlerAgg(s *State, cmd Command) error {
	feedURL := "https://www.wagslane.dev/index.xml"
	ctx := context.Background()
	rss_ptr, err := xmlfetcher.FetchFeed(ctx, feedURL)
	if err != nil {
		fmt.Println("Error in fetching feed: ", err)
		return err
	}

	fmt.Printf("%+v\n", rss_ptr)

	return nil
}

func HandlerAddFeed(s *State, cmd Command) error {

	ctx := context.Background()

	currentUser := s.ConfigPointer.CurrentUserName
	currentUserData, err := s.DataBase.GetUser(context.Background(), currentUser)
	if err != nil {
		fmt.Println("Current user not exist")
		os.Exit(1)
	}

	var feedArg database.CreateFeedParams

	feedArg.ID = uuid.New()
	feedArg.CreatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	feedArg.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	feedArg.Name = sql.NullString{String: cmd.Args[0], Valid: true}
	feedArg.Url = sql.NullString{String: cmd.Args[1], Valid: true}
	feedArg.UserID = uuid.NullUUID{UUID: currentUserData.ID, Valid: true}

	resultFeed, err := s.DataBase.CreateFeed(ctx, feedArg)
	if err != nil {
		fmt.Println("Error in adding feed: ", err)
		return err
	}

	fmt.Printf("ID: %v\n", resultFeed.ID)
	fmt.Printf("CreatedAt: %v (Valid: %v)\n", resultFeed.CreatedAt.Time, resultFeed.CreatedAt.Valid)
	fmt.Printf("UpdatedAt: %v (Valid: %v)\n", resultFeed.UpdatedAt.Time, resultFeed.UpdatedAt.Valid)
	fmt.Printf("Name: %s (Valid: %v)\n", resultFeed.Name.String, resultFeed.Name.Valid)
	fmt.Printf("Url: %s (Valid: %v)\n", resultFeed.Url.String, resultFeed.Url.Valid)
	fmt.Printf("UserID: %v (Valid: %v)\n", resultFeed.UserID.UUID, resultFeed.UserID.Valid)

	return nil
}

const configFileName = ".gatorconfig.json"

func Read() *Config {
	// homePath, err := os.UserHomeDir()
	// if err != nil {
	// 	fmt.Println("Error - home directory:", err)
	// }

	// fmt.Println(homePath)

	// filePath := homePath + "/.gatorconfig.json"
	// fmt.Println(filePath)

	filePath, _ := getConfigFilePath()

	bytes, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error - file read", err)
	}

	var result Config
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		fmt.Println("Error - unmarshal:", err)
	}

	// fmt.Println(result.DbURL)
	// fmt.Println(result.CurrentUserName)

	return &result

}

func (c *Config) SetUser(username string) error {
	// c.CurrentUserName = "prymbd"
	c.CurrentUserName = username

	err := write(*c)
	if err != nil {
		fmt.Println("Error - writing function", err)
		return err
	}

	return nil
}

func write(cfg Config) error {

	filePath, _ := getConfigFilePath()

	jsonData, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		fmt.Println("Error - marshal:", err)
		return err
	}

	err = os.WriteFile(filePath, jsonData, 0644)
	// 0644 translates to -rw-r--r--, which is the standard, secure default for creating files that the owner can modify but others can only read
	if err != nil {
		fmt.Println("Error - file writing:", err)
		return err
	}
	return nil
}

func getConfigFilePath() (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error - home directory:", err)
		return "", err
	}

	// fmt.Println("Home path:", homePath)

	filePath := homePath + "/" + configFileName
	// fmt.Println("File path:", filePath)

	return filePath, nil
}
