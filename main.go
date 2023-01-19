package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

//go:embed ts-template.txt
//go:embed tsx-template.txt
//go:embed hpp-template.txt
//go:embed cpp-template.txt
//go:embed go-template.txt
//go:embed sh-template.txt
var content embed.FS

func template_name(t string) string {
	switch t {
	case "ts":
		return "ts-template.txt"
	case "tsx":
		return "tsx-template.txt"
	case "go":
		return "go-template.txt"
	case "cpp":
		return "cpp-template.txt"
	case "hpp":
		return "hpp-template.txt"
	case "sh":
		return "sh-template.txt"
	}
	panic("Unsupported type: " + t)
	return ""
}

// Info collects the info that will be used for the template.
type Info struct {
	Author  string
	Company string `json: "company"`
	File    string
	Name    string
	Title   string
	Project string `json: "project"`
	Today   string
	Year    string
}

const YYYYMMDD = "2006/01/02"
const YYYY = "2006"

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: new <TYPE> <NAME>")
		fmt.Println("Available Types:\n\tts\n\ttsx\n\thpp\n\tcpp\n\tgo\n")
		return
	}
	i := info("project.json")
	i.Today = time.Now().Format(YYYYMMDD)
	i.Year = time.Now().Format(YYYY)
	i.File = os.Args[1]
	i.Name = os.Args[2]
	i.Title = strings.Title(os.Args[2])

	/*
		// using the function
		mydir, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(mydir)
	*/

	out, err := exec.Command("git", "config", "--get", "user.name").Output()
	if err != nil {
		panic(err)
	}
	i.Author = strings.Replace(string(out), "\n", "", 1)

	t := template.Must(template.ParseFS(content, template_name(i.File)))
	err = t.Execute(os.Stdout, i)
	if err != nil {
		panic(err)
	}
}

func info(filename string) Info {
	o := Info{}
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return o
	}
	_ = json.Unmarshal([]byte(file), &o)
	return o
}
