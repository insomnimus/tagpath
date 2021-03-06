package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"

	"github.com/insomnimus/tagpath/engine"
)

var (
	flagA   = false
	exeName = "tagpath"
	file    string
	elem    string
)

func showHelp() {
	log.Printf(`%s, generate query selectors
usage:
	%s [options] filename|url '<html element>'
options are:
	-a, --all: print all matches instead of just the first
	-h, --help: show this message`,
		exeName, exeName)
	os.Exit(0)
}

func showAll(nodes []*html.Node) {
	for i, n := range nodes {
		fmt.Printf("##%d:\n", i)
		fmt.Printf("query selector:\n%s\n", engine.QuerySelector(n))
		fmt.Println("\npath:")
		for _, p := range engine.NodePath(n) {
			fmt.Println(p)
			fmt.Println("-")
		}
	}
}

func main() {
	exeName = filepath.Base(os.Args[0])
	exeName = strings.TrimSuffix(exeName, ".exe")
	log.SetFlags(0)
	log.SetPrefix("")
	if len(os.Args) == 1 {
		showHelp()
	}

	for _, a := range os.Args[1:] {
		if a == "-a" || a == "--all" {
			flagA = true
			continue
		}
		if a == "-h" || a == "--help" {
			showHelp()
		}

		if file == "" {
			file = a
		} else {
			elem = a
		}
	}
	if file == "" {
		log.Fatal("missing arguments: file and target")
	}
	if elem == "" {
		log.Fatal("missing argument: target")
	}
	var r io.Reader
	if strings.HasPrefix(file, "https://") ||
		strings.HasPrefix(file, "http://") {
		resp, err := http.Get(file)
		if err != nil {
			log.Fatal(err)
		}
		r = resp.Body
		defer resp.Body.Close()
	} else {
		f, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		r = f
	}

	q, err := engine.NewQuery(elem)
	if err != nil {
		log.Fatal(err)
	}
	nodes, err := q.FindIn(r)
	if err != nil {
		log.Fatal(err)
	}
	if len(nodes) == 0 {
		log.Fatal("No matches found")
	}
	if flagA && len(nodes) > 1 {
		showAll(nodes)
		return
	}

	fmt.Printf("single selector:\n%s\n", q.Selector())
	fmt.Printf("full selector:\n%s\n", engine.QuerySelector(nodes[0]))
	fmt.Println("path:")
	for _, n := range engine.NodePath(nodes[0]) {
		fmt.Printf("%s\n-\n", n)
	}
}
