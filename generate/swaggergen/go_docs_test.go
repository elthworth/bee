// Copyright 2013 bee authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package swaggergen

import (
	"go/ast"
	"go/build"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

//package model
//
//import (
//	"github.com/shopspring/decimal"
//)
//
//type Object struct{
//	Total  decimal.Decimal
//}
func TestCheckAndLoadPackageOnGoMod(t *testing.T) {
	var (
		pkgName       = "decimal"
		pkgImportPath = "github.com/shopspring/decimal"
	)

	defer os.Setenv("GO111MODULE", os.Getenv("GO111MODULE"))
	os.Setenv("GO111MODULE", "on")

	imports := []*ast.ImportSpec{
		{
			Path: &ast.BasicLit{
				Value: pkgImportPath,
			},
		},
	}
	checkAndLoadPackage(imports, "decimal.Decimal", "model")
	if len(astPkgs) == 0 {
		t.Fatalf("failed to load module: %s", pkgImportPath)
	}
	notLoadFlag := true
	for _, v := range astPkgs {
		if v.Name == pkgName {
			notLoadFlag = false
		}
	}
	if notLoadFlag {
		t.Fatalf("failed to load module: %s", pkgImportPath)
	}
}

//package model
//
//import (
//"example.com/comm"
//)
//
//type Object struct {
//	Total comm.Common
//}
func TestCheckAndLoadPackageOnGoPath(t *testing.T) {
	var (
		pkgName       = "comm"
		pkgImportPath = "example.com/comm"

		testCommPkg = `
package comm

type Common struct {
	Code  string
	Error string
}
`
	)

	gopath, err := ioutil.TempDir("", "gobuild-gopath")
	if err != nil {
		t.Fatal(err)
	}

	defer os.RemoveAll(gopath)

	if err := os.MkdirAll(filepath.Join(gopath, "src/example.com/comm"), 0777); err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile(filepath.Join(gopath, "src/example.com/comm/comm.go"), []byte(testCommPkg), 0666); err != nil {
		t.Fatal(err)
	}

	defer os.Setenv("GO111MODULE", os.Getenv("GO111MODULE"))
	os.Setenv("GO111MODULE", "off")
	defer os.Setenv("GOPATH", os.Getenv("GOPATH"))
	os.Setenv("GOPATH", gopath)
	build.Default.GOPATH = gopath

	imports := []*ast.ImportSpec{
		{
			Path: &ast.BasicLit{
				Value: pkgImportPath,
			},
		},
	}
	checkAndLoadPackage(imports, "comm.Common", "model")
	if len(astPkgs) == 0 {
		t.Fatalf("failed to load module: %s", pkgImportPath)
	}
	notLoadFlag := true
	for _, v := range astPkgs {
		if v.Name == pkgName {
			notLoadFlag = false
		}
	}
	if notLoadFlag {
		t.Fatalf("failed to load module: %s", pkgImportPath)
	}
}
