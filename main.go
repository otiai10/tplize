package main

import (
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	if err := Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Run ...
func Run(args []string) error {

	gendir, err := os.Getwd()
	if err != nil {
		return err
	}
	// self := filepath.Join(gendir, os.Getenv("GOFILE"))

	parsed, err := ParseArgs(args[1:])
	if err != nil {
		return err
	}

	if parsed.Path != "" {
		gendir = filepath.Join(gendir, parsed.Path)
	}

	files, err := parsed.ToTargetFiles(gendir)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("no files found on this directory: %s", gendir)
	}

	pkg, err := build.ImportDir(gendir, 0)
	if err != nil {
		return err
	}

	gen := &Generator{
		FileNames: files,
		Pkg:       pkg,
		VarName:   parsed.VarName,
		Output:    parsed.Output,
	}

	if _, err = os.Stat(gen.Destination()); err == nil {
		// TODO: Refactor
		if err = os.Remove(gen.Destination()); err != nil {
			return err
		}
	}
	f, err := os.Create(gen.Destination())
	if err != nil {
		return err
	}
	defer f.Close()

	if err = gen.Stack(); err != nil {
		return err
	}
	if err = gen.Flush(f); err != nil {
		return err
	}

	return nil
}

// Args ...
type Args struct {
	Pattern *string
	Files   []string
	Path    string
	VarName string
	Output  string
}

// ToTargetFiles ...
func (args *Args) ToTargetFiles(dir string) ([]string, error) {

	files := []string{}

	// Files specified by names
	for _, fname := range args.Files {
		files = append(files, filepath.Join(dir, fname))
	}
	if len(files) != 0 {
		return files, nil
	}

	// Files specified by pattern
	if args.Pattern != nil {
		matches, err := filepath.Glob(*args.Pattern)
		if err != nil {
			return files, err
		}
		for _, fname := range matches {
			files = append(files, filepath.Join(dir, fname))
		}
		return files, nil
	}

	// Implicitly specified as current directory
	info, err := ioutil.ReadDir(dir)
	if err != nil {
		return files, err
	}
	for _, f := range info {
		name := f.Name()
		if filepath.Ext(name) != ".go" {
			files = append(files, filepath.Join(dir, name))
		}
	}

	return files, nil
}

// ParseArgs ...
func ParseArgs(args []string) (*Args, error) {
	parsed := new(Args)
	var arg string
	for {
		if len(args) == 0 {
			break
		}
		arg, args = args[0], args[1:]
		switch arg {
		case "-e", "--regex":
			if len(args) == 0 {
				continue
			}
			parsed.Pattern = &args[0]
			args = args[1:] // Unshift again
		case "-V", "--var":
			if len(args) == 0 {
				continue
			}
			parsed.VarName = args[0]
			args = args[1:] // Unshift again
		case "-p", "--path":
			if len(args) == 0 {
				continue
			}
			parsed.Path = args[0]
			args = args[1:] // Unshift again
		case "-o", "--out":
			if len(args) == 0 {
				continue
			}
			parsed.Output = args[0]
			args = args[1:] // Unshift again
		default:
			parsed.Files = append(parsed.Files, arg)
		}
	}
	return parsed, nil
}
