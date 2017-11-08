package main

import (
	"bytes"
	"fmt"
	"go/build"
	"go/format"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Generator ...
type Generator struct {
	FileNames []string
	Pkg       *build.Package
	Files     map[string]string
	VarName   string
	Output    string
	buf       *bytes.Buffer
}

// Destination ...
func (g *Generator) Destination() string {
	if g.Output != "" {
		return g.Output
	}
	return filepath.Join(g.Pkg.Dir, fmt.Sprintf("%s_templated.go", g.Pkg.Name))
}

// Stack stacks contents of target files.
func (g *Generator) Stack() error {
	if g.Files == nil {
		g.Files = make(map[string]string)
	}
	for _, entry := range g.FileNames {
		if err := g.stack(entry); err != nil {
			return err
		}
	}
	return nil
}

// read
func (g *Generator) stack(entry string) error {

	info, err := os.Stat(entry)
	if err != nil {
		return err
	}

	if info.IsDir() {
		matches, err := filepath.Glob(filepath.Join(entry, "*"))
		if err != nil {
			return err
		}
		for _, e := range matches {
			if err = g.stack(e); err != nil {
				return err
			}
		}
	} else {
		f, err := os.Open(entry)
		if err != nil {
			return err
		}
		defer f.Close()
		b, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}
		g.Files[entry] = string(b)
	}

	return nil
}

// Flush write all stacked contents to destination template file.
func (g *Generator) Flush(w io.Writer) (err error) {

	if g.VarName == "" {
		g.VarName = "Tpl"
	}
	if g.buf == nil {
		g.buf = bytes.NewBuffer(nil)
	}

	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("failed to flush stacked contents to writer: %v", e)
		}
	}()

	g.write("package %s\n\n", g.Pkg.Name)

	g.write("// %s is auto-generated template pool.\n", g.VarName)
	g.write("var %s = map[string]string{\n", g.VarName)
	for name, content := range g.Files {
		name = strings.Replace(name, g.Pkg.Dir+"/", "", 1)
		g.write("// %s\n", name)
		g.write("\"%s\": `%s`,\n", name, content)
	}
	g.write("}\n")

	formatted, err := format.Source(g.buf.Bytes())
	if err != nil {
		return err
	}

	_, err = w.Write(formatted)
	return err
}

func (g *Generator) write(format string, v ...interface{}) {
	_, err := fmt.Fprintf(g.buf, format, v...)
	if err != nil {
		panic(err)
	}
}
