package gocode

import (
	"go/types"

	"golang.org/x/tools/go/packages"
)

type (
	// FunctionName は、関数名を表す。
	FunctionName string

	// ParameterName は、関数のパラメータ名を表す。
	ParameterName string

	// Parameter は、関数のパラメータを表す。
	Parameter struct {
		name string
		typ  *Type
	}

	// Parameters は、関数のパラメータリスト。
	Parameters []*Parameter

	// ReturnValueName は、関数の戻り値名を表す。
	ReturnValueName string

	// ReturnValue は、関数の戻り値を表す。
	ReturnValue struct {
		name string
		typ  *Type
	}

	// ReturnValues は、関数の戻り値リスト。
	ReturnValues []*ReturnValue

	// Function は、関数を表す。
	Function struct {
		name         FunctionName
		parameters   Parameters
		returnValues ReturnValues
	}

	// FunctionList は、関数リストを表す。
	FunctionList struct {
		functions []*Function
	}
)

func (f FunctionName) String() string {
	return string(f)
}

func newParameter(obj types.Object) *Parameter {
	return &Parameter{
		name: obj.Name(),
		typ:  newType(newPackageSummaryFromGoTypes(obj.Pkg()), obj.Type()),
	}
}

func (p Parameter) Name() string {
	return p.name
}

func (p Parameter) Type() *Type {
	return p.typ
}

func newReturnValue(obj types.Object) *ReturnValue {
	return &ReturnValue{
		name: obj.Name(),
		typ:  newType(newPackageSummaryFromGoTypes(obj.Pkg()), obj.Type()),
	}
}

func (rv ReturnValue) Name() string {
	return rv.name
}

func (rv ReturnValue) Type() *Type {
	return rv.typ
}

func newFunctionIfSignatureType(f *types.Func) (*Function, bool) {
	fn := &Function{
		name: FunctionName(f.Name()),
	}

	s, ok := f.Type().(*types.Signature)
	if !ok {
		return fn, false
	}

	fn.parameters = newParameters(s)
	fn.returnValues = newReturnValues(s)

	return fn, true
}

func newParameters(signature *types.Signature) Parameters {
	res := make(Parameters, 0)
	params := signature.Params()
	if params != nil {
		for i := 0; i < params.Len(); i++ {
			p := signature.Params().At(i)
			res = append(res, newParameter(p))
		}
	}
	return res
}

func newReturnValues(signature *types.Signature) ReturnValues {
	res := make(ReturnValues, 0)
	results := signature.Results()
	if results != nil {
		for i := 0; i < results.Len(); i++ {
			r := signature.Results().At(i)
			res = append(res, newReturnValue(r))
		}
	}
	return res
}

func (f *Function) Name() FunctionName {
	return f.name
}

func (f *Function) Parameters() Parameters {
	return append(Parameters{}, f.parameters...)
}

func (f *Function) ReturnValues() ReturnValues {
	return append(ReturnValues{}, f.returnValues...)
}

func newFunctionListFromInterface(interfaceType *types.Interface) *FunctionList {
	var functions []*Function
	for i := 0; i < interfaceType.NumMethods(); i++ {
		m := interfaceType.Method(i)
		if fn, ok := newFunctionIfSignatureType(m); ok {
			functions = append(functions, fn)
		}
	}
	return &FunctionList{functions: functions}
}

func (fl *FunctionList) asSlice() []*Function {
	var slice []*Function
	for i := range fl.functions {
		slice = append(slice, fl.functions[i])
	}
	return slice
}

func newMethodsFromObject(pkg *packages.Package, obj types.Object) *FunctionList {
	var methods []*Function
	namedObj := pkg.Types.Scope().Lookup(obj.Name())
	if named, ok := namedObj.Type().(*types.Named); ok && named != nil {
		for i := 0; i < named.NumMethods(); i++ {
			funcObj := named.Method(i)
			if fn, ok := newFunctionIfSignatureType(funcObj); ok {
				methods = append(methods, fn)
			}
		}
	}
	return &FunctionList{functions: methods}
}
