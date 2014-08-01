package handlers

import (
	"github.com/russmack/elrepl/types"
)

var HandlerRegistry = make(map[string]*types.Handler)
