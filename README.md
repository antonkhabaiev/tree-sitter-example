# tree-sitter-example

Sample application for parsing javascript.

# run
- go to directory root
- make sure you have go 1.17 installed
- `go mod tidy`
- `go run main.go -file-path=parser/javascript_test.js`
- file-path flag determines path for javascript source code to be parsed

# test
- go to directory root
- `go test ./...`