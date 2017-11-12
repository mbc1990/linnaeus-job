package main

import "fmt"
import "encoding/json"
import "os"

// Configuration struct that conf json file is read into
type Configuration struct {
	LinnaeusPort    int
	ImageDirAbsPath string
	PGHost          string
	PGPort          int
	PGUser          string
	PGPassword      string
	PGDbname        string
}

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("Usage: ./main <absolute path to configuration file>")
		return
	}
	file, _ := os.Open(args[0])
	decoder := json.NewDecoder(file)
	var conf = Configuration{}
	err := decoder.Decode(&conf)
	if err != nil {
		fmt.Println("error:", err)
	}
	job := NewLinnaeusJob(&conf)
	job.ClassifyImages()
}
