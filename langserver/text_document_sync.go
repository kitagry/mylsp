package langserver

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sourcegraph/go-lsp"
	"github.com/sourcegraph/jsonrpc2"
)

func (h *handler) handleTextDocumentDidOpen(
	_ context.Context,
	_ *jsonrpc2.Conn,
	req *jsonrpc2.Request,
) (result any, err error) {
	var params lsp.DidOpenTextDocumentParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}

	h.cache.Set(params.TextDocument.URI, params.TextDocument)

	h.diagnosticRequest <- params.TextDocument.URI

	return nil, nil
}

func (h *handler) handleTextDocumentDidChange(
	_ context.Context,
	_ *jsonrpc2.Conn,
	req *jsonrpc2.Request,
) (result any, err error) {
	var params lsp.DidChangeTextDocumentParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}

	file, ok := h.cache.Get(params.TextDocument.URI)
	if !ok {
		err = fmt.Errorf("document not found: %s", params.TextDocument.URI)
		return nil, err
	}
	file.Text = params.ContentChanges[0].Text
	file.Version = params.TextDocument.Version
	h.cache.Set(params.TextDocument.URI, file)

	h.diagnosticRequest <- params.TextDocument.URI

	return nil, nil
}

func (h *handler) handleTextDocumentDidClose(
	_ context.Context,
	_ *jsonrpc2.Conn,
	req *jsonrpc2.Request,
) (result any, err error) {
	var params lsp.DidCloseTextDocumentParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}

	h.cache.Delete(params.TextDocument.URI)

	return nil, nil
}

func (h *handler) handleTextDocumentDidSave(
	_ context.Context,
	_ *jsonrpc2.Conn,
	req *jsonrpc2.Request,
) (result any, err error) {
	var params lsp.DidSaveTextDocumentParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}

	return nil, nil
}
