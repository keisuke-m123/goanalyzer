package gocode

import "go/types"

type (
	// Embed は埋め込みされた型を表現する。
	Embed struct {
		typ *Type
	}

	// EmbedList は埋め込みされた型のリストを表現する。
	EmbedList struct {
		embeds []*Embed
	}
)

func newEmbed(currentPkgSummary *PackageSummary, typ types.Type) *Embed {
	return &Embed{typ: newType(currentPkgSummary, typ)}
}

func (e *Embed) Type() *Type {
	return e.typ
}

func newEmbedListFromInterfaceType(currentPkgSummary *PackageSummary, interfaceType *types.Interface) *EmbedList {
	var embeds []*Embed
	for i := 0; i < interfaceType.NumEmbeddeds(); i++ {
		e := interfaceType.EmbeddedType(i)
		embeds = append(embeds, newEmbed(currentPkgSummary, e))
	}
	return &EmbedList{embeds: embeds}
}

func (el *EmbedList) asSlice() []*Embed {
	var slice []*Embed
	for i := range el.embeds {
		slice = append(slice, el.embeds[i])
	}
	return slice
}
