package dirtree

import (
	"testing"

	"github.com/Sirupsen/logrus"
)

type testItem struct {
	Success bool
	URL     string
	Path    string
}

var tests = []testItem{
	{true, "", "index.md"},
	{true, "dir1", "dir1/index.md"},
	{true, "dir1/file2", "dir1/file2.md"},
	{true, "dir1/dir2/file4", "dir1/dir2/file4.md"},
	{false, "dir1/dir2", ""}, // dir1/dir2/index.md - path exists but no index
	{false, "dir1/dXX/file4", "dir1/index.md"},
}

func pathOrEmpty(n *Node) (str string) {
	if n != nil && n.File != nil {
		str = n.File.Path
	}
	return
}

func TestDirTree(t *testing.T) {

	log := logrus.New().WithField("in", "test")
	tree, _ := New(log, "test/", ".md")

	for _, tt := range tests {
		node, ok := tree.Node(tt.URL, nil)
		if ok != tt.Success || tt.Path != pathOrEmpty(node) {
			t.Errorf("%s: expected (%v) %s, actual (%v) %s", tt.URL, tt.Success, tt.Path, ok, pathOrEmpty(node))
		}
	}
}
