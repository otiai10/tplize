package main

import (
	"fmt"
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

	files, err := parsed.ToTargetFiles(gendir)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("no files found on this directory: %s", gendir)
	}

	return nil
}

// Args ...
type Args struct {
	Pattern *string
	Files   []string
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
		default:
			parsed.Files = append(parsed.Files, arg)
		}
	}
	return parsed, nil
}
