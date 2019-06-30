package godocmd

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type DocMD struct {
	templateFile string
}

type Package struct {
	*doc.Package
	// FuncsFiltered filtered out examples and benchmarks
	FuncsFiltered []*doc.Func
	FuncsName     map[string]*doc.Func
}

func NewPackage(docPkg *doc.Package) *Package {
	pkg := &Package{docPkg,
		make([]*doc.Func, 0),
		make(map[string]*doc.Func)}
	pkg.init()
	return pkg
}

func (pkg *Package) init() {
	for _, f := range pkg.Funcs {
		matched := false
		for _, prefix := range []string{"Example", "Benchmark", "Test"} {
			if strings.HasPrefix(f.Name, prefix) {
				if prefix != "Test" {
					pkg.FuncsName[f.Name] = f
				}
				matched = true
				break
			}
		}
		if !matched {
			pkg.FuncsFiltered = append(pkg.FuncsFiltered, f)
		}
	}
}

func New(templateFile string) (*DocMD, error) {
	return &DocMD{templateFile: templateFile}, nil
}

func (d *DocMD) writeOutPackageMD(docPkg *doc.Package, fset *token.FileSet, name, outDir string) error {
	pkg := NewPackage(docPkg)
	outfile := filepath.Join(outDir, fmt.Sprintf("%s.md", pkg.Name))
	f, err := os.Create(outfile)
	if err != nil {
		return err
	}
	defer func() {
		f.Sync()
		f.Close()
	}()
	temp, err := template.New(filepath.Base(d.templateFile)).
		Funcs(*d.templateFuncMap(fset, pkg)).
		ParseFiles(d.templateFile)
	if err != nil {
		return fmt.Errorf("failed to parse template: %s", err)
	}

	err = temp.Execute(f, pkg)
	if err != nil {
		return err
	}
	return nil
}

func (d *DocMD) processDir(outDir, packageBasePath, sourcePath string) error {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, sourcePath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse source: %s", err)
	}
	for pkgName, pkgAst := range pkgs {
		if strings.HasSuffix(pkgName, "_test") {
			continue
		}
		pkg := doc.New(pkgAst, packageBasePath, doc.PreserveAST)
		err := d.writeOutPackageMD(pkg, fset, pkgName, outDir)
		if err != nil {
			return fmt.Errorf("failed to write out: %s", err)
		}
	}
	return nil
}

func (d *DocMD) ProcessPackageDirs(outDir, packageBasePath string, dirs ...string) error {
	for _, dir := range dirs {
		abspath, err := filepath.Abs(dir)
		if err != nil {
			fmt.Printf("failed abs on path %s: %s", dir, err)
			continue
		}
		prefixDir, rootDir := filepath.Split(abspath)
		fmt.Println("Processing", prefixDir, rootDir)
		// dirBase := filepath.Base(dir)
		err = filepath.Walk(abspath, func(sourcePath string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("walk failure accessing a path %q: %v\n", sourcePath, err)
				return err
			}
			if info.IsDir() {
				relPath, err := filepath.Rel(prefixDir, sourcePath)
				currentOutDir := filepath.Dir(filepath.Join(outDir, relPath))
				if _, err := os.Stat(currentOutDir); os.IsNotExist(err) {
					// output path doesn't exist, need to make
					os.Mkdir(currentOutDir, os.ModePerm)
				}
				fmt.Println("  * Doing", sourcePath, relPath)
				packageImportPath := fmt.Sprintf("%s/%s", packageBasePath, relPath)
				if packageBasePath == "" {
					packageImportPath = relPath
				}
				err = d.processDir(currentOutDir, packageImportPath, sourcePath)
				if err != nil {
					fmt.Printf("failed to process path %q: %v\n", sourcePath, err)
					return err
				}
			}
			return nil
		})
		if err != nil {
			fmt.Printf("error walking the path %q: %v\n", dir, err)
			return err
		}
	}
	return nil
}

func anyTypeSourceString(fset *token.FileSet, _type interface{}) string {
	var buffer bytes.Buffer
	printer.Fprint(&buffer, fset, _type)
	return buffer.String()
}

func (d *DocMD) anyTypeSourceString(fset *token.FileSet) func(interface{}) string {
	return func(_type interface{}) string {
		return anyTypeSourceString(fset, _type)
	}
}

func getParamNames(idents []*ast.Ident, paramTypeString string) string {
	if len(idents) == 0 {
		return paramTypeString
	}
	paramNames := []string{}
	for _, i := range idents {
		if i.Name == "" {
			continue
		}
		paramNames = append(paramNames, i.Name)
	}
	return fmt.Sprint(strings.Join(paramNames, ", "), " ", paramTypeString)
}

func (d *DocMD) functionParam(fset *token.FileSet) func(string, *ast.FuncDecl) string {
	return func(_type string, decl *ast.FuncDecl) string {
		var params []string
		if decl != nil {
			offset := decl.Pos()
			if decl.Doc != nil {
				offset = decl.Doc.Pos()
			}
			paramDeclString := anyTypeSourceString(fset, decl)
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

func (d *DocMD) sourceFileLink(fset *token.FileSet) func(string) string {
	return func(_filepath string) string {
		return fmt.Sprintf("[%s](%s)", filepath.Base(_filepath), _filepath)
	}
}

func (d *DocMD) functionSignature(fset *token.FileSet) func(*doc.Func) string {
	return func(f *doc.Func) string {
		functionGet := d.functionParam(fset)
		return fmt.Sprintf("func %s%s%s", f.Decl.Name.Name, functionGet("params", f.Decl), functionGet("results", f.Decl))
	}
}

func (d *DocMD) functionAnchor(fset *token.FileSet) func(string, interface{}) string {
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

func (d *DocMD) getExampleForFunc(pkg *Package) func(*doc.Type, *doc.Func) []*doc.Func {
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
			}
		}
		return ret
	}
}

func (d *DocMD) templateFuncMap(fset *token.FileSet, pkg *Package) *template.FuncMap {
	return &template.FuncMap{
		"functionSignature":   d.functionSignature(fset),
		"anyTypeSourceString": d.anyTypeSourceString(fset),
		"sourceFileLink":      d.sourceFileLink(fset),
		"anchorFunc":          d.functionAnchor(fset),
		"getExampleForFunc":   d.getExampleForFunc(pkg),
	}
}
