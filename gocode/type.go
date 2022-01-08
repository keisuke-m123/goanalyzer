package gocode

import (
	"fmt"
	"go/types"
	"strings"
)

type (
	// TypeName は、 types.Type を文字列に変換した型名を表す。
	TypeName string

	// RelativeFullTypeName は Type を保持する構造体が所属するパッケージから見た相対的なパッケージ名付きの型名を表す。
	//
	// Type を保持する構造体と表現される型が所属するパッケージが同一であればパッケージ名は付与されず型名のみとなる。
	//
	// Type を保持する構造体と表現される型が所属するパッケージが同一でなければパッケージ名は付与されず型名のみとなる。
	RelativeFullTypeName string

	// Type は型を表す。
	Type struct {
		// goType は解析元の types.Type 。
		goType types.Type
		// typeName は、 goType を文字列に変換した型名を表す。
		typeName TypeName
		// relativeFullTypeName は、 goType を保持する構造体が所属するパッケージから見た相対的なパッケージ名付きの型名を表す。
		relativeFullTypeName RelativeFullTypeName
		// pkgSummary は types.Type の所属するパッケージのサマリ。
		pkgSummary *PackageSummary
		// fundamentalTypes は types.Type の基底となる Type 情報の一覧。
		fundamentalTypes []*Type
	}

	// typeConverter は、 types.Type から Type を生成するためのコンバータ。
	typeConverter struct {
		currentPkgSummary *PackageSummary
	}
)

func (tn TypeName) String() string {
	return string(tn)
}

func (tn TypeName) builtin() bool {
	return builtin(tn.String())
}

func (ftn RelativeFullTypeName) String() string {
	return string(ftn)
}

func newTypeWithoutFundamentalTypes(currentPkgSummary *PackageSummary, typ types.Type) *Type {
	t := &Type{
		goType:     typ,
		pkgSummary: currentPkgSummary,
	}
	switch goType := typ.(type) {
	case *types.Named:
		if goType.Obj().Pkg() != nil {
			t.pkgSummary = newPackageSummaryFromGoTypes(goType.Obj().Pkg())
		}
	}

	c := newTypeConverter(currentPkgSummary)
	t.typeName = c.typeName(typ)
	if t.typeName.builtin() {
		t.pkgSummary = &PackageSummary{}
	}

	t.resetFullTypeName(currentPkgSummary, t.pkgSummary)

	return t
}

func newType(currentPkgSummary *PackageSummary, typ types.Type) *Type {
	t := newTypeWithoutFundamentalTypes(currentPkgSummary, typ)

	c := newTypeConverter(currentPkgSummary)
	t.fundamentalTypes = append(t.fundamentalTypes, c.fundamentalTypes(typ)...)

	return t
}

func (t *Type) PackageSummary() *PackageSummary {
	return t.pkgSummary
}

func (t *Type) TypeName() TypeName {
	return t.typeName
}

func (t *Type) RelativeFullTypeName() RelativeFullTypeName {
	return t.relativeFullTypeName
}

func (t *Type) FundamentalTypes() []*Type {
	return append([]*Type{}, t.fundamentalTypes...)
}

func (t *Type) ContainsBuiltinInFundamentalTypes() bool {
	for _, ft := range t.fundamentalTypes {
		if ft.typeName.builtin() {
			return true
		}
	}
	return false
}

func (t *Type) Builtin() bool {
	return t.typeName.builtin()
}

func (t *Type) resetFullTypeName(
	currentPkgSummary *PackageSummary,
	typePackageSummary *PackageSummary,
) {
	switch {
	case t.typeName.builtin():
		t.relativeFullTypeName = RelativeFullTypeName(t.typeName.String())
	case currentPkgSummary.Equal(typePackageSummary):
		t.relativeFullTypeName = RelativeFullTypeName(t.typeName.String())
	default:
		t.relativeFullTypeName = RelativeFullTypeName(fmt.Sprintf("%s.%s", typePackageSummary.Name(), t.typeName.String()))
	}
}

func newTypeConverter(currentPkgSummary *PackageSummary) *typeConverter {
	return &typeConverter{currentPkgSummary: currentPkgSummary}
}

// _typeName の戻り値を TypeName として返す。
func (tc *typeConverter) typeName(typ types.Type) TypeName {
	return TypeName(tc._typeName(typ))
}

// typ を表示可能な形式の文字列に変換可能して返す。
func (tc *typeConverter) _typeName(typ types.Type) string {
	switch t := typ.(type) {
	case *types.Basic:
		return tc.typeNameBasic(t)
	case *types.Slice:
		return tc.typeNameSlice(t)
	case *types.Array:
		return tc.typeNameArray(t)
	case *types.Map:
		return tc.typeNameMap(t)
	case *types.Pointer:
		return tc.typeNamePointer(t)
	case *types.Chan:
		return tc.typeNameChan(t)
	case *types.Struct:
		return tc.typeNameStruct(t)
	case *types.Interface:
		return tc.typeNameInterface(t)
	case *types.Signature:
		return tc.typeNameSignature(t)
	case *types.Named:
		return tc.typeNameNamed(t)
	case *types.Tuple:
		return ""
	default:
		return ""
	}
}

func (tc *typeConverter) typeNameBasic(t *types.Basic) string {
	return t.Name()
}

func (tc *typeConverter) typeNameSlice(t *types.Slice) string {
	eType := tc._typeName(t.Elem())
	return fmt.Sprintf("[]%s", eType)
}

func (tc *typeConverter) typeNameArray(t *types.Array) string {
	eType := tc._typeName(t.Elem())
	return fmt.Sprintf("[]%s", eType)
}

func (tc *typeConverter) typeNameMap(t *types.Map) string {
	kType := tc._typeName(t.Key())
	eType := tc._typeName(t.Elem())
	return fmt.Sprintf("map[%s]%s", kType, eType)
}

func (tc *typeConverter) typeNamePointer(t *types.Pointer) string {
	eType := tc._typeName(t.Elem())
	return fmt.Sprintf("*%s", eType)
}

func (tc *typeConverter) typeNameChan(t *types.Chan) string {
	eType := tc._typeName(t.Elem())
	return fmt.Sprintf("chan %s", eType)
}

func (tc *typeConverter) typeNameStruct(t *types.Struct) string {
	fieldList := make([]string, 0)
	for i := 0; i < t.NumFields(); i++ {
		fType := tc._typeName(t.Field(i).Type())
		fieldList = append(fieldList, fType)
	}
	return fmt.Sprintf("struct{%s}", strings.Join(fieldList, ", "))
}

func (tc *typeConverter) typeNameInterface(t *types.Interface) string {
	methods := make([]string, 0)
	for i := 0; i < t.NumMethods(); i++ {
		m := t.Method(i)
		methods = append(methods, fmt.Sprintf("%s %s", m.Name(), tc._typeName(m.Type())))
	}
	return fmt.Sprintf("interface{%s}", strings.Join(methods, "; "))
}

func (tc *typeConverter) typeNameSignature(t *types.Signature) string {
	paramTypes := make([]string, 0)
	for _, p := range newParameters(t) {
		paramTypes = append(paramTypes, p.Type().RelativeFullTypeName().String())
	}

	returnValues := make([]string, 0)
	for _, r := range newReturnValues(t) {
		returnValues = append(returnValues, r.Type().RelativeFullTypeName().String())
	}

	var returns string
	if len(returnValues) > 1 {
		returns = fmt.Sprintf("(%s)", strings.Join(returnValues, ", "))
	} else {
		returns = strings.Join(returnValues, "")
	}
	return fmt.Sprintf("func(%s) %s", strings.Join(paramTypes, ", "), returns)
}

func (tc *typeConverter) typeNameNamed(t *types.Named) string {
	return t.Obj().Name()
}

// typ の基底となる型情報を解析して一覧で返す。
func (tc *typeConverter) fundamentalTypes(typ types.Type) []*Type {
	switch t := typ.(type) {
	case *types.Basic:
		return []*Type{newTypeWithoutFundamentalTypes(tc.currentPkgSummary, t)}
	case *types.Slice:
		return tc.fundamentalTypes(t.Elem())
	case *types.Array:
		return tc.fundamentalTypes(t.Elem())
	case *types.Map:
		keyRes := tc.fundamentalTypes(t.Key())
		elmRes := tc.fundamentalTypes(t.Elem())
		return append(keyRes, elmRes...)
	case *types.Pointer:
		return tc.fundamentalTypes(t.Elem())
	case *types.Chan:
		return tc.fundamentalTypes(t.Elem())
	case *types.Struct:
		return []*Type{}
	case *types.Interface:
		return []*Type{}
	case *types.Signature:
		return tc.signatureFundamentalTypes(t)
	case *types.Named:
		return []*Type{newTypeWithoutFundamentalTypes(tc.currentPkgSummary, t)}
	case *types.Tuple:
		return []*Type{}
	default:
		return []*Type{}
	}
}

func (tc *typeConverter) signatureFundamentalTypes(t *types.Signature) []*Type {
	var ts []*Type
	for i := 0; i < t.Params().Len(); i++ {
		ts = append(ts, tc.fundamentalTypes(t.Params().At(i).Type())...)
	}
	for i := 0; i < t.Results().Len(); i++ {
		ts = append(ts, tc.fundamentalTypes(t.Results().At(i).Type())...)
	}
	return ts
}
