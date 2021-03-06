package gocode

import (
	"go/token"
	"go/types"
	"strings"
)

type (
	// StructName は、Goのstruct名を表す。
	StructName string

	// PackageStructName は、パッケージ名付きのstruct名を表す。
	PackageStructName string

	// Struct は、Goのstructを表す。
	Struct struct {
		definedPos token.Pos
		typ        *Type
		structName StructName
		pkgSummary *PackageSummary
		methods    *FunctionList
		fields     *FieldList
		implements *PackageInterfaceMap
	}

	// StructList は、Goのstructのリストを表す。
	StructList struct {
		structs []*Struct
	}
)

func (sn StructName) String() string {
	return string(sn)
}

func (sn StructName) EqualString(s string) bool {
	return sn.String() == s
}

func (psn PackageStructName) String() string {
	return string(psn)
}

func (psn PackageStructName) EqualString(s string) bool {
	return psn.String() == s
}

func NewPackageStructName(pkgName PackageName, structName StructName) PackageStructName {
	return PackageStructName(strings.Join([]string{pkgName.String(), structName.String()}, "."))
}

func newStructList(pkg packageIn) *StructList {
	var structs []*Struct
	for _, obj := range pkg.Typed() {
		if s, ok := newStructIfStructType(pkg, obj); ok {
			structs = append(structs, s)
		}
	}
	return &StructList{structs: structs}
}

func (s *StructList) asSlice() []*Struct {
	var slice []*Struct
	for i := range s.structs {
		slice = append(slice, s.structs[i])
	}
	return slice
}

func newStructIfStructType(pkg packageIn, obj types.Object) (res *Struct, ok bool) {
	structType, ok := obj.Type().Underlying().(*types.Struct)
	if !ok {
		return &Struct{}, false
	}

	pkgSummary := newPackageSummaryFromGoTypes(obj.Pkg())

	s := &Struct{
		definedPos: obj.Pos(),
		pkgSummary: pkgSummary,
		typ:        newType(pkgSummary, obj.Type()),
		structName: StructName(obj.Name()),
		fields:     newFieldListFromStructType(structType),
		methods:    newMethodsFromObject(pkg, obj),
		implements: newPackageInterfaceMap(),
	}

	return s, true
}

func (s *Struct) DefinedPos() token.Pos {
	return s.definedPos
}

func (s *Struct) PackageSummary() *PackageSummary {
	return s.pkgSummary
}

func (s *Struct) Name() StructName {
	return s.structName
}

func (s *Struct) Type() *Type {
	return s.typ
}

func (s *Struct) PackageStructName() PackageStructName {
	return NewPackageStructName(s.pkgSummary.Name(), s.Name())
}

func (s *Struct) Methods() []*Function {
	return s.methods.asSlice()
}

func (s *Struct) Fields() []*Field {
	return s.fields.asSlice()
}

func (s *Struct) ImplementInterfaces() *PackageInterfaceMap {
	return s.implements
}

func (s *Struct) Implements(i *Interface) bool {
	return implements(s.Type().GoType(), i.goInterface)
}

func (s *Struct) ImplementsGoTypes(i *types.Interface) bool {
	return implements(s.Type().GoType(), i)
}

func (s *Struct) addInterfaceIfImplements(i *Interface) {
	if s.Implements(i) {
		s.implements.put(i)
	}
}
