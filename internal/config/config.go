package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

type State struct {
	ConfigPointer *Config
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

	(s.ConfigPointer).SetUser(cmd.Args[0])

	fmt.Println("The user has been set!")

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
