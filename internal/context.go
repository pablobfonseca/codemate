package internal

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetProjectContext() (string, error) {
	root, err := FindGitRoot()
	if err != nil {
		return "", err
	}

	files := []string{"main.go", "go.mod", "README.md"}
	context := ""
	for _, file := range files {
		path := filepath.Join(root, file)
		if _, err := os.Stat(path); err == nil {
			data, _ := os.ReadFile(path)
			context += fmt.Sprintf("\n--- %s ---\n%s\n", file, string(data))
		}
	}

	return context, nil
}

func FindGitRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}

		dir = parent
	}

	return "", fmt.Errorf("no Git repository found")
}
