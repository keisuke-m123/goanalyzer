package gocode

import "go/types"

// golangのbuiltin型であるかを判定する
func builtin(name string) bool {
	return types.Universe.Lookup(name) != nil
}
