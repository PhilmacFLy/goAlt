package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Alt struct {
	Name  string
	Comma bool
}

type Character struct {
	Main string
	Alts []string
}

type CharExport struct {
	Main   string
	Alts   []Alt
	IsMain bool
	Comma  bool
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

func (c *Character) GetExport() []CharExport {
	var res []CharExport
	//Erstmal den Main Char
	cm := CharExport{Main: c.Main, Comma: true}
	for _, a := range c.Alts {
		al := Alt{Name: a, Comma: true}
		cm.Alts = append(cm.Alts, al)
	}
	cm.Alts[len(cm.Alts)-1].Comma = false
	cm.IsMain = true
	res = append(res, cm)

	//Dann der Rest
	for _, c1 := range c.Alts {
		cres := CharExport{Main: c1}
		for _, c2 := range c.Alts {
			if c2 != c1 {
				al := Alt{Name: c2, Comma: true}
				cres.Alts = append(cres.Alts, al)
			}
		}
		al := Alt{Name: c.Main, Comma: false}
		cres.Alts = append(cres.Alts, al)
		cres.IsMain = false
		cres.Comma = true
		res = append(res, cres)
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	text := "blub"
	var c Character
	err := c.Load("Gullofbog.json")
	if err != nil {
		log.Fatal(err)
	}

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
	Chars := loadAllChars()

	t, err := template.ParseFiles("template.txt")
	if err != nil {
		log.Fatal(err)
	}
	outfile, err := os.Create("AltTracker.plugindata")
	if err != nil {
		log.Println(err)
	}
	fmt.Println(Chars)
	w := bufio.NewWriter(outfile)
	t.Execute(w, &Chars)
	w.Flush()
}

func loadAllChars() []CharExport {
	var res []CharExport
	files, _ := ioutil.ReadDir("chars")
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".json" {
			c := Character{}
			c.Load(f.Name())
			add := c.GetExport()
			res = append(res, add...)
		}
	}
	res[len(res)-1].Comma = false
	return res
}
