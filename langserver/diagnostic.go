package langserver

import (
	"context"
	"errors"
	"fmt"
	"go/parser"
	"go/scanner"
	"go/token"
	"log"

	"github.com/sourcegraph/go-lsp"
)

func (h *handler) handleDiagnostics() {
	for uri := range h.diagnosticRequest {
		ctx := context.Background()

		go func() {
			diagnostics, err := h.diagnose(uri)
			if err != nil {
				log.Println(err.Error())
				return
			}

			h.conn.Notify(
				ctx,
				"textDocument/publishDiagnostics",
				lsp.PublishDiagnosticsParams{
					URI:         uri,
					Diagnostics: diagnostics,
				},
			)
		}()
	}
}

func (h *handler) diagnose(uri lsp.DocumentURI) ([]lsp.Diagnostic, error) {
	file, ok := h.cache.Get(uri)
	if !ok {
		return nil, fmt.Errorf("document not found: %s", uri)
	}

	_, err := parser.ParseFile(token.NewFileSet(), string(uri), file.Text, parser.AllErrors)
	if err == nil {
		return []lsp.Diagnostic{}, nil
	}

	var errList scanner.ErrorList
	if ok := errors.As(err, &errList); !ok {
		return nil, err
	}

	result := convertErrorListToDiagnostics(errList)
	return result, nil
}

func convertErrorListToDiagnostics(errs scanner.ErrorList) []lsp.Diagnostic {
	result := make([]lsp.Diagnostic, len(errs))
	for i, e := range errs {
		result[i] = convertErrorToDiagnostic(e)
	}
	return result
}

func convertErrorToDiagnostic(err *scanner.Error) lsp.Diagnostic {
	return lsp.Diagnostic{
		Severity: lsp.Error,
		Range: lsp.Range{
			Start: lsp.Position{
				Line:      err.Pos.Line - 1,
				Character: err.Pos.Column - 1,
			},
			End: lsp.Position{
				Line:      err.Pos.Line - 1,
				Character: err.Pos.Column - 1,
			},
		},
		Message: err.Msg,
	}
}
