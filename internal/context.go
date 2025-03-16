package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
)

func GetProjectContext() (string, error) {
	root, err := FindGitRoot()
	if err != nil {
		return "", err
	}

	files, err := getAllFiles(root)
	if err != nil {
		return "", fmt.Errorf("Can't file the files\n")
	}

	context := ""
	for _, file := range files {
		if isBinaryOrLargeFile(filepath.Join(root, file)) {
			continue
		}

		data, err := os.ReadFile(filepath.Join(root, file))
		if err != nil {
			continue
		}
		context += fmt.Sprintf("\n--- %s ---\n%s\n", file, string(data))
	}

	return context, nil
}

func getAllFiles(root string) ([]string, error) {
	cmd := exec.Command("git", "-C", root, "ls-files")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var files []string
	for file := range strings.SplitSeq(strings.TrimSpace(string(output)), "\n") {
		if file != "" {
			files = append(files, file)
		}
	}

	return files, nil
}

func isBinaryOrLargeFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil || info.Size() > 1024*1024 {
		return true
	}

	ext := strings.ToLower(filepath.Ext(path))
	binaryExts := []string{
		".png", ".jpg", ".jpeg", ".gif", ".bmp", ".svg",
		".zip", ".tar", ".gz", ".rar",
		".exe", ".dll", ".so", ".dylib",
		".pdf", ".doc", ".docx",
	}

	return slices.Contains(binaryExts, ext)
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
