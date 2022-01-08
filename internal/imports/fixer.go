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

	getLine := func(pos token.Pos) int {
		return fset.Position(pos).Line
	}
	hasImportC := func(decl *ast.GenDecl) bool {
		for _, spec := range decl.Specs {
			spec := spec.(*ast.ImportSpec)
			importedPath, _ := strconv.Unquote(spec.Path.Value)
			if importedPath == "C" {
				return true
			}
		}
		return false
	}

	commentByLine := make(map[int]*ast.CommentGroup, len(file.Comments))
	for _, c := range file.Comments {
		commentByLine[getLine(c.Pos())] = c
	}

	var importDecls []*ast.GenDecl
	imports := make(map[string]string, 8)
	usages := make(map[string]struct{}, 8)
	inlineComments := make(map[*ast.ImportSpec]*ast.CommentGroup)
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
			comment := commentByLine[getLine(n.Pos())]
			if comment != nil {
				inlineComments[n] = comment
			}

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
		if _, ok := imports[pkgPath]; ok {
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

	if len(unusedImports) == 0 && len(missingImports) == 0 {
		return src, nil
	}

	var buf bytes.Buffer
	continueFrom := 0
	if len(imports) != 0 {
		startFrom := 0
		addedImports := false
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
			importedC := hasImportC(decl)
			addParens := (!addedImports && !importedC) || decl.Lparen.IsValid()
			if addParens {
				buf.WriteString(" (\n")
			} else {
				buf.WriteByte(' ')
			}
			prevLine := 0
			if len(decl.Specs) != 0 {
				prevLine = getLine(decl.Specs[0].Pos()) - 1
			}
			for _, spec := range decl.Specs {
				spec := spec.(*ast.ImportSpec)
				if spec.Doc != nil {
					commentFrom := fset.Position(spec.Doc.Pos()).Offset
					commentTo := fset.Position(spec.Doc.End()).Offset
					buf.WriteByte('\t')
					buf.Write(src[commentFrom:commentTo])
					buf.WriteByte('\n')
				}
				line := getLine(spec.Pos())
				if line-prevLine != 1 {
					// Grouping line break. Write extra imports here.
					if !addedImports && !importedC {
						addedImports = true
						for _, pkgPath := range missingImports {
							fmt.Fprintf(&buf, "\t%q\n", pkgPath)
						}
					}
					buf.WriteByte('\n')
				}
				prevLine = line
				importedPath, _ := strconv.Unquote(spec.Path.Value)
				if _, ok := unusedImports[importedPath]; ok {
					continue
				}
				if addParens {
					buf.WriteByte('\t')
				}
				if spec.Name != nil {
					fmt.Fprintf(&buf, "%s %s", spec.Name.Name, spec.Path.Value)
				} else {
					buf.WriteString(spec.Path.Value)
				}
				comment := inlineComments[spec]
				if comment != nil {
					commentFrom := fset.Position(comment.Pos()).Offset
					commentTo := fset.Position(comment.End()).Offset
					buf.WriteByte(' ')
					buf.Write(src[commentFrom:commentTo])
				}
				if addParens {
					buf.WriteByte('\n')
				}
			}
			// Put extra imports into the first import decl,
			// if not done so yet.
			if !addedImports && !importedC {
				addedImports = true
				for _, pkgPath := range missingImports {
					fmt.Fprintf(&buf, "\t%q\n", pkgPath)
				}
			}
			if addParens {
				buf.WriteByte(')')
			}
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
