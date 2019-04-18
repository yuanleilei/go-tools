package stylecheck

import (
	"flag"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"honnef.co/go/tools/config"
	"honnef.co/go/tools/internal/passes/buildssa"
	"honnef.co/go/tools/lint"
)

func newFlagSet() flag.FlagSet {
	statusCodes := append(lint.StringSliceVar{}, config.DefaultConfig.HTTPStatusCodeWhitelist...)
	initialisms := append(lint.StringSliceVar{}, config.DefaultConfig.Initialisms...)
	dotImports := append(lint.StringSliceVar{}, config.DefaultConfig.DotImportWhitelist...)

	fs := flag.NewFlagSet("", flag.PanicOnError)
	fs.Int("go", 0, "Target minor Go version")
	fs.Var(&dotImports, "dot-import-whitelist", "")
	fs.Var(&statusCodes, "http-status-code-whitelist", "")
	fs.Var(&initialisms, "initialisms", "")
	return *fs
}

var Analyzers = map[string]*analysis.Analyzer{
	"ST1000": {
		Name:  "ST1000",
		Run:   CheckPackageComment,
		Doc:   docST1000,
		Flags: newFlagSet(),
	},
	"ST1001": {
		Name:  "ST1001",
		Run:   CheckDotImports,
		Doc:   docST1001,
		Flags: newFlagSet(),
	},
	"ST1003": {
		Name:  "ST1003",
		Run:   CheckNames,
		Doc:   docST1003,
		Flags: newFlagSet(),
	},
	"ST1005": {
		Name:     "ST1005",
		Run:      CheckErrorStrings,
		Doc:      docST1005,
		Requires: []*analysis.Analyzer{buildssa.Analyzer},
		Flags:    newFlagSet(),
	},
	"ST1006": {
		Name:     "ST1006",
		Run:      CheckReceiverNames,
		Doc:      docST1006,
		Requires: []*analysis.Analyzer{buildssa.Analyzer},
		Flags:    newFlagSet(),
	},
	"ST1008": {
		Name:     "ST1008",
		Run:      CheckErrorReturn,
		Doc:      docST1008,
		Requires: []*analysis.Analyzer{buildssa.Analyzer},
		Flags:    newFlagSet(),
	},
	"ST1011": {
		Name:  "ST1011",
		Run:   CheckTimeNames,
		Doc:   docST1011,
		Flags: newFlagSet(),
	},
	"ST1012": {
		Name:  "ST1012",
		Run:   CheckErrorVarNames,
		Doc:   docST1012,
		Flags: newFlagSet(),
	},
	"ST1013": {
		Name:  "ST1013",
		Run:   CheckHTTPStatusCodes,
		Doc:   docST1013,
		Flags: newFlagSet(),
	},
	"ST1015": {
		Name:     "ST1015",
		Run:      CheckDefaultCaseOrder,
		Doc:      docST1015,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
		Flags:    newFlagSet(),
	},
	"ST1016": {
		Name:     "ST1016",
		Run:      CheckReceiverNamesIdentical,
		Doc:      docST1016,
		Requires: []*analysis.Analyzer{buildssa.Analyzer},
		Flags:    newFlagSet(),
	},
	"ST1017": {
		Name:     "ST1017",
		Run:      CheckYodaConditions,
		Doc:      docST1017,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
		Flags:    newFlagSet(),
	},
	"ST1018": {
		Name:     "ST1018",
		Run:      CheckInvisibleCharacters,
		Doc:      docST1018,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
		Flags:    newFlagSet(),
	},
}
