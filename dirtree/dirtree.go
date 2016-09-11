package dirtree

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type NodeFile struct {
	Modified time.Time `json:",omitempty"`
	Path     string
}
type Node struct {
	Childs *Nodes    `json:",omitempty"`
	File   *NodeFile `json:",omitempty"`
	//TODO: Do we need name?
}

type Nodes map[string]Node

type RegisterFunc func(file os.FileInfo) error

type Tree struct {
	Root         *Node
	Prefix       string
	Ext          string
	Log          *log.Logger
	RegisterFunc *RegisterFunc
}

// -----------------------------------------------------------------------------
// Functional options

// Debug sets sql tracing to on when "on" argument is true
func OnRegister(f RegisterFunc) func(tree *Tree) error {
	return func(tree *Tree) error {
		return tree.setRegisterFunc(f)
	}
}

// -----------------------------------------------------------------------------
// Internal setters

func (tree *Tree) setRegisterFunc(f RegisterFunc) error {
	tree.RegisterFunc = &f
	return nil
}

// -----------------------------------------------------------------------------

func New(logger *log.Logger, root, ext string, options ...func(tree *Tree) error) (*Tree, error) {
	tree := Tree{Log: logger, Prefix: root, Ext: ext, Root: &Node{}}
	for _, option := range options {
		err := option(&tree)
		if err != nil {
			return nil, err
		}
	}
	tree.ReadLevel(root, nil)
	return &tree, nil
}

func (tree Tree) ReadLevel(root string, level *Node) {

	//	log.Printf("Scan dir %s", root)
	files, err := ioutil.ReadDir(root)
	if err != nil {
		tree.Log.Fatalf("Read dir (%s) error: %+v", root, err)
	}
	if level == nil {
		level = tree.Root
	}
	nodes := Nodes{}
	for _, file := range files {
		name := file.Name()
		node := Node{}
		if file.IsDir() {
			tree.ReadLevel(filepath.Join(root, name), &node)
			if node.Childs != nil || node.File != nil {
				nodes[name] = node
			}
		} else if strings.HasSuffix(name, tree.Ext) {
			tree.Log.Printf("debug: Add file %s %s", root, name)
			path0 := strings.TrimPrefix(filepath.Join(root, name), tree.Prefix)
			f := NodeFile{Path: path0}
			name0 := strings.TrimSuffix(name, tree.Ext)
			if tree.RegisterFunc != nil {
				fu := tree.RegisterFunc
				(*fu)(file)
			}
			if name0 == "index" {
				// node describes level
				level.File = &f
			} else if _, ok := nodes[name0]; !ok { // if no such dir already
				node.File = &f
				nodes[name0] = node
			}
		}
	}
	//	log.Printf("Level: %+v", nodes)
	if len(nodes) > 0 {
		level.Childs = &nodes
	}
}

func (tree Tree) Node(path string, root *Node) (*Node, bool) {

	var current *Node
	if root == nil {
		current = tree.Root
	} else {
		current = root
	}
	if path == "" {
		return current, current.File != nil
	}
	nodeList := []string{}
	dirList := strings.Split(strings.TrimSuffix(path, "/"), "/")
	for _, d := range dirList {
		if current.Childs == nil {
			tree.Log.Printf("deug: no childs: %s", d)
			break
		}
		if node, ok := (*current.Childs)[d]; ok {
			nodeList = append(nodeList, d)
			current = &node
		} else {
			tree.Log.Printf("debug: no node: %s", d)
			break
		}
	}
	found := len(nodeList) == len(dirList) && current.File != nil
	return current, found
}
