package main

import "log"
import "os"
import "path/filepath"
import "regexp"

type fileHandler func(path string) error

func isDotPath(p string) bool {
	b := filepath.Base(p)
	if ".." != b && len(b) > 1 && '.' == b[0] {
		return true
	}
	return false
}

func makeWalkFunc(fh fileHandler) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if isDotPath(path) {
				return filepath.SkipDir
			}
			return nil
		}
		return fh(path)
	}
}

func main() {
	if len(os.Args) < 4 {
		log.Fatalln(`Syntax error. Usage: regexpmv ./ "text" "new $1"`)
	}
	root := os.Args[1]
	r, err := regexp.Compile(os.Args[2])
	if err != nil {
		log.Fatalln(err)
	}
	newText := os.Args[3]
	walker := makeWalkFunc(func(path string) error {
		if !r.MatchString(path) {
			return nil
		}
		newPath := r.ReplaceAllString(path, newText)
		return os.Rename(path, newPath)
	})
	err = filepath.Walk(root, walker)
	if err != nil {
		log.Fatalln(err)
	}
}
