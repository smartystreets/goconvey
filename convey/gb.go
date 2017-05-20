package convey

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func SetGBPath(workDir string) {
	root, err := findProjectroot(workDir)
	if err != nil {
		log.Fatal(err)
	}

	goPath := strings.Join([]string{
		root,
		filepath.Join(root, "vendor")}, ":")

	err = os.Setenv("GOPATH", goPath)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Enabling GB compatibility: GOPATH set to '%s'\n", goPath)
}

// findProjectroot works upwards from path seaching for the
// src/ directory which identifies the project root.
func findProjectroot(path string) (string, error) {
	if path == "" {
		return "", errors.New("project root is blank")
	}
	start := path
	for path != filepath.Dir(path) {
		root := filepath.Join(path, "src")
		if _, err := os.Stat(root); err != nil {
			if os.IsNotExist(err) {
				path = filepath.Dir(path)
				continue
			}
			return "", err
		}
		return path, nil
	}
	return "", fmt.Errorf(`could not find project root in "%s" or its parents`, start)
}
