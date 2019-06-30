// docmd will generate Mark Down formated package documentation for provided paths
package main

import (
	"flag"
	"log"
	"os"

	"github.com/nickchen/godocmd"
)

func main() {
	var outDir, packageBasePath, templateFile string
	defaultOutDir := "./docs"
	defaultPackagePath := "github.com/nickchen/godocmd"
	defaultTemplateFile := "./markdown.tmpl"
	flag.StringVar(&outDir, "output-dir", defaultOutDir, "output directory")
	flag.StringVar(&packageBasePath, "package-base", defaultPackagePath, `package import basepath, output from "go list -m"`)
	flag.StringVar(&templateFile, "template", defaultTemplateFile, `template file`)
	flag.Parse()
	d, err := godocmd.New(templateFile)
	if err != nil {
		log.Fatal(err)
	}
	args := flag.Args()
	if len(args) == 0 {
		args = []string{"./fixture"}
	}
	err = d.ProcessPackageDirs(outDir, packageBasePath, args...)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
