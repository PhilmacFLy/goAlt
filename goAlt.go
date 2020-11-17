package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"text/template"
)

const (
	startLine = 3
	mainCol   = 0
	altStart  = 2
)

type alt struct {
	Name  string
	Comma bool
}

type character struct {
	Main string
	Alts []string
}

type charExport struct {
	Main   string
	Alts   []alt
	IsMain bool
	Comma  bool
}

func (c *character) getExport() []charExport {
	var res []charExport
	//Erstmal den Main Char
	cm := charExport{Main: c.Main, Comma: true}
	for _, a := range c.Alts {
		al := alt{Name: a, Comma: true}
		cm.Alts = append(cm.Alts, al)
	}
	cm.Alts[len(cm.Alts)-1].Comma = false
	cm.IsMain = true
	res = append(res, cm)

	//Dann der Rest
	for _, c1 := range c.Alts {
		cres := charExport{Main: c1}
		for _, c2 := range c.Alts {
			if c2 != c1 {
				al := alt{Name: c2, Comma: true}
				cres.Alts = append(cres.Alts, al)
			}
		}
		al := alt{Name: c.Main, Comma: false}
		cres.Alts = append(cres.Alts, al)
		cres.IsMain = false
		cres.Comma = true
		res = append(res, cres)
	}
	return res
}

func main() {
	var filename string
	if len(os.Args) > 1 {
		filename = os.Args[1]
	} else {
		filename = "characters.csv"
	}
	generatePlugindata(filename)
}

func generatePlugindata(filename string) {
	Chars := loadAllChars(filename)

	t, err := template.ParseFiles("template.txt")
	if err != nil {
		log.Fatal("Error parsing template:", err)
	}
	outfile, err := os.Create("AltTracker.plugindata")
	if err != nil {
		log.Println("Error parsing plugindata:", err)
	}
	fmt.Println(Chars)
	w := bufio.NewWriter(outfile)
	t.Execute(w, &Chars)
	w.Flush()
}

func loadAllChars(filename string) []charExport {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal("Error opening File:", err)
	}
	var res []charExport
	line := 0
	r := csv.NewReader(bufio.NewReader(f))
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		line++
		if (line <= startLine) || (record[mainCol] == "") {
			continue
		}
		var c character
		c.Main = record[mainCol]
		col := altStart
		for {
			if record[col] != "" {
				c.Alts = append(c.Alts, record[col])
				col++
			} else {
				break
			}
		}
		fmt.Println(c)
		res = append(res, c.getExport()...)
	}
	return res
}
