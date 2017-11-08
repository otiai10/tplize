# tplize

`tplize` is a command line to generate Go variable which pools contents of other files.

When you want to embed any template in your cli tool, for example, it must be included in the exacutables,
it means non-Go files are not included in it.

`tplize` solves it __in f**king easy way__ by generating just a map variable.

# Example

Given

```sh
./example
├── index.html
└── messages
    └── hello.txt
```

then hit

```
% cd ./example
% tplize
```

generates

```sh
./example
├── example_templated.go # <- This is generated
├── index.html
└── messages
    └── hello.txt
```

with Go variable like this

```go
package example

// Tpl is auto-generated template pool.
var Tpl = map[string]string{
	// index.html
	"index.html": `...`,         // Here contents of "./example/index.html"
	// messages/hello.txt
	"messages/hello.txt": `...`, // Here contents of "./example/messsages/hello.txt"
}
```

so that you can use the contents in your Go code like this

```go
example.Tpl["index.html"]
```

# Installation

```sh
go get -u github.com/otiai10/tplize
```

# Go generate

`tplize` is just a cli, therefore, it can be hit by `go:generate`.

When you place any Go file in the directory like this

```sh
./example
├── generate.go # <-
├── index.html
└── messages
    └── hello.txt
```

```go
package example

//go:generate tplize
```

then hit

```sh
% go generate ./example
```

would generate `example_templated.go` as well.
