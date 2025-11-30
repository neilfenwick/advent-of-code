package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"

	data "github.com/neilfenwick/advent-of-code/data_structures"
)

func main() {
	var (
		file *os.File
		err  error
	)

	switch len(os.Args) {
	case 1:
		file = os.Stdin
	case 3:
		_, _ = strconv.Atoi(os.Args[2])
		fallthrough
	case 2:
		file, err = os.Open(os.Args[1])
		if err != nil {
			log.Fatalf("Error opening file: %s", os.Args[1])
		}
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	/*
	   Because this problem looked like a file tree calculation, I just used a tree data
	   structure from previous years.  In hindsight, that old tree wasn't the best
	   implementation. And it could probably do with being re-written with generics.

	   This solution could probably have been achieved by just passing a hashtable around
	   and doing a recursive parse of the input.

	   The solution feels a bit scatter-brained and unelegant 😕
	*/

	fileTree := processTerminalOutput(file)

	fmt.Printf("*************\nPart 1:\n*************\n")
	results := searchDirectoriesMaxSize(fileTree, 100000)

	totalSize := 0
	for _, size := range results {
		totalSize += size
	}
	fmt.Printf("Total size: %d\n", totalSize)

	fmt.Printf("*************\nPart 2:\n*************\n")
	results = searchDirectoriesMaxSize(fileTree, math.MaxInt)

	rootSize, _ := results["$root"]
	spaceRemaining := 70_000_000 - rootSize
	requiredToFree := 30_000_000 - spaceRemaining

	largeDirectories := make([]directory, 0)
	for name, size := range results {
		if size >= requiredToFree {
			largeDirectories = append(largeDirectories, directory{name: name, size: size})
			fmt.Printf("%s: %d\n", name, size)
		}
	}
	sort.Slice(largeDirectories, func(i, j int) bool {
		return largeDirectories[i].size < largeDirectories[j].size
	})
	fmt.Printf("Directory to delete: %v\n", largeDirectories[0])
}

type termLine struct {
	prefix string
	suffix string
}

type file struct {
	name string
	size int
}

type directory struct {
	name        string
	files       []file
	directories []directory
	size        int
}

func processTerminalOutput(inputFile *os.File) *data.GenericTree[any] {
	s := bufio.NewScanner(inputFile)
	var t *data.GenericTree[any]
	var currentNode *data.GenericTreeNode[any]

	for s.Scan() {
		line := s.Text()

		if line == "$ cd /" {
			log.Println("$root")
			t = data.NewGenericTree(any(directory{name: "$root"}))
			currentNode = t.Root
			continue
		}

		if strings.HasPrefix(line, "$ cd") {
			var dest string
			fmt.Sscanf(line, "$ cd %s", &dest)
			if dest == ".." {
				if currentNode.Parent != nil {
					currentNode = currentNode.Parent
				}
			} else {
				var found *data.GenericTreeNode[any]
				for _, child := range currentNode.Children {
					if dir, ok := child.Value.(directory); ok && dir.name == dest {
						found = child
						break
					}
				}
				if found != nil {
					currentNode = found
				}
			}
		} else if strings.HasPrefix(line, "$ ls") {
			for s.Scan() {
				line := s.Text()
				if strings.HasPrefix(line, "$") {
					// This is a command line, process it in the next iteration
					// But we already consumed it, so we need to handle it now
					if strings.HasPrefix(line, "$ cd") {
						var dest string
						fmt.Sscanf(line, "$ cd %s", &dest)
						if dest == ".." {
							if currentNode.Parent != nil {
								currentNode = currentNode.Parent
							}
						} else {
							var found *data.GenericTreeNode[any]
							for _, child := range currentNode.Children {
								if dir, ok := child.Value.(directory); ok && dir.name == dest {
									found = child
									break
								}
							}
							if found != nil {
								currentNode = found
							}
						}
					}
					break
				}
				parts := strings.Fields(line)
				if len(parts) < 2 {
					continue
				}
				if parts[0] == "dir" {
					log.Printf("Appending directory '%s'\n", parts[1])
					dir := directory{name: parts[1]}
					currentNode.AddChild(any(dir))
				} else {
					size, err := strconv.Atoi(parts[0])
					if err == nil {
						log.Printf("Parsing file: %s, size: %d\n", parts[1], size)
						f := file{name: parts[1], size: size}
						currentNode.AddChild(any(f))
					}
				}
			}
		}
	}

	return t
}

func searchDirectoriesMaxSize(fileTree *data.GenericTree[any], maxSizeThreshold int) map[string]int {
	allDirectorySizes := make(map[string]int)
	root := fileTree.Root
	if rootDir, ok := root.Value.(directory); ok {
		allDirectorySizes[rootDir.name] = 0
	}
	searchDirectoriesRecursive(root, allDirectorySizes)

	results := make(map[string]int)
	for key, value := range allDirectorySizes {
		if value <= maxSizeThreshold {
			results[key] = value
		}
	}

	return results
}

func searchDirectoriesRecursive(treeNode *data.GenericTreeNode[any], directorySizeMap map[string]int) {
	dir, ok := treeNode.Value.(directory)
	if !ok {
		return // skip non-directory nodes
	}

	totalSize := 0
	for _, child := range treeNode.Children {
		switch childVal := child.Value.(type) {
		case directory:
			searchDirectoriesRecursive(child, directorySizeMap)
			if size, ok := directorySizeMap[childVal.name]; ok {
				totalSize += size
			}
		case file:
			totalSize += childVal.size
		}
	}
	directorySizeMap[dir.name] = totalSize
}
