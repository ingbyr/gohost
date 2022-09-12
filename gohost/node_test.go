package gohost

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"testing"
)

func TestService_LoadTree(t *testing.T) {
	svc := GetService()
	//printTree(svc.Tree().children)
	svc.Node(001).isFolded = false
	printNodes(svc.TreeNodeItem())
}

func printTree(nodes []*TreeNode) {
	if len(nodes) == 0 {
		return
	}
	for _, node := range nodes {
		//fmt.Println(strings.Repeat(" ", node.depth) + node.Title())
		fmt.Println(node.Title())
		printTree(node.Children())
	}
}

func printNodes(nodes []list.Item) {
	for _, node := range nodes {
		fmt.Println(node.FilterValue())
	}
}
