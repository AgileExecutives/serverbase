package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func findModuleYAMLs(root string) ([]string, error) {
	var out []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.EqualFold(filepath.Base(path), "module.yaml") {
			out = append(out, path)
		}
		return nil
	})
	return out, err
}

func parseNameVersion(r io.Reader) (string, string, error) {
	s := bufio.NewScanner(r)
	var name, version string
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if strings.HasPrefix(line, "name:") {
			name = strings.TrimSpace(strings.TrimPrefix(line, "name:"))
		}
		if strings.HasPrefix(line, "version:") {
			version = strings.TrimSpace(strings.TrimPrefix(line, "version:"))
		}
		if name != "" && version != "" {
			break
		}
	}
	if err := s.Err(); err != nil {
		return "", "", err
	}
	return name, version, nil
}

func cmdModuleList(root string) error {
	files, err := findModuleYAMLs(root)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		fmt.Println("No module.yaml files found")
		return nil
	}
	fmt.Printf("Found %d module descriptors:\n", len(files))
	for _, f := range files {
		r, err := os.Open(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to open %s: %v\n", f, err)
			continue
		}
		name, version, _ := parseNameVersion(r)
		r.Close()
		rel, _ := filepath.Rel(root, filepath.Dir(f))
		if name == "" {
			name = "(unknown)"
		}
		if version == "" {
			version = "(unknown)"
		}
		fmt.Printf("- %s: %s (path: %s)\n", name, version, rel)
	}
	return nil
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: sbctl module list [--root PATH]\n")
		flag.PrintDefaults()
	}
	root := flag.String("root", ".", "project root to scan for modules")
	flag.Parse()
	if flag.NArg() < 1 || flag.Arg(0) != "module" {
		flag.Usage()
		os.Exit(2)
	}
	if flag.NArg() < 2 || flag.Arg(1) != "list" {
		flag.Usage()
		os.Exit(2)
	}
	if err := cmdModuleList(*root); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
