package gocode

import (
	"go/types"
	"strings"

	"golang.org/x/tools/go/packages"
)

type (
	// TypeAliasName は、型別名を表す。
	TypeAliasName string

	// PackageTypeAliasName は、パッケージ内の型別名を表す。
	PackageTypeAliasName string

	// TypeAlias は、型別名を表す。
	TypeAlias struct {
		name       TypeAliasName
		pkgSummary *PackageSummary
		typ        *Type
	}

	// TypeAliasList は、型別名のリストを表す。
	TypeAliasList struct {
		aliases []*TypeAlias
	}
)

func (an TypeAliasName) String() string {
	return string(an)
}

func NewPackageAliasName(pkgName PackageName, aliasName TypeAliasName) PackageTypeAliasName {
	return PackageTypeAliasName(strings.Join([]string{
		pkgName.String(),
		aliasName.String(),
	}, "."))
}

func (pan PackageTypeAliasName) String() string {
	return string(pan)
}

func newTypeAliasIfObjectTypeAlias(obj types.Object) (res *TypeAlias, ok bool) {
	tn, ok := obj.(*types.TypeName)
	if !ok || !tn.IsAlias() {
		return &TypeAlias{}, false
	}

	pkgSummary := newPackageSummaryFromGoTypes(obj.Pkg())

	return &TypeAlias{
		name:       TypeAliasName(obj.Name()),
		pkgSummary: pkgSummary,
		typ:        newType(pkgSummary, obj.Type().Underlying()),
	}, true
}

func (a *TypeAlias) PackageSummary() *PackageSummary {
	return a.pkgSummary
}

func (a *TypeAlias) PackageAliasName() PackageTypeAliasName {
	return NewPackageAliasName(a.pkgSummary.name, a.name)
}

func (a *TypeAlias) Name() TypeAliasName {
	return a.name
}

func (a *TypeAlias) Type() *Type {
	return a.typ
}

func newAliasList(pkg *packages.Package) *TypeAliasList {
	var aliases []*TypeAlias
	scope := pkg.Types.Scope()
	for _, name := range scope.Names() {
		obj := scope.Lookup(name)
		if a, ok := newTypeAliasIfObjectTypeAlias(obj); ok {
			aliases = append(aliases, a)
		}
	}
	return &TypeAliasList{aliases: aliases}
}

func (al *TypeAliasList) asSlice() []*TypeAlias {
	var slice []*TypeAlias
	for i := range al.aliases {
		slice = append(slice, al.aliases[i])
	}
	return slice
}
