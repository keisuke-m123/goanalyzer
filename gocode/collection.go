package gocode

type PackageMap struct {
	m map[PackagePath]*Package
}

func newPackageMap() *PackageMap {
	return &PackageMap{
		m: make(map[PackagePath]*Package),
	}
}

func (p *PackageMap) AsSlice() []*Package {
	var packages []*Package
	for _, pkg := range p.m {
		packages = append(packages, pkg)
	}
	return packages
}

func (p *PackageMap) add(pkg *Package) {
	p.m[pkg.Summary().Path()] = pkg
}

func (p PackageMap) Contains(pkgPath PackagePath) bool {
	_, ok := p.m[pkgPath]
	return ok
}

type PackageStructureMap struct {
	m map[PackageName]map[StructName]*Struct
}

func newPackageStructureMap() *PackageStructureMap {
	return &PackageStructureMap{m: make(map[PackageName]map[StructName]*Struct)}
}

func (p *PackageStructureMap) Get(pkgName PackageName, structName StructName) (s *Struct, ok bool) {
	structMap, ok := p.m[pkgName]
	if !ok {
		return nil, false
	}
	s, ok = structMap[structName]
	return s, ok
}

func (p *PackageStructureMap) PackageNames() []PackageName {
	var names []PackageName
	for pkgName := range p.m {
		names = append(names, pkgName)
	}
	return names
}

func (p *PackageStructureMap) PackageStructNames() []PackageStructName {
	var names []PackageStructName
	for pkgName, structMap := range p.m {
		for structName := range structMap {
			names = append(names, NewPackageStructName(pkgName, structName))
		}
	}
	return names
}

func (p *PackageStructureMap) PackageStructs(pkgName PackageName) []*Struct {
	var structs []*Struct
	structMap, ok := p.m[pkgName]
	if !ok {
		return structs
	}
	for key := range structMap {
		structs = append(structs, structMap[key])
	}
	return structs
}

func (p *PackageStructureMap) StructAll() []*Struct {
	var structs []*Struct
	for pkgName := range p.m {
		for stName := range p.m[pkgName] {
			structs = append(structs, p.m[pkgName][stName])
		}
	}
	return structs
}

func (p *PackageStructureMap) Contains(pkgName PackageName, structName StructName) bool {
	_, ok := p.Get(pkgName, structName)
	return ok
}

func (p *PackageStructureMap) put(s *Struct) {
	pkgName := s.PackageSummary().Name()
	if _, ok := p.m[pkgName]; !ok {
		p.m[pkgName] = make(map[StructName]*Struct)
	}
	p.m[pkgName][s.Name()] = s
}

type PackageInterfaceMap struct {
	m map[PackageName]map[InterfaceName]*Interface
}

func newPackageInterfaceMap() *PackageInterfaceMap {
	return &PackageInterfaceMap{m: make(map[PackageName]map[InterfaceName]*Interface)}
}

func (p *PackageInterfaceMap) Get(pkgName PackageName, interfaceName InterfaceName) (iface *Interface, ok bool) {
	interfaceMap, ok := p.m[pkgName]
	if !ok {
		return nil, false
	}
	iface, ok = interfaceMap[interfaceName]
	return iface, ok
}

func (p *PackageInterfaceMap) PackageNames() []PackageName {
	var names []PackageName
	for pkgName := range p.m {
		names = append(names, pkgName)
	}
	return names
}

func (p *PackageInterfaceMap) PackageInterfaceNames() []PackageInterfaceName {
	var names []PackageInterfaceName
	for pkgName, interfaceMap := range p.m {
		for interfaceName := range interfaceMap {
			names = append(names, NewPackageInterfaceName(pkgName, interfaceName))
		}
	}
	return names
}

func (p *PackageInterfaceMap) PackageInterfaces(pkgName PackageName) []*Interface {
	var interfaces []*Interface
	interfaceMap, ok := p.m[pkgName]
	if !ok {
		return interfaces
	}
	for key := range interfaceMap {
		interfaces = append(interfaces, interfaceMap[key])
	}
	return interfaces
}

func (p *PackageInterfaceMap) InterfaceAll() []*Interface {
	var interfaces []*Interface
	for pkgName := range p.m {
		for itName := range p.m[pkgName] {
			interfaces = append(interfaces, p.m[pkgName][itName])
		}
	}
	return interfaces
}

func (p *PackageInterfaceMap) Contains(pkgName PackageName, interfaceName InterfaceName) bool {
	_, ok := p.Get(pkgName, interfaceName)
	return ok
}

func (p *PackageInterfaceMap) put(iface *Interface) {
	pkgName := iface.PackageSummary().Name()
	_, ok := p.m[pkgName]
	if !ok {
		p.m[pkgName] = make(map[InterfaceName]*Interface)
	}
	p.m[pkgName][iface.Name()] = iface
}

type PackageTypeAliasMap struct {
	m map[PackageName]map[TypeAliasName]*TypeAlias
}

func newPackageTypeAliasMap() *PackageTypeAliasMap {
	return &PackageTypeAliasMap{m: make(map[PackageName]map[TypeAliasName]*TypeAlias)}
}

func (p *PackageTypeAliasMap) Get(pkgName PackageName, aliasName TypeAliasName) (al *TypeAlias, ok bool) {
	aliasMap, ok := p.m[pkgName]
	if !ok {
		return nil, false
	}
	al, ok = aliasMap[aliasName]
	return al, ok
}

func (p *PackageTypeAliasMap) PackageNames() []PackageName {
	var names []PackageName
	for pkgName := range p.m {
		names = append(names, pkgName)
	}
	return names
}

func (p *PackageTypeAliasMap) PackageAliasNames() []PackageTypeAliasName {
	var names []PackageTypeAliasName
	for pkgName, aliasMap := range p.m {
		for aliasName := range aliasMap {
			names = append(names, NewPackageAliasName(pkgName, aliasName))
		}
	}
	return names
}

func (p *PackageTypeAliasMap) PackageAliases(pkgName PackageName) []*TypeAlias {
	var aliases []*TypeAlias
	aliasMap, ok := p.m[pkgName]
	if !ok {
		return aliases
	}
	for key := range aliasMap {
		aliases = append(aliases, aliasMap[key])
	}
	return aliases
}

func (p *PackageTypeAliasMap) AliasAll() []*TypeAlias {
	var aliases []*TypeAlias
	for pkgName := range p.m {
		for aliasName := range p.m[pkgName] {
			aliases = append(aliases, p.m[pkgName][aliasName])
		}
	}
	return aliases
}

func (p *PackageTypeAliasMap) Contains(pkgName PackageName, aliasName TypeAliasName) bool {
	_, ok := p.Get(pkgName, aliasName)
	return ok
}

func (p *PackageTypeAliasMap) put(al *TypeAlias) {
	pkgName := al.PackageSummary().Name()
	_, ok := p.m[pkgName]
	if !ok {
		p.m[pkgName] = make(map[TypeAliasName]*TypeAlias)
	}
	p.m[pkgName][al.Name()] = al
}

type PackageDefinedTypeMap struct {
	m map[PackageName]map[DefinedTypeName]*DefinedType
}

func newPackageDefinedTypeMap() *PackageDefinedTypeMap {
	return &PackageDefinedTypeMap{m: make(map[PackageName]map[DefinedTypeName]*DefinedType)}
}

func (p *PackageDefinedTypeMap) Get(pkgName PackageName, definedType DefinedTypeName) (*DefinedType, bool) {
	definedTypeMap, ok := p.m[pkgName]
	if !ok {
		return nil, false
	}
	dt, ok := definedTypeMap[definedType]
	return dt, ok
}

func (p *PackageDefinedTypeMap) PackageDefinedTypes(pkgName PackageName) []*DefinedType {
	var definedTypes []*DefinedType
	definedTypeMap, ok := p.m[pkgName]
	if !ok {
		return definedTypes
	}
	for key := range definedTypeMap {
		definedTypes = append(definedTypes, definedTypeMap[key])
	}
	return definedTypes
}

func (p *PackageDefinedTypeMap) put(definedType *DefinedType) {
	pkgName := definedType.PackageSummary().Name()
	_, ok := p.m[pkgName]
	if !ok {
		p.m[pkgName] = make(map[DefinedTypeName]*DefinedType)
	}
	p.m[pkgName][definedType.Name()] = definedType
}
