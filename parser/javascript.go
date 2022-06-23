package parser

import (
	"context"
	"fmt"

	"github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/javascript"
)

type Javascript struct {
	parser *sitter.Parser
}

func NewJavascript() *Javascript {
	p := sitter.NewParser()
	p.SetLanguage(javascript.GetLanguage())

	return &Javascript{
		parser: p,
	}
}

func (j Javascript) Parse(ctx context.Context, contents []byte) *sitter.Tree {
	t, err := j.parser.ParseCtx(ctx, nil, contents)
	if err != nil {
		panic(err)
	}

	return t
}

// FindFsPromisesVarName returns name of the variable for require("fs/promises") module
func (j Javascript) FindFsPromisesVarName(tree *sitter.Tree, contents []byte) string {
	results := searchRecursive(tree.RootNode(),
		make([]*sitter.Node, 0),
		contents,
		"call_expression",
		`require("fs/promises")`)

	if len(results) != 1 {
		panic(`want exactly one require("fs/promises"), got ` + fmt.Sprint(len(results)))
	}

	// assumption fs = require("fs/promises"), can be smarter if need to
	return results[0].NamedChild(0).Child(0).Content(contents)
}

// FindReadFile returns node that contains lexical_declaration for readFile func
func (j Javascript) FindReadFile(tree *sitter.Tree, contents []byte, varName string) []*sitter.Node {
	return searchRecursive(
		tree.RootNode(),
		make([]*sitter.Node, 0),
		contents,
		"member_expression",
		varName+".readFile")
}

func searchRecursive(n *sitter.Node, results []*sitter.Node,
	contents []byte, typeMatch, contentMatch string) []*sitter.Node {
	if n == nil {
		return results
	}

	for i := 0; i < int(n.NamedChildCount()); i++ {
		child := n.NamedChild(i)

		if child.Type() == typeMatch &&
			child.Content(contents) == contentMatch {
			fullNode := getDeclarationParent(child)

			results = append(results, fullNode)
		}

		results = searchRecursive(child, results, contents, typeMatch, contentMatch)
	}

	return results
}

func getDeclarationParent(n *sitter.Node) *sitter.Node {
	if n == nil {
		return nil
	}

	if n.Type() == "lexical_declaration" {
		return n
	}

	return getDeclarationParent(n.Parent())
}
