package gocode

import (
	"go/token"
	"go/types"
)

type (
	// DefinedTypeName はdefined typeの名前を表す。
	DefinedTypeName string

	// DefinedType はdefined typeを表す。
	DefinedType struct {
		// definedPos はコード中で DefinedType が定義された位置。
		definedPos token.Pos
		// typ は DefinedType 自体の型情報。
		typ *Type
		// underlyingTyp はtypeされた型情報。
		underlyingTyp *Type
		// name はtypeされた名前。
		name DefinedTypeName
		// pkgSummary は定義されているパッケージのサマリ情報。
		pkgSummary *PackageSummary
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

func newDefinedTypeIfObjectDefinedType(pkg packageIn, obj types.Object) (res *DefinedType, ok bool) {
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
		definedPos:    obj.Pos(),
		typ:           newType(pkgSummary, obj.Type()),
		underlyingTyp: newType(pkgSummary, obj.Type().Underlying()),
		pkgSummary:    pkgSummary,
		name:          DefinedTypeName(obj.Name()),
		methods:       newMethodsFromObject(pkg, obj),
	}, true
}

func (dt *DefinedType) DefinedPos() token.Pos {
	return dt.definedPos
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

func (dt *DefinedType) UnderlyingType() *Type {
	return dt.underlyingTyp
}

func (dt *DefinedType) Methods() []*Function {
	return dt.methods.asSlice()
}

func (dt *DefinedType) Implements(i *Interface) bool {
	return implements(dt.Type().GoType(), i.goInterface)
}

func (dt *DefinedType) ImplementsGoTypes(i *types.Interface) bool {
	return implements(dt.Type().GoType(), i)
}

func newDefinedList(pkg packageIn) *DefinedTypeList {
	var definedTypes []*DefinedType
	for _, obj := range pkg.Typed() {
		if a, ok := newDefinedTypeIfObjectDefinedType(pkg, obj); ok {
			definedTypes = append(definedTypes, a)
		}
	}
	return &DefinedTypeList{definedTypes: definedTypes}
}

func (dtl *DefinedTypeList) asSlice() []*DefinedType {
	return append([]*DefinedType{}, dtl.definedTypes...)
}
