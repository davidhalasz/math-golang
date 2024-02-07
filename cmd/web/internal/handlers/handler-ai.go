package handlers

import (
	"fmt"
	"net/http"

	"github.com/davidhalasz/gomath/cmd/web/internal/render"
)

func AiPage(w http.ResponseWriter, r *http.Request) {
	if err := render.Template(w, r, "ai-basics.page.gohtml", nil); err != nil {
		app.ErrorLog.Println(err)
	}
}

type Node struct {
	Val   int
	Left  *Node
	Right *Node
	Next  *Node
}

type NodeWithLevel struct {
	node  *Node
	level int
}

func BFS(root *Node) *Node {
	if root == nil {
		return root
	}

	queue := []NodeWithLevel{
		{
			node:  root,
			level: 0,
		},
	}

	visited := []*Node{}

	for len(queue) > 0 {
		vertex := queue[0]
		node, level := vertex.node, vertex.level
		visited = append(visited, node)
		queue = queue[1:]

		if node.Left != nil {
			leftNode := NodeWithLevel{
				node:  node.Left,
				level: level + 1,
			}
			queue = append(queue, leftNode)

			rightNode := NodeWithLevel{
				node:  node.Right,
				level: level + 1,
			}
			queue = append(queue, rightNode)
		}

		fmt.Printf("Level: %d", queue[len(queue)-1].level)
		fmt.Printf("Visited: [ ")
		for _, node := range visited {
			fmt.Printf("%d ", node.Val)
		}
		fmt.Printf("] Queue: [ ")
		for _, vertex := range queue {
			fmt.Printf("%d ", vertex.node.Val)
		}
		fmt.Printf("]\n")
	}

	return root
}

func CallBFS() {
	root := &Node{Val: 1}
	root.Left = &Node{Val: 2}
	root.Right = &Node{Val: 3}
	root.Left.Left = &Node{Val: 4}
	root.Left.Right = &Node{Val: 5}
	root.Right.Left = &Node{Val: 6}
	root.Right.Right = &Node{Val: 7}

	fmt.Println("BFS traversal of the binary tree:")
	BFS(root)

	// The result
	// BFS traversal of the binary tree:
	// Visited: [ 1 ] Queue: [ 2 3 ]
	// Visited: [ 1 2 ] Queue: [ 3 4 5 ]
	// Visited: [ 1 2 3 ] Queue: [ 4 5 6 7 ]
	// Visited: [ 1 2 3 4 ] Queue: [ 5 6 7 ]
	// Visited: [ 1 2 3 4 5 ] Queue: [ 6 7 ]
	// Visited: [ 1 2 3 4 5 6 ] Queue: [ 7 ]
	// Visited: [ 1 2 3 4 5 6 7 ] Queue: [ ]
}

// Deep-First Search
func DFS(root *Node) []*Node {
	visited := []*Node{}

	if root == nil {
		return visited
	}

	return recurse(root, visited)
}

func recurse(root *Node, visited []*Node) []*Node {
	visited = append(visited, root)
	fmt.Printf("visited: [")
	for i, v := range visited {
		fmt.Printf(" %d", v.Val)
		if i < len(visited)-1 {
			fmt.Print(",")
		}
	}
	fmt.Print(" ]\n")

	if root.Left != nil {
		visited = recurse(root.Left, visited)
	}

	if root.Right != nil {
		visited = recurse(root.Right, visited)
	}

	return visited
}

func CallDFS(w http.ResponseWriter, r *http.Request) {
	root := &Node{Val: 1}
	root.Left = &Node{Val: 2}
	root.Right = &Node{Val: 3}
	root.Left.Left = &Node{Val: 4}
	root.Left.Right = &Node{Val: 5}
	root.Right.Left = &Node{Val: 6}
	root.Right.Right = &Node{Val: 7}

	fmt.Println("DFS traversal of the binary tree:")
	DFS(root)

	// The result
	// DFS traversal of the binary tree:
	// visited: [ 1 ]
	// visited: [ 1, 2 ]
	// visited: [ 1, 2, 4 ]
	// visited: [ 1, 2, 4, 5 ]
	// visited: [ 1, 2, 4, 5, 3 ]
	// visited: [ 1, 2, 4, 5, 3, 6 ]
	// visited: [ 1, 2, 4, 5, 3, 6, 7 ]
}
