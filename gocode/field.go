package gocode

import (
	"go/types"
)

type (
	// FieldName はフィールド名を表す。
	FieldName string

	// Field はstructのフィールドを表す。
	Field struct {
		goVar      *types.Var
		name       FieldName
		pkgSummary *PackageSummary
		typ        *Type
	}

	// FieldList はstructのフィールドのリストを表す。
	FieldList struct {
		fields []*Field
	}
)

func (fn FieldName) String() string {
	return string(fn)
}

func newFieldListFromStructType(structType *types.Struct) *FieldList {
	var fields []*Field
	for i := 0; i < structType.NumFields(); i++ {
		f := structType.Field(i)
		fields = append(fields, newField(f))
	}
	return &FieldList{fields: fields}
}

func (fl *FieldList) asSlice() []*Field {
	var slice []*Field
	for i := range fl.fields {
		slice = append(slice, fl.fields[i])
	}
	return slice
}

func newField(field *types.Var) *Field {
	pkgSummary := newPackageSummaryFromGoTypes(field.Pkg())

	return &Field{
		goVar:      field,
		pkgSummary: pkgSummary,
		name:       FieldName(field.Name()),
		typ:        newType(pkgSummary, field.Type()),
	}
}

func (f *Field) Embedded() bool {
	return f.goVar.Embedded()
}

func (f *Field) PackageSummary() *PackageSummary {
	return f.pkgSummary
}

func (f *Field) Name() FieldName {
	return f.name
}

func (f *Field) Type() *Type {
	return f.typ
}
