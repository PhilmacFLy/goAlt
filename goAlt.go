package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Character struct {
	Main string
	Alts []string
}

func (c *Character) Save() error {
	b, err := json.MarshalIndent(&c, "", "    ")
	if err != nil {
		return err
	}
	ioutil.WriteFile("chars/"+c.Main+".json", b, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (c *Character) Load(f string) error {
	body, err := ioutil.ReadFile("chars/" + f)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &c)
	if err != nil {
		return err
	}
	return nil
}

func (c *Character) GetChars() string {
	result := "[\"" + c.Main + "\"]\n"
	result = result + "{\n"
	result = result + "\t[\"Alt\"] = \n"
	result = result + "\t{\n"
	for _, a := range c.Alts {
		result = result + "\t\t[\"" + a + "\"],\n"
	}
	result = result + "\t},\n"

	return result
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	text := "blub"
	var c Character
	err := c.Load("Gullofbog.json")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(c)
	fmt.Println(c.GetChars())

	for text != "q\n" {
		fmt.Println("===== Alttracker Generation ======")
		fmt.Println("Please choose a Command:")
		fmt.Println("1) Enter a New Character")
		fmt.Println("2) Generate Plugindata")
		fmt.Println("q) Exit")
		text, _ = reader.ReadString('\n')
		switch text {
		case "1\n":
			enterCharacter(reader)
		case "2\n":
			generatePlugindata()
		}
	}
}

func enterCharacter(reader *bufio.Reader) {
	var c Character
	fmt.Println("Please enter Main Character:")
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	c.Main = text
	for text != "" {
		fmt.Println("Please enter Alt Character:")
		text, _ = reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		if text != "" {
			c.Alts = append(c.Alts, text)
		}
	}
	err := c.Save()
	if err != nil {
		log.Println("Error:", err)
	}

}

func generatePlugindata() {

}
