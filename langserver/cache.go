package langserver

import (
	"sync"

	"github.com/sourcegraph/go-lsp"
)

type cache struct {
	mux   *sync.RWMutex
	files map[lsp.DocumentURI]lsp.TextDocumentItem
}

func newCache() *cache {
	return &cache{
		mux:   &sync.RWMutex{},
		files: make(map[lsp.DocumentURI]lsp.TextDocumentItem),
	}
}

func (c *cache) Set(uri lsp.DocumentURI, item lsp.TextDocumentItem) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.files[uri] = item
}

func (c *cache) Get(uri lsp.DocumentURI) (lsp.TextDocumentItem, bool) {
	c.mux.RLock()
	defer c.mux.RUnlock()
	item, ok := c.files[uri]
	return item, ok
}

func (c *cache) Delete(uri lsp.DocumentURI) {
	c.mux.Lock()
	defer c.mux.Unlock()
	delete(c.files, uri)
}
