package internal

import (
	"sort"
	"testing"
)

func TestFileSortingByPathname(t *testing.T) {

	path1 := "b.txt"
	path2 := "a/b.txt"
	path3 := "a/a.txt"
	path4 := "a.txt"

	paths := []string{path1, path2, path3, path4}
	correctOrder := []string{path3, path2, path4, path1}

	sort.Sort(ByFilePath(paths))
	for i, _ := range paths {
		pathToTest := paths[i]
		correctPath := correctOrder[i]
		if pathToTest != correctPath {
			t.Fatalf("Incorrect order! i: %v", i)
		}
	}

}
