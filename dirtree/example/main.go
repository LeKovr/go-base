package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"

	"github.com/LeKovr/go-base/dirtree"
)

func main() {

	log = logrus.New()
	if len(os.Args) < 1 { // {3 {
		log.Printf("Use: %s root ext", os.Args[0])
		os.Exit(1)
	}
	root := "../test" //os.Args[1] // .
	ext := ".md"      //os.Args[2]  // .md

	t := dirtree.New(log, root, ext)

	b, _ := json.MarshalIndent(&t, "", "  ")
	fmt.Println(string(b))

}
