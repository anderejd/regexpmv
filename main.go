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

func getMatchingPaths(root string, re *regexp.Regexp) ([]string, error) {
	paths := []string{}
	walker := makeWalkFunc(func(path string) error {
		if !re.MatchString(path) {
			return nil
		}
		paths = append(paths, path)
		return nil
	})
	return paths, filepath.Walk(root, walker)
}

func main() {
	if len(os.Args) < 4 {
		log.Fatalln(`Syntax error. Usage: regexpmv ./ "text" "new $1"`)
	}
	root := os.Args[1]
	re, err := regexp.Compile(os.Args[2])
	if err != nil {
		log.Fatalln(err)
	}
	newText := os.Args[3]
	paths, err := getMatchingPaths(root, re)
	if err != nil {
		log.Fatalln(err)
	}
	for _, path := range paths {
		newPath := re.ReplaceAllString(path, newText)
		err := os.Rename(path, newPath)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
