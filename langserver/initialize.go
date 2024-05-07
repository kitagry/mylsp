package langserver

import (
	"context"

	"github.com/sourcegraph/go-lsp"
	"github.com/sourcegraph/jsonrpc2"
)

func (h *handler) handleInitialize(
	_ context.Context,
	conn *jsonrpc2.Conn,
	_ *jsonrpc2.Request,
) (result any,
	err error,
) {
	h.conn = conn

	return lsp.InitializeResult{
		Capabilities: lsp.ServerCapabilities{
			TextDocumentSync: &lsp.TextDocumentSyncOptionsOrKind{
				Kind: toPtr(lsp.TDSKFull),
			},
			CompletionProvider: &lsp.CompletionOptions{
				TriggerCharacters: []string{"."},
			},
		},
	}, nil
}

func toPtr[T any](t T) *T {
	return &t
}
