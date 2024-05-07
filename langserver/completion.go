package langserver

import (
	"context"
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"log"

	"github.com/sourcegraph/go-lsp"
	"github.com/sourcegraph/jsonrpc2"
)

func (h *handler) handleTextDocumentCompletion(
	_ context.Context,
	_ *jsonrpc2.Conn,
	req *jsonrpc2.Request,
) (any, error) {
	var params lsp.TextDocumentPositionParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}

	file, ok := h.cache.Get(params.TextDocument.URI)
	if !ok {
		return []lsp.CompletionItem{}, nil
	}

	parsed, err := parser.ParseFile(token.NewFileSet(), string(params.TextDocument.URI), file.Text, parser.AllErrors)
	if err != nil {
		return nil, err
	}

	funcs := listFunctions(parsed)
	result := make([]lsp.CompletionItem, len(funcs))
	for i, f := range funcs {
		doc := file.Text[f.Pos()-1 : f.End()-1]
		if f.Doc != nil {
			log.Println(f.Doc.Text())
			doc = f.Doc.Text()
		}
		result[i] = lsp.CompletionItem{
			Label:            f.Name.Name,
			Kind:             lsp.CIKFunction,
			Detail:           doc,
			InsertTextFormat: lsp.ITFPlainText,
			InsertText:       f.Name.Name + "()",
		}
	}

	return result, nil
}

func listFunctions(f *ast.File) []*ast.FuncDecl {
	results := make([]*ast.FuncDecl, 0)
	for _, d := range f.Decls {
		switch d := d.(type) {
		case *ast.FuncDecl:
			results = append(results, d)
		}
	}
	return results
}
