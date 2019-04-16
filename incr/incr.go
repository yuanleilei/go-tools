package main

import (
	"go/types"
	"log"
	"os"

	"golang.org/x/tools/go/packages"
	"honnef.co/go/tools/lint"
	"honnef.co/go/tools/lint/lintutil/format"
	"honnef.co/go/tools/loader"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
	"honnef.co/go/tools/unused"
)

func main() {
	ld := &loader.Loader{}
	pkgs, err := ld.Graph("sandbox/bar")
	if err != nil {
		log.Fatal(err)
	}
	if packages.PrintErrors(pkgs) > 0 {
		log.Fatal("Encountered errors")
	}
	initial := map[*packages.Package]struct{}{}
	for _, pkg := range pkgs {
		initial[pkg] = struct{}{}
	}
	packages.Visit(pkgs, nil, func(pkg *packages.Package) {
		if pkg.PkgPath == "unsafe" {
			pkg.Types = types.Unsafe
			return
		}
		var fromSource bool
		if _, ok := initial[pkg]; ok {
			fromSource = true
		} else {
			for _, dep := range pkg.Imports {
				if ld.IsFromSource(dep) {
					fromSource = true
					break
				}
			}
		}
		// XXX also load from source if we need to update our facts
		if fromSource {
			//fmt.Println("Loading", pkg.ID, "from source")
			if err := ld.LoadFromSource(pkg); err != nil {
				log.Println(err)
			}
		} else {
			//fmt.Println("Loading", pkg.ID, "from export data")
			if err := ld.LoadFromExport(pkg); err != nil {
				log.Println(err)
			}
		}
	})

	l := &lint.Linter{
		Checkers: []lint.Checker{
			simple.NewChecker(),
			staticcheck.NewChecker(),
			stylecheck.NewChecker(),
			&unused.Checker{},
		},
		GoVersion: 13,
	}
	p := l.Lint(pkgs, nil)
	f := format.Text{W: os.Stdout}
	for _, pp := range p {
		f.Format(pp)
	}

	if false {
		explicit := make(map[*packages.Package]struct{}, len(pkgs))
		for _, pkg := range pkgs {
			explicit[pkg] = struct{}{}
		}
		packages.Visit(pkgs, nil, func(pkg *packages.Package) {
			log.Println(pkg.PkgPath, pkg.CompiledGoFiles)
			// Step 2: if ExportFile is newer than cached facts, or
			// package was specified explicitly, run analyses, produce
			// facts and diagnoses

			update := false
			if _, ok := explicit[pkg]; ok {
				// always check explicitly provided packages, we want to
				// output diagnoses for them
				update = true
			} else {
				// XXX get mtime of cached facts, get mtime of export
				// file, recompute facts if export file is newer than
				// facts. don't output diagnoses reported for these
				// packages.
			}
			if update {
				// XXX load package from source, run analyses

				// XXX special care must be taken for 0) tests 1) test variants 2)
				// packages specified as a list of files.
				//
				// For tests we can use the import path for
				// loading and the ID for matching. For a list of files,
				// there should only be ever at most one of these
				// packages in the graph. For test variants we're out of luck.
			}
		})
	}
}
