package loader

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"os"

	"golang.org/x/tools/go/gcexportdata"
	"golang.org/x/tools/go/packages"
)

type Loader struct {
	fromSource map[*packages.Package]struct{}
}

func (ld *Loader) Graph(patterns ...string) ([]*packages.Package, error) {
	cfg := &packages.Config{
		Mode:  packages.NeedName | packages.NeedImports | packages.NeedDeps | packages.NeedExportsFile | packages.NeedFiles | packages.NeedCompiledGoFiles,
		Tests: true,
	}
	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		return nil, err
	}
	fset := token.NewFileSet()
	packages.Visit(pkgs, nil, func(pkg *packages.Package) {
		pkg.Fset = fset
	})
	return pkgs, nil
}

func (ld *Loader) IsFromSource(pkg *packages.Package) bool {
	_, ok := ld.fromSource[pkg]
	return ok
}

func (ld *Loader) LoadFromExport(pkg *packages.Package) error {
	pkg.IllTyped = true
	for path, pkg := range pkg.Imports {
		if pkg.Types == nil {
			return fmt.Errorf("dependency %q hasn't been loaded yet", path)
		}
		if _, ok := ld.fromSource[pkg]; ok {
			return fmt.Errorf("dependency %q was loaded from source", path)
		}
	}
	if pkg.ExportFile == "" {
		return fmt.Errorf("no export data for %q", pkg.ID)
	}
	f, err := os.Open(pkg.ExportFile)
	if err != nil {
		return err
	}
	defer f.Close()

	r, err := gcexportdata.NewReader(f)
	if err != nil {
		return err
	}

	view := make(map[string]*types.Package)  // view seen by gcexportdata
	seen := make(map[*packages.Package]bool) // all visited packages
	var visit func(pkgs map[string]*packages.Package)
	visit = func(pkgs map[string]*packages.Package) {
		for _, pkg := range pkgs {
			if !seen[pkg] {
				seen[pkg] = true
				view[pkg.PkgPath] = pkg.Types
				visit(pkg.Imports)
			}
		}
	}
	visit(pkg.Imports)
	tpkg, err := gcexportdata.Read(r, pkg.Fset, view, pkg.PkgPath)
	if err != nil {
		return err
	}
	pkg.Types = tpkg
	pkg.IllTyped = false
	return nil
}

func (ld *Loader) LoadFromSource(pkg *packages.Package) error {
	pkg.IllTyped = true
	for _, dep := range pkg.Imports {
		if dep.IllTyped {
			return fmt.Errorf("ill-typed dependency %q", dep.ID)
		}
	}
	if ld.fromSource == nil {
		ld.fromSource = map[*packages.Package]struct{}{}
	}
	ld.fromSource[pkg] = struct{}{}
	pkg.Types = types.NewPackage(pkg.PkgPath, pkg.Name)
	// XXX error handling a la go/packages
	// XXX parallelize file parsing
	pkg.Syntax = make([]*ast.File, len(pkg.CompiledGoFiles))
	for i, file := range pkg.CompiledGoFiles {
		f, err := parser.ParseFile(pkg.Fset, file, nil, parser.ParseComments)
		if err != nil {
			return err
		}
		pkg.Syntax[i] = f
	}
	pkg.TypesInfo = &types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Implicits:  make(map[ast.Node]types.Object),
		Scopes:     make(map[ast.Node]*types.Scope),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
	}
	// XXX set pkg.TypeSizes

	importer := func(path string) (*types.Package, error) {
		if path == "unsafe" {
			return types.Unsafe, nil
		}
		// XXX all sorts of error handling
		return pkg.Imports[path].Types, nil
	}
	tc := &types.Config{
		Importer: importerFunc(importer),
		// XXX ignore function bodies
	}
	err := types.NewChecker(tc, pkg.Fset, pkg.Types, pkg.TypesInfo).Files(pkg.Syntax)
	if err != nil {
		return err
	}
	pkg.IllTyped = false
	return nil
}

type importerFunc func(path string) (*types.Package, error)

func (f importerFunc) Import(path string) (*types.Package, error) { return f(path) }
