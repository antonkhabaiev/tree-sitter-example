package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"

	"tree-sitter-example/parser"
)

func main() {
	println("start")
	defer println("finished")

	filePath := flag.String(
		"file-path",
		"parser/javascript_unit_test_mod.js",
		"file path to source code that needs to be parsed")

	if filePath == nil {
		panic("file-path is required")
	}

	contents, err := ioutil.ReadFile(*filePath)
	if err != nil {
		panic(fmt.Errorf("open file: %w", err))
	}

	p := parser.NewJavascript()
	ctx := context.Background()

	tree := p.Parse(ctx, contents)

	promisesVarName := p.FindFsPromisesVarName(tree, contents)

	println("fs/promises variable named: " + promisesVarName)

	readFileMentions := p.FindReadFile(tree, contents, promisesVarName)

	if len(readFileMentions) == 0 {
		println("no readFile mentions found")

		return
	}

	println("readFile mentions:")

	for _, m := range readFileMentions {
		fmt.Println(fmt.Sprintf("<line %d, %s>", m.StartPoint().Row+1, m.Content(contents)))
	}
}
