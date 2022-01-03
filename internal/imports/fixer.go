package imports

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path"
	"sort"
	"strconv"
)

type fixer struct {
	config FixConfig
}

func newFixer(config FixConfig) *fixer {
	return &fixer{
		config: config,
	}
}

func (f *fixer) Fix(src []byte) ([]byte, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var importDecls []*ast.GenDecl
	imports := make(map[string]string, 8)
	usages := make(map[string]struct{}, 8)
	var walkError error
	ast.Inspect(file, func(n ast.Node) bool {
		if n == nil {
			return false
		}
		switch n := n.(type) {
		case *ast.GenDecl:
			if n.Tok != token.IMPORT {
				return true
			}
			importDecls = append(importDecls, n)

		case *ast.ImportSpec:
			importedPath, err := strconv.Unquote(n.Path.Value)
			if err != nil {
				walkError = fmt.Errorf("unquote %q: %w", n.Path.Value, err)
			}
			var localName string
			if n.Name != nil {
				localName = n.Name.Name
			} else {
				localName = path.Base(importedPath)
			}
			imports[localName] = importedPath

		case *ast.SelectorExpr:
			xident, ok := n.X.(*ast.Ident)
			if !ok {
				return true
			}
			if xident.Obj != nil {
				// If the parser can resolve it, it's not a package ref.
				return true
			}
			pkgName := xident.Name
			usages[pkgName] = struct{}{}
		}
		return true
	})
	if walkError != nil {
		return nil, walkError
	}

	var missingImports []string
	for pkgName := range usages {
		pkgPath := f.config.Packages[pkgName]
		if pkgPath == "" {
			pkgPath = f.config.StdlibPackages[pkgName]
		}
		if pkgPath == "" {
			continue
		}
		missingImports = append(missingImports, pkgPath)
	}
	sort.Strings(missingImports)

	unusedImports := make(map[string]struct{})
	for pkgName, pkgPath := range imports {
		if pkgName == "_" || pkgName == "." {
			continue
		}
		if pkgPath == "C" {
			continue
		}
		_, isUsed := usages[pkgName]
		if !isUsed {
			unusedImports[pkgPath] = struct{}{}
		}
	}

	var buf bytes.Buffer
	continueFrom := 0
	if len(imports) != 0 {
		startFrom := 0
		for i, decl := range importDecls {
			numUsedImports := 0
			for _, spec := range decl.Specs {
				spec := spec.(*ast.ImportSpec)
				importedPath, _ := strconv.Unquote(spec.Path.Value)
				if _, ok := unusedImports[importedPath]; !ok {
					numUsedImports++
				}
			}
			insertImports := i == 0 && len(missingImports) != 0
			isEmpty := !insertImports && numUsedImports == 0
			continueFrom = fset.Position(decl.End()).Offset
			if isEmpty {
				buf.Write(src[startFrom:fset.Position(decl.Pos()).Offset])
				if bytes.HasPrefix(src[continueFrom:], []byte("\n")) {
					continueFrom++
				} else if bytes.HasPrefix(src[continueFrom:], []byte("\r\n")) {
					continueFrom += 2
				}
				startFrom = continueFrom
				continue
			}

			buf.Write(src[startFrom : fset.Position(decl.Pos()).Offset+len("import")])
			buf.WriteString(" (\n")
			for _, spec := range decl.Specs {
				spec := spec.(*ast.ImportSpec)
				importedPath, _ := strconv.Unquote(spec.Path.Value)
				if _, ok := unusedImports[importedPath]; ok {
					continue
				}
				if spec.Name != nil {
					fmt.Fprintf(&buf, "\t%s %s\n", spec.Name.Name, spec.Path.Value)
				} else {
					fmt.Fprintf(&buf, "\t%s\n", spec.Path.Value)
				}
			}
			// Put extra imports into the first import decl.
			if i == 0 {
				for _, pkgPath := range missingImports {
					fmt.Fprintf(&buf, "\t%q\n", pkgPath)
				}
			}
			buf.WriteByte(')')
			startFrom = continueFrom
		}
	} else {
		buf.Write(src[:fset.Position(file.Name.End()).Offset+1])
		buf.WriteString("\nimport (\n")
		for _, pkgPath := range missingImports {
			fmt.Fprintf(&buf, "\t%q\n", pkgPath)
		}
		buf.WriteString(")\n")
		continueFrom = fset.Position(file.Name.End()).Offset + 1
	}
	buf.Write(src[continueFrom:])

	return buf.Bytes(), nil
}
