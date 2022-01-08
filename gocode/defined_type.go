package gocode

import (
	"go/types"

	"golang.org/x/tools/go/packages"
)

type (
	// DefinedTypeName はdefined typeの名前を表す。
	DefinedTypeName string

	// DefinedType はdefined typeを表す。
	DefinedType struct {
		// name はtypeされた名前。
		name DefinedTypeName
		// pkgSummary は定義されているパッケージのサマリ情報。
		pkgSummary *PackageSummary
		// typ はtypeされた型。
		typ *Type
		// methods は定義されたメソッドの一覧。
		methods *FunctionList
	}

	// DefinedTypeList はdefined typeの一覧を表す。
	DefinedTypeList struct {
		definedTypes []*DefinedType
	}
)

func (dtn DefinedTypeName) String() string {
	return string(dtn)
}

func newDefinedTypeIfObjectDefinedType(pkg *packages.Package, obj types.Object) (res *DefinedType, ok bool) {
	tn, ok := obj.(*types.TypeName)
	if !ok || tn.IsAlias() {
		return &DefinedType{}, false
	}

	switch obj.Type().Underlying().(type) {
	case *types.Struct, *types.Interface:
		return &DefinedType{}, false
	}

	pkgSummary := newPackageSummaryFromGoTypes(obj.Pkg())

	return &DefinedType{
		pkgSummary: pkgSummary,
		name:       DefinedTypeName(obj.Name()),
		typ:        newType(pkgSummary, obj.Type().Underlying()),
		methods:    newMethodsFromObject(pkg, obj),
	}, true
}

func (dt *DefinedType) Name() DefinedTypeName {
	return dt.name
}

func (dt *DefinedType) PackageSummary() *PackageSummary {
	return dt.pkgSummary
}

func (dt *DefinedType) Type() *Type {
	return dt.typ
}

func (dt *DefinedType) Methods() []*Function {
	return dt.methods.asSlice()
}

func newDefinedList(pkg *packages.Package) *DefinedTypeList {
	var definedTypes []*DefinedType
	scope := pkg.Types.Scope()
	for _, name := range scope.Names() {
		obj := scope.Lookup(name)
		if a, ok := newDefinedTypeIfObjectDefinedType(pkg, obj); ok {
			definedTypes = append(definedTypes, a)
		}
	}
	return &DefinedTypeList{definedTypes: definedTypes}
}

func (dtl *DefinedTypeList) asSlice() []*DefinedType {
	return append([]*DefinedType{}, dtl.definedTypes...)
}
