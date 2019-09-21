package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	override = flag.Bool("w", false, "write result to (source) file instead of stdout")
)

func visitFile(filename string, f os.FileInfo, err error) error {
	if err == nil && strings.HasSuffix(f.Name(), ".json") {
		err = processFile(filename)
	}
	if err != nil && !os.IsNotExist(err) {
		log.Println(err)
	}
	return nil
}

func walkDir(filename string) {
	_ = filepath.Walk(filename, visitFile)
}

func processFile(filename string) (err error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	m := []interface{}{}
	err = json.Unmarshal(data, &m)
	if err != nil {
		return
	}

	dataFmt, _ := json.MarshalIndent(m, "", "  ")
	if *override {
		err = ioutil.WriteFile(filename, dataFmt, 0755)
	} else {
		fmt.Println(string(dataFmt))
	}
	return
}

func main() {
	flag.Parse()

	if len(os.Args) == 1 {
		flag.Usage()
	}

	for i := 0; i < flag.NArg(); i++ {
		path := flag.Arg(i)

		switch dir, err := os.Stat(path); {
		case err != nil:
			log.Println(err)
		case dir.IsDir():
			walkDir(path)
		default:
			if err := processFile(path); err != nil {
				log.Println(err)
			}

		}
	}

}
