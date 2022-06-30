package parser

import (
	"context"

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

type VarNameMatch struct {
	Type    string
	Content string
}

func NewVarNameMatch(t, c string) VarNameMatch {
	return VarNameMatch{
		Type:    t,
		Content: c,
	}
}

// FindFsPromisesVarName returns name of the variable for require("fs/promises") module
func (j Javascript) FindFsPromisesVarName(tree *sitter.Tree, contents []byte) string {
	results := searchRecursive(tree.RootNode(),
		make([]*sitter.Node, 0),
		contents,
		[]VarNameMatch{
			NewVarNameMatch("call_expression", `require("fs/promises")`),
			NewVarNameMatch("member_expression", `require('fs').promises`),
			NewVarNameMatch("call_expression", `require('fs')`),
		},
	)

	for _, res := range results {
		if res.NamedChild(0).NamedChild(1).Content(contents) == `require('fs')` {
			if res.NamedChild(0).Child(0).Child(1).NamedChild(0).Content(contents) == "promises" {
				return res.NamedChild(0).Child(0).Child(1).NamedChild(1).Content(contents)
			}
		} else {
			return results[0].NamedChild(0).Child(0).Content(contents)
		}
	}

	panic("couldn't find any occurrences of fs.promises")
}

// FindReadFile returns node that contains lexical_declaration for readFile func
func (j Javascript) FindReadFile(tree *sitter.Tree, contents []byte, varName string) []*sitter.Node {
	return searchRecursive(
		tree.RootNode(),
		make([]*sitter.Node, 0),
		contents,
		[]VarNameMatch{
			NewVarNameMatch("member_expression", varName+".readFile"),
		},
	)
}

func searchRecursive(n *sitter.Node, results []*sitter.Node,
	contents []byte, matches []VarNameMatch) []*sitter.Node {
	if n == nil {
		return results
	}

	for i := 0; i < int(n.NamedChildCount()); i++ {
		child := n.NamedChild(i)

		for _, match := range matches {
			if child.Type() == match.Type &&
				child.Content(contents) == match.Content {
				fullNode := getDeclarationParent(child)

				results = append(results, fullNode)
			}
		}

		results = searchRecursive(child, results, contents, matches)
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
