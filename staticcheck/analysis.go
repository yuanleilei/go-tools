package staticcheck

import (
	"flag"

	"honnef.co/go/tools/internal/passes/buildssa"
	"honnef.co/go/tools/lint/lintdsl"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

func newFlagSet() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.PanicOnError)
	fs.Int("go", 0, "Target minor Go version")
	return fs
}

var Analyzers = map[string]*analysis.Analyzer{
	"SA1000": &analysis.Analyzer{
		Name:     "SA1000",
		Run:      callChecker(checkRegexpRules),
		Doc:      docSA1000,
		Requires: []*analysis.Analyzer{buildssa.Analyzer},
	},
	"SA1001": &analysis.Analyzer{
		Name:     "SA1001",
		Run:      CheckTemplate,
		Doc:      docSA1001,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	},
	"SA1002": &analysis.Analyzer{
		Name:     "SA1002",
		Run:      callChecker(checkTimeParseRules),
		Doc:      docSA1002,
		Requires: []*analysis.Analyzer{buildssa.Analyzer},
	},
	"SA1003": &analysis.Analyzer{
		Name:     "SA1003",
		Run:      callChecker(checkEncodingBinaryRules),
		Doc:      docSA1003,
		Requires: []*analysis.Analyzer{buildssa.Analyzer},
	},
	"SA1004": &analysis.Analyzer{
		Name:     "SA1004",
		Run:      CheckTimeSleepConstant,
		Doc:      docSA1004,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	},
	"SA1005": &analysis.Analyzer{
		Name:     "SA1005",
		Run:      CheckExec,
		Doc:      docSA1005,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	},
	"SA1006": &analysis.Analyzer{
		Name:     "SA1006",
		Run:      CheckUnsafePrintf,
		Doc:      docSA1006,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	},
	"SA1007": &analysis.Analyzer{
		Name:     "SA1007",
		Run:      callChecker(checkURLsRules),
		Doc:      docSA1007,
		Requires: []*analysis.Analyzer{buildssa.Analyzer},
	},
	"SA1008": &analysis.Analyzer{
		Name:     "SA1008",
		Run:      CheckCanonicalHeaderKey,
		Doc:      docSA1008,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	},
	"SA1010": &analysis.Analyzer{
		Name:     "SA1010",
		Run:      callChecker(checkRegexpFindAllRules),
		Doc:      docSA1010,
		Requires: []*analysis.Analyzer{buildssa.Analyzer},
	},
	"SA1011": &analysis.Analyzer{
		Name:     "SA1011",
		Run:      callChecker(checkUTF8CutsetRules),
		Doc:      docSA1011,
		Requires: []*analysis.Analyzer{buildssa.Analyzer},
	},
	"SA1012": &analysis.Analyzer{
		Name:     "SA1012",
		Run:      CheckNilContext,
		Doc:      docSA1012,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	},
	"SA1013": &analysis.Analyzer{
		Name:     "SA1013",
		Run:      CheckSeeker,
		Doc:      docSA1013,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	},
	"SA1014": &analysis.Analyzer{
		Name:     "SA1014",
		Run:      callChecker(checkUnmarshalPointerRules),
		Doc:      docSA1014,
		Requires: []*analysis.Analyzer{buildssa.Analyzer},
	},
	"SA1015": &analysis.Analyzer{
		Name:     "SA1015",
		Run:      CheckLeakyTimeTick,
		Doc:      docSA1015,
		Requires: []*analysis.Analyzer{buildssa.Analyzer},
	},
	"SA1016": &analysis.Analyzer{
		Name:     "SA1016",
		Run:      CheckUntrappableSignal,
		Doc:      docSA1016,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	},
	"SA1017": &analysis.Analyzer{
		Name:     "SA1017",
		Run:      callChecker(checkUnbufferedSignalChanRules),
		Doc:      docSA1017,
		Requires: []*analysis.Analyzer{buildssa.Analyzer},
	},
	"SA1018": &analysis.Analyzer{
		Name:     "SA1018",
		Run:      callChecker(checkStringsReplaceZeroRules),
		Doc:      docSA1018,
		Requires: []*analysis.Analyzer{buildssa.Analyzer},
	},
	// "SA1019": &analysis.Analyzer{Name: "SA1019", Run: CheckDeprecated, Doc: docSA1019},
	"SA1020": &analysis.Analyzer{
		Name:     "SA1020",
		Run:      callChecker(checkListenAddressRules),
		Doc:      docSA1020,
		Requires: []*analysis.Analyzer{buildssa.Analyzer},
	},
	"SA1021": &analysis.Analyzer{
		Name:     "SA1021",
		Run:      callChecker(checkBytesEqualIPRules),
		Doc:      docSA1021,
		Requires: []*analysis.Analyzer{buildssa.Analyzer},
	},
	"SA1023": &analysis.Analyzer{
		Name:     "SA1023",
		Run:      CheckWriterBufferModified,
		Doc:      docSA1023,
		Requires: []*analysis.Analyzer{buildssa.Analyzer},
	},
	"SA1024": &analysis.Analyzer{
		Name:     "SA1024",
		Run:      callChecker(checkUniqueCutsetRules),
		Doc:      docSA1024,
		Requires: []*analysis.Analyzer{buildssa.Analyzer},
	},
	"SA1025": &analysis.Analyzer{
		Name:     "SA1025",
		Run:      CheckTimerResetReturnValue,
		Doc:      docSA1025,
		Requires: []*analysis.Analyzer{buildssa.Analyzer},
	},
	"SA1026": &analysis.Analyzer{
		Name:     "SA1026",
		Run:      callChecker(checkUnsupportedMarshal),
		Doc:      docSA1026,
		Requires: []*analysis.Analyzer{buildssa.Analyzer},
	},
	"SA1027": &analysis.Analyzer{
		Name:     "SA1027",
		Run:      callChecker(checkAtomicAlignment),
		Doc:      docSA1027,
		Requires: []*analysis.Analyzer{buildssa.Analyzer},
	},

	"SA2000": &analysis.Analyzer{
		Name:     "SA2000",
		Run:      CheckWaitgroupAdd,
		Doc:      docSA2000,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	},
	"SA2001": &analysis.Analyzer{
		Name:     "SA2001",
		Run:      CheckEmptyCriticalSection,
		Doc:      docSA2001,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	},
	"SA2002": &analysis.Analyzer{
		Name:     "SA2002",
		Run:      CheckConcurrentTesting,
		Doc:      docSA2002,
		Requires: []*analysis.Analyzer{buildssa.Analyzer},
	},
	"SA2003": &analysis.Analyzer{
		Name:     "SA2003",
		Run:      CheckDeferLock,
		Doc:      docSA2003,
		Requires: []*analysis.Analyzer{buildssa.Analyzer},
	},

	"SA3000": &analysis.Analyzer{
		Name:     "SA3000",
		Run:      CheckTestMainExit,
		Doc:      docSA3000,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	},
	"SA3001": &analysis.Analyzer{
		Name:     "SA3001",
		Run:      CheckBenchmarkN,
		Doc:      docSA3001,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	},

	"SA4000": &analysis.Analyzer{
		Name:     "SA4000",
		Run:      CheckLhsRhsIdentical,
		Doc:      docSA4000,
		Requires: []*analysis.Analyzer{inspect.Analyzer, lintdsl.TokenFileAnalyzer},
	},
	"SA4001": &analysis.Analyzer{
		Name:     "SA4001",
		Run:      CheckIneffectiveCopy,
		Doc:      docSA4001,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	},
	"SA4002": &analysis.Analyzer{
		Name:     "SA4002",
		Run:      CheckDiffSizeComparison,
		Doc:      docSA4002,
		Requires: []*analysis.Analyzer{buildssa.Analyzer},
	},
	"SA4003": &analysis.Analyzer{
		Name:     "SA4003",
		Run:      CheckExtremeComparison,
		Doc:      docSA4003,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	},
	"SA4004": &analysis.Analyzer{
		Name:     "SA4004",
		Run:      CheckIneffectiveLoop,
		Doc:      docSA4004,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	},
	"SA4006": &analysis.Analyzer{
		Name:     "SA4006",
		Run:      CheckUnreadVariableValues,
		Doc:      docSA4006,
		Requires: []*analysis.Analyzer{buildssa.Analyzer},
	},
	"SA4008": &analysis.Analyzer{
		Name:     "SA4008",
		Run:      CheckLoopCondition,
		Doc:      docSA4008,
		Requires: []*analysis.Analyzer{buildssa.Analyzer},
	},
	"SA4009": &analysis.Analyzer{
		Name:     "SA4009",
		Run:      CheckArgOverwritten,
		Doc:      docSA4009,
		Requires: []*analysis.Analyzer{buildssa.Analyzer},
	},
	"SA4010": &analysis.Analyzer{
		Name:     "SA4010",
		Run:      CheckIneffectiveAppend,
		Doc:      docSA4010,
		Requires: []*analysis.Analyzer{buildssa.Analyzer},
	},
	"SA4011": &analysis.Analyzer{
		Name:     "SA4011",
		Run:      CheckScopedBreak,
		Doc:      docSA4011,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	},
	"SA4012": &analysis.Analyzer{
		Name:     "SA4012",
		Run:      CheckNaNComparison,
		Doc:      docSA4012,
		Requires: []*analysis.Analyzer{buildssa.Analyzer},
	},
	"SA4013": &analysis.Analyzer{
		Name:     "SA4013",
		Run:      CheckDoubleNegation,
		Doc:      docSA4013,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	},
	"SA4014": &analysis.Analyzer{
		Name:     "SA4014",
		Run:      CheckRepeatedIfElse,
		Doc:      docSA4014,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	},
	"SA4015": &analysis.Analyzer{
		Name:     "SA4015",
		Run:      callChecker(checkMathIntRules),
		Doc:      docSA4015,
		Requires: []*analysis.Analyzer{buildssa.Analyzer},
	},
	"SA4016": &analysis.Analyzer{
		Name:     "SA4016",
		Run:      CheckSillyBitwiseOps,
		Doc:      docSA4016,
		Requires: []*analysis.Analyzer{buildssa.Analyzer, lintdsl.TokenFileAnalyzer},
	},
	"SA4017": &analysis.Analyzer{
		Name:     "SA4017",
		Run:      CheckPureFunctions,
		Doc:      docSA4017,
		Requires: []*analysis.Analyzer{buildssa.Analyzer},
	},
	"SA4018": &analysis.Analyzer{
		Name:     "SA4018",
		Run:      CheckSelfAssignment,
		Doc:      docSA4018,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	},
	"SA4019": &analysis.Analyzer{
		Name: "SA4019",
		Run:  CheckDuplicateBuildConstraints,
		Doc:  docSA4019,
	},
	"SA4020": &analysis.Analyzer{
		Name:     "SA4020",
		Run:      CheckUnreachableTypeCases,
		Doc:      docSA4020,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	},
	"SA4021": &analysis.Analyzer{
		Name:     "SA4021",
		Run:      CheckSingleArgAppend,
		Doc:      docSA4021,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	},

	"SA5000": &analysis.Analyzer{
		Name:     "SA5000",
		Run:      CheckNilMaps,
		Doc:      docSA5000,
		Requires: []*analysis.Analyzer{buildssa.Analyzer},
	},
	"SA5001": &analysis.Analyzer{
		Name:     "SA5001",
		Run:      CheckEarlyDefer,
		Doc:      docSA5001,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	},
	"SA5002": &analysis.Analyzer{
		Name:     "SA5002",
		Run:      CheckInfiniteEmptyLoop,
		Doc:      docSA5002,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	},
	"SA5003": &analysis.Analyzer{
		Name:     "SA5003",
		Run:      CheckDeferInInfiniteLoop,
		Doc:      docSA5003,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	},
	"SA5004": &analysis.Analyzer{
		Name:     "SA5004",
		Run:      CheckLoopEmptyDefault,
		Doc:      docSA5004,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	},
	"SA5005": &analysis.Analyzer{
		Name:     "SA5005",
		Run:      CheckCyclicFinalizer,
		Doc:      docSA5005,
		Requires: []*analysis.Analyzer{buildssa.Analyzer},
	},
	"SA5007": &analysis.Analyzer{
		Name:     "SA5007",
		Run:      CheckInfiniteRecursion,
		Doc:      docSA5007,
		Requires: []*analysis.Analyzer{buildssa.Analyzer},
	},
	"SA5008": &analysis.Analyzer{
		Name:     "SA5008",
		Run:      CheckStructTags,
		Doc:      `XXX`,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	},
	"SA5009": &analysis.Analyzer{
		Name:     "SA5009",
		Run:      callChecker(checkPrintfRules),
		Doc:      `XXX`,
		Requires: []*analysis.Analyzer{buildssa.Analyzer},
	},

	"SA6000": &analysis.Analyzer{
		Name:     "SA6000",
		Run:      callChecker(checkRegexpMatchLoopRules),
		Doc:      docSA6000,
		Requires: []*analysis.Analyzer{buildssa.Analyzer},
	},
	"SA6001": &analysis.Analyzer{
		Name:     "SA6001",
		Run:      CheckMapBytesKey,
		Doc:      docSA6001,
		Requires: []*analysis.Analyzer{buildssa.Analyzer},
	},
	"SA6002": &analysis.Analyzer{
		Name:     "SA6002",
		Run:      callChecker(checkSyncPoolValueRules),
		Doc:      docSA6002,
		Requires: []*analysis.Analyzer{buildssa.Analyzer},
	},
	"SA6003": &analysis.Analyzer{
		Name:     "SA6003",
		Run:      CheckRangeStringRunes,
		Doc:      docSA6003,
		Requires: []*analysis.Analyzer{buildssa.Analyzer},
	},
	"SA6005": &analysis.Analyzer{
		Name:     "SA6005",
		Run:      CheckToLowerToUpperComparison,
		Doc:      docSA6005,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	},

	"SA9001": &analysis.Analyzer{
		Name:     "SA9001",
		Run:      CheckDubiousDeferInChannelRangeLoop,
		Doc:      docSA9001,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	},
	"SA9002": &analysis.Analyzer{
		Name:     "SA9002",
		Run:      CheckNonOctalFileMode,
		Doc:      docSA9002,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	},
	"SA9003": &analysis.Analyzer{
		Name:     "SA9003",
		Run:      CheckEmptyBranch,
		Doc:      docSA9003,
		Requires: []*analysis.Analyzer{buildssa.Analyzer, lintdsl.TokenFileAnalyzer},
	},
	"SA9004": &analysis.Analyzer{
		Name:     "SA9004",
		Run:      CheckMissingEnumTypesInDeclaration,
		Doc:      docSA9004,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	},
	// Filtering generated code because it may include empty structs generated from data models.
	"SA9005": &analysis.Analyzer{
		Name:     "SA9005",
		Run:      callChecker(checkNoopMarshal),
		Doc:      docSA9005,
		Requires: []*analysis.Analyzer{buildssa.Analyzer},
	},
}
