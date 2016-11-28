package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/hyperhq/hyperd/client"
	"github.com/hyperhq/hyperd/client/api"
)

func nodeList(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var nodeList []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		nodeList = append(nodeList, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return nodeList
}

func main() {
	var (
		proto = "http"
		nodeListPath = "/tmp/hyper/nodelist"
	)
	nodes := nodeList(nodeListPath)
	cli := client.NewHyperClient(proto, nodes[0], nil)

	// set the flag to output
	flHelp := flag.Bool("help", false, "Help Message")
	flVersion := flag.Bool("version", false, "Version Message")
	flag.Usage = func() { cli.Cmd("help") }
	flag.Parse()
	if flag.NArg() == 0 {
		cli.Cmd("help")
		return
	}
	if *flHelp == true {
		cli.Cmd("help")
	}
	if *flVersion == true {
		cli.Cmd("version")
	}

	if err := cli.Cmd(flag.Args()...); err != nil {
		if sterr, ok := err.(api.StatusError); ok {
			if sterr.Status != "" {
				fmt.Printf("%s ERROR: %s\n", os.Args[0], err.Error())
				os.Exit(-1)
			}
			os.Exit(sterr.StatusCode)
		}

		fmt.Printf("%s ERROR: %s\n", os.Args[0], err.Error())
		os.Exit(-1)
	}
}
