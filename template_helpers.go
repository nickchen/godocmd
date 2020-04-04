package godocmd

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/doc"
	"go/printer"
	"go/token"
	"path/filepath"
	"strings"
	"text/template"
)

type Package struct {
	*doc.Package
	// FuncsFiltered filtered out examples and benchmarks
	FuncsFiltered []*doc.Func
	FuncsName     map[string]*doc.Func
}

func NewPackage(mainPkg, testPkg *doc.Package) *Package {
	pkg := &Package{mainPkg,
		make([]*doc.Func, 0),
		make(map[string]*doc.Func)}
	pkg.init(testPkg)
	return pkg
}

func filter(ary []*doc.Func, f func(fun *doc.Func) bool) []*doc.Func {
	fs := make([]*doc.Func, 0)
	for _, a := range ary {
		if f(a) {
			fs = append(fs, a)
		}
	}
	return fs
}

func stringStartsWith(in string, patterns []string) bool {
	for _, p := range patterns {
		if strings.HasPrefix(in, p) {
			return true
		}
	}
	return false
}

func (pkg *Package) init(testPkg *doc.Package) {
	pkg.FuncsFiltered = filter(pkg.Funcs, func(fun *doc.Func) bool {
		pkg.FuncsName[fun.Name] = fun

		return !stringStartsWith(fun.Name, []string{"Benchmark", "Test", "Example"})
	})
	if testPkg != nil {
		for _, f := range testPkg.Funcs {
			if stringStartsWith(f.Name, []string{"Example"}) {
				pkg.FuncsName[f.Name] = f
			}
		}
	}
}

func templateFuncMap(fset *token.FileSet, pkg *Package, workingDir, outfile string) *template.FuncMap {
	return &template.FuncMap{
		"functionSignature":       functionSignature(fset),
		"anyTypeSourceString":     anyTypeSourceString(fset),
		"sourceFileLink":          sourceFileLink(fset, workingDir, outfile),
		"anchorFunc":              functionAnchor(fset),
		"getExampleForFunc":       getExampleForFunc(pkg),
		"getOtherExamplesForType": getOtherExamplesForType(pkg),
		"codeBlock":               func() string { return "```" },
		"codeBlockGolang":         func() string { return "```go" },
	}
}

func anyTypeSourceStringFunc(fset *token.FileSet, _type interface{}) string {
	var buffer bytes.Buffer
	printer.Fprint(&buffer, fset, _type)
	return buffer.String()
}

func anyTypeSourceString(fset *token.FileSet) func(interface{}) string {
	return func(_type interface{}) string {
		return anyTypeSourceStringFunc(fset, _type)
	}
}

func functionSignature(fset *token.FileSet) func(*doc.Func) string {
	return func(f *doc.Func) string {
		functionGet := functionParam(fset)
		return fmt.Sprintf("func %s%s%s", f.Decl.Name.Name, functionGet("params", f.Decl), functionGet("results", f.Decl))
	}
}

// functionParam from ast.FuncDecl, this is for both parameters, and results
func functionParam(fset *token.FileSet) func(string, *ast.FuncDecl) string {
	return func(_type string, decl *ast.FuncDecl) string {
		var params []string
		if decl != nil {
			offset := decl.Pos()
			if decl.Doc != nil {
				offset = decl.Doc.Pos()
			}
			paramDeclString := anyTypeSourceStringFunc(fset, decl)
			if fields, ok := map[string]*ast.FieldList{"params": decl.Type.Params, "results": decl.Type.Results}[_type]; ok && fields != nil {
				for _, field := range fields.List {
					fieldType := paramDeclString[field.Type.Pos()-offset : field.Type.End()-offset]
					params = append(params, getParamNames(field.Names, fieldType))
				}
			}
		}
		if len(params) > 0 {
			return fmt.Sprintf("(%s)", strings.Join(params, ", "))
		} else if _type == "params" {
			return "()"
		}
		return ""
	}
}

// sourceFileLink generate links for the source file, need to rework base on relative path of the doc
func sourceFileLink(fset *token.FileSet, workingDir, outfile string) func(string) string {
	outLevel := len(filepath.SplitList(outfile))
	sourceRel := strings.Repeat("../", outLevel)
	return func(_filepath string) string {

		relPath, _ := filepath.Rel(workingDir, _filepath)
		return fmt.Sprintf("[%s](%s%s)", filepath.Base(_filepath), sourceRel, relPath)
	}
}

func functionAnchor(fset *token.FileSet) func(string, interface{}) string {
	return func(linkType string, _f interface{}) string {
		var name, objType string
		switch f := _f.(type) {
		case *doc.Type:
			name = f.Name
			objType = "type"
		case *doc.Func:
			name = f.Decl.Name.Name
			objType = "func"
		default:
			panic(fmt.Errorf("unhandle type %v", _f))
		}

		if t, ok := map[string]string{
			"link":   strings.ToLower(fmt.Sprintf("#%s-%s", objType, name)),
			"anchor": fmt.Sprintf("%s-%s", objType, name),
		}[linkType]; ok {
			return t
		} else {
			return ""
		}
	}
}

func getExampleForFunc(pkg *Package) func(*doc.Type, *doc.Func) []*doc.Func {
	return func(t *doc.Type, f *doc.Func) []*doc.Func {
		ret := make([]*doc.Func, 0)
		for _, prefix := range []string{"Example", "Benchmark"} {
			var funcKey string
			switch {
			case t != nil && f != nil:
				funcKey = fmt.Sprintf("%s%s_%s", prefix, t.Name, f.Name)
			case t == nil && f != nil:
				funcKey = fmt.Sprintf("%s%s", prefix, f.Name)
			case t != nil:
				funcKey = fmt.Sprintf("%s%s", prefix, t.Name)
			default:
				return nil
			}
			if e, ok := pkg.FuncsName[funcKey]; ok {
				ret = append(ret, e)
				//delete (pkg.FuncsName, funcKey)
			}
		}
		return ret
	}
}

func getOtherExamplesForType(pkg *Package) func(*doc.Type) []*doc.Func {
	return func(t *doc.Type) []*doc.Func {
		ret := make([]*doc.Func, 0)
		for key, e := range pkg.FuncsName {
			if strings.HasPrefix(key, "Example"+t.Name) || strings.HasPrefix(key, "Benchmark"+t.Name) {
				fmt.Printf("FUN %s\n", t.Name, e)
				ret = append(ret, e)
			}
		}
		return ret
	}
}
