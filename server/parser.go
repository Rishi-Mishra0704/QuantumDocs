package server

import (
	"go/ast"
	"go/parser"
	"go/token"

	"strings"

	"github.com/Rishi-Mishra0704/QuantumDocs/models"
)

func ParseAPIDoc(filePath string) (*models.APIDoc, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	apiDoc := &models.APIDoc{}

	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			if x.Doc != nil {
				endpoint := parseEndpoint(x)
				if endpoint != nil {
					apiDoc.Endpoints = append(apiDoc.Endpoints, *endpoint)
				}
			}
		}
		return true
	})

	return apiDoc, nil
}

func parseEndpoint(fn *ast.FuncDecl) *models.Endpoint {
	endpoint := &models.Endpoint{}
	for _, comment := range fn.Doc.List {
		text := strings.TrimSpace(strings.TrimPrefix(comment.Text, "//"))
		if strings.HasPrefix(text, "@Method") {
			endpoint.Method = strings.TrimSpace(strings.TrimPrefix(text, "@Method"))
		} else if strings.HasPrefix(text, "@Path") {
			endpoint.Path = strings.TrimSpace(strings.TrimPrefix(text, "@Path"))
		} else if strings.HasPrefix(text, "@Description") {
			endpoint.Description = strings.TrimSpace(strings.TrimPrefix(text, "@Description"))
		}
	}

	if endpoint.Method != "" && endpoint.Path != "" {
		return endpoint
	}
	return nil
}
