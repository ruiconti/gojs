package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"

	"github.com/awalterschulze/gographviz"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <filename>\n", os.Args[0])
		os.Exit(1)
	}

	filename := os.Args[1]

	// Parse the source file and build the AST
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing file %s: %v\n", filename, err)
		os.Exit(1)
	}

	// Generate a graph of the AST
	graph := gographviz.NewGraph()
	if err := graph.SetName("AST"); err != nil {
		fmt.Fprintf(os.Stderr, "error setting graph name: %v\n", err)
		os.Exit(1)
	}
	if err := graph.SetDir(true); err != nil {
		fmt.Fprintf(os.Stderr, "error setting graph direction: %v\n", err)
		os.Exit(1)
	}

	// Visit each Expr in the AST and add it to the graph
	visitor := &visitor{graph: graph}
	ast.Walk(visitor, node)

	// Output the graph as a DOT file
	dot := graph.String()

	if err := ioutil.WriteFile("ast.dot", []byte(dot), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "error writing DOT file: %v\n", err)
		os.Exit(1)
	}

	// Convert the DOT file to a PNG image using Graphviz
	cmd := exec.Command("dot", "-Tpng", "-o", "ast.png", "ast.dot")
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error rendering graph: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("AST visualization generated at ast.png")
}

type visitor struct {
	graph     *gographviz.Graph
	parent    string
	nodeCount int
}

func encodeNode(node ast.Node) (string, bool) {
	var name string

	switch n := node.(type) {
	case *ast.FuncDecl:
		name = n.Name.String()
	case *ast.ExprStmt:
		if callExpr, isCallExpr := n.X.(*ast.CallExpr); isCallExpr {
			switch fn := callExpr.Fun.(type) {
			case *ast.SelectorExpr:
				name = fmt.Sprintf("%s", fn.Sel.String())
			case *ast.Ident:
				name = fmt.Sprintf("%s", fn.String())
			}
		} else {
			name = fmt.Sprintf(`ExprStmt%v`, n.Pos())
		}
	case *ast.FuncLit:
		name = fmt.Sprintf(`FuncLit%v`, n.Pos())
	case *ast.CallExpr:
		switch fn := n.Fun.(type) {
		case *ast.SelectorExpr:
			name = fmt.Sprintf("%s", fn.Sel.String())
		case *ast.Ident:
			name = fmt.Sprintf("%s", fn.String())
		}
	default:
		return "", false
	}

	reg, err := regexp.Compile(`[\*\.]`)
	if err != nil {
		panic(err)
	} else {
		name = reg.ReplaceAllString(name, "")
	}

	return name, true
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	switch node.(type) {
	case *ast.File:
	case *ast.ExprStmt:
		// fmt.Printf("exprstmt: %s\n", encodeExpr(Expr))
	case *ast.FuncDecl:
	case *ast.CallExpr:
		// fmt.Printf("callexpr: %s\n", encodeExpr(Expr))
	case *ast.FuncLit:
	}

	var prevParent string
	v.nodeCount++
	// Create a new Node for the current AST Node
	if ExprName, ok := encodeNode(node); ok {
		if err := v.graph.AddNode("G", ExprName, map[string]string{
			"label": ExprName,
		}); err != nil {
			fmt.Fprintf(os.Stderr, "error adding Expr: %v\n", err)
			return nil
		}

		// Connect the new Expr to its parent
		if v.parent != "" {
			if err := v.graph.AddEdge(v.parent, ExprName, true, nil); err != nil {
				fmt.Fprintf(os.Stderr, "error adding edge: %v\n", err)
				return nil
			}
		}

		// Save the new Expr as the parent for the next Expr
		prevParent = v.parent
		v.parent = ExprName

		// Return a new visitor to visit the children of this Expr
	}
	return &visitor{
		graph:     v.graph,
		parent:    prevParent,
		nodeCount: v.nodeCount,
	}
}

func (v *visitor) Leave(Expr ast.Expr) {
	// Pop the current Expr off the stack by restoring the parent
	v.parent = ""
}
