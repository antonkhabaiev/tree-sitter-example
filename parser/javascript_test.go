package parser_test

import (
	"context"
	"io/ioutil"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"tree-sitter-example/parser"
)

func TestNewJavascript(t *testing.T) {
	Convey("NewJavascript should return new object", t, func() {
		j := parser.NewJavascript()

		So(j, ShouldNotBeNil)
	})
}

func TestJavascript_FindReadFile(t *testing.T) {
	Convey("FindReadFile", t, func() {
		ctx := context.Background()
		j := parser.NewJavascript()

		contents, err := ioutil.ReadFile("javascript_unit_test.js")
		if err != nil {
			panic(err)
		}

		Convey("should find all occurrences", func() {
			nodes := j.FindReadFile(ctx, contents, "fs")
			So(nodes, ShouldHaveLength, 2)
			So(nodes[0].StartPoint().Row, ShouldEqual, 7)
			So(nodes[0].Content(contents), ShouldEqual, `const f = await fs.readFile("/tmp/somefile.txt");`)
			So(nodes[1].StartPoint().Row, ShouldEqual, 12)
			So(nodes[1].Content(contents), ShouldEqual, `const f2 = await fs.readFile("/tmp/somefile_2.txt");`)
		})
		Convey("should return no results when nothing found", func() {
			nodes := j.FindReadFile(ctx, contents, "f")
			So(nodes, ShouldHaveLength, 0)
		})
	})
}

func TestJavascript_FindFsPromisesVarName(t *testing.T) {
	Convey("FindFsPromisesVarName", t, func() {
		ctx := context.Background()
		j := parser.NewJavascript()

		contents, err := ioutil.ReadFile("javascript_unit_test.js")
		if err != nil {
			panic(err)
		}
		Convey("should find first occurrence of var name definition", func() {
			actual := j.FindFsPromisesVarName(ctx, contents)
			So(actual, ShouldEqual, "fs")
		})
		// Convey("")
	})
}
