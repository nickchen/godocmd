package godocmd

import (
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/nickchen/godocmd/templates"
)

type DocMD struct {
	workingDir string
	template   *template.Template
}

func New() (*DocMD, error) {
	d := &DocMD{}
	var err error
	d.workingDir, err = filepath.Abs("./")
	return d, err
}

type pkgWithTest struct {
	main *doc.Package
	test *doc.Package
}

func (p *pkgWithTest) writeOutPackageMD(fset *token.FileSet, workingDir, outDir string) error {
	pkg := NewPackage(p.main, p.test)
	// this should be a path base on current working directory
	outfile := filepath.Join(outDir, fmt.Sprintf("%s.md", pkg.Name))
	directory := filepath.Dir(outfile)
	_ = os.MkdirAll(directory, os.ModePerm)
	f, err := os.Create(outfile)
	if err != nil {
		return err
	}
	defer func() {
		f.Sync()
		f.Close()
	}()

	temp := template.New("markdown.tmpl")
	if err != nil {
		return err
	}
	temp, err = templates.Parse(temp.Funcs(*templateFuncMap(fset, pkg, workingDir, outfile)))
	if err != nil {
		return err
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
	pkgWithTests := make(map[string]*pkgWithTest)
	for pkgName, pkgAst := range pkgs {
		if pkgName == "main" {
			continue
		}
		if strings.HasSuffix(pkgName, "_test") {
			mainPkgName := strings.TrimSuffix(pkgName, "_test")
			if p, ok := pkgWithTests[mainPkgName]; ok {
				p.test = doc.New(pkgAst, packageBasePath, doc.PreserveAST)
			} else {
				pkgWithTests[pkgName] = &pkgWithTest{test: doc.New(pkgAst, packageBasePath, doc.PreserveAST)}
			}
		} else {
			if p, ok := pkgWithTests[pkgName]; ok {
				p.main = doc.New(pkgAst, packageBasePath, doc.PreserveAST)
			} else {
				pkgWithTests[pkgName] = &pkgWithTest{main: doc.New(pkgAst, packageBasePath, doc.PreserveAST)}
			}
		}
	}
	for _, pkg := range pkgWithTests {
		pkg.writeOutPackageMD(fset, d.workingDir, outDir)
	}
	return nil
}

func (d *DocMD) ProcessPackageDirs(outDir, packageBasePath string, dirs ...string) error {
	// for each package dir, walk the directory tree, and create fileset for each, process
	// fileset using ast
	outFullPath, err := filepath.Abs(outDir)
	if err != nil {
		return err
	}
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
				packageImportPath := fmt.Sprintf("%s/%s", packageBasePath, relPath)
				if packageBasePath == "" {
					packageImportPath = relPath
				}
				if strings.HasSuffix(sourcePath, ".git") {
					return filepath.SkipDir
				}
				if strings.HasPrefix(sourcePath, outFullPath) {
					return filepath.SkipDir
				}
				fmt.Println("  * Doing", sourcePath, relPath)
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

// getParamNames from ast.Ident, ie: a, b string or a string
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
