package gocode

import (
	"go/types"
	"strings"

	"golang.org/x/tools/go/packages"
)

type (
	// StructName は、Goのstruct名を表す。
	StructName string

	// PackageStructName は、パッケージ名付きのstruct名を表す。
	PackageStructName string

	// Struct は、Goのstructを表す。
	Struct struct {
		goType     types.Type
		structName StructName
		pkgSummary *PackageSummary
		typ        *Type
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

func newStructList(pkg *packages.Package) *StructList {
	var structs []*Struct
	scope := pkg.Types.Scope()
	for _, name := range scope.Names() {
		obj := scope.Lookup(name)
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

func newStructIfStructType(pkg *packages.Package, obj types.Object) (res *Struct, ok bool) {
	structType, ok := obj.Type().Underlying().(*types.Struct)
	if !ok {
		return &Struct{}, false
	}

	pkgSummary := newPackageSummaryFromGoTypes(obj.Pkg())

	s := &Struct{
		goType:     obj.Type(),
		pkgSummary: pkgSummary,
		typ:        newType(pkgSummary, obj.Type()),
		structName: StructName(obj.Name()),
		fields:     newFieldListFromStructType(structType),
		methods:    newMethodsFromObject(pkg, obj),
		implements: newPackageInterfaceMap(),
	}

	return s, true
}

func (s *Struct) PackageSummary() *PackageSummary {
	return s.pkgSummary
}

func (s *Struct) Name() StructName {
	return s.structName
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
	if len(i.Methods()) == 0 {
		return false
	}

	typ := s.goType
	switch t := typ.(type) {
	case *types.Pointer:
	default:
		// pointerにしておかないとtypes.Implementsで正しく判定されない
		typ = types.NewPointer(t)
	}
	return types.Implements(typ, i.goInterface)
}

func (s *Struct) addInterfaceIfImplements(i *Interface) {
	if s.Implements(i) {
		s.implements.put(i)
	}
}
