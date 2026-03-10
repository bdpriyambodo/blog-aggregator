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

func (c *Config) SetUser() {
	c.CurrentUserName = "prymbd"

	err := write(*c)
	if err != nil {
		fmt.Println("Error - writing function", err)
	}
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
