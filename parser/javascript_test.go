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
	ctx := context.Background()
	j := parser.NewJavascript()

	contents, err := ioutil.ReadFile("javascript_unit_test.js")
	if err != nil {
		panic(err)
	}

	tree := j.Parse(ctx, contents)
	Convey("FindReadFile", t, func() {
		Convey("should find all occurrences", func() {
			nodes := j.FindReadFile(tree, contents, "fs")
			So(nodes, ShouldHaveLength, 2)
			So(nodes[0].StartPoint().Row, ShouldEqual, 7)
			So(nodes[0].Content(contents), ShouldEqual, `const f = await fs.readFile("/tmp/somefile.txt");`)
			So(nodes[1].StartPoint().Row, ShouldEqual, 12)
			So(nodes[1].Content(contents), ShouldEqual, `const f2 = await fs.readFile("/tmp/somefile_2.txt");`)
		})
		Convey("should return no results when nothing found", func() {
			nodes := j.FindReadFile(tree, contents, "f")
			So(nodes, ShouldHaveLength, 0)
		})
	})
}

func TestJavascript_FindFsPromisesVarName(t *testing.T) {
	contents, err := ioutil.ReadFile("javascript_unit_test.js")
	if err != nil {
		panic(err)
	}

	Convey("FindFsPromisesVarName", t, func() {
		ctx := context.Background()
		j := parser.NewJavascript()
		Convey("should find first occurrence of var name definition", func() {
			tree := j.Parse(ctx, contents)

			actual := j.FindFsPromisesVarName(tree, contents)
			So(actual, ShouldEqual, "fs")
		})
		Convey("should panic when", func() {
			Convey("no var definition found", func() {
				contents = []byte(`
"use strict";

const app = express();

const router = express.Router();
router.get('/', async (req, res) => {
    res.send("ok");
});
app.use(router);

app.listen(3000, async () => {
    console.log('App listening locally at :3000');
});
`)
				tree := j.Parse(ctx, contents)
				So(func() {
					_ = j.FindFsPromisesVarName(tree, contents)
				}, ShouldPanicWith,
					"couldn't find any occurrences of fs.promises")
			})
		})
	})
}
