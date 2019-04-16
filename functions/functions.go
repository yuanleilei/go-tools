package functions

import (
	"sync"

	"honnef.co/go/tools/callgraph"
	"honnef.co/go/tools/callgraph/static"
	"honnef.co/go/tools/ssa"
	"honnef.co/go/tools/staticcheck/vrp"
)

var stdlibDescs = map[string]Description{
	"errors.New": {Pure: true},

	"fmt.Errorf":  {Pure: true},
	"fmt.Sprintf": {Pure: true},
	"fmt.Sprint":  {Pure: true},

	"sort.Reverse": {Pure: true},

	"strings.Map":            {Pure: true},
	"strings.Repeat":         {Pure: true},
	"strings.Replace":        {Pure: true},
	"strings.Title":          {Pure: true},
	"strings.ToLower":        {Pure: true},
	"strings.ToLowerSpecial": {Pure: true},
	"strings.ToTitle":        {Pure: true},
	"strings.ToTitleSpecial": {Pure: true},
	"strings.ToUpper":        {Pure: true},
	"strings.ToUpperSpecial": {Pure: true},
	"strings.Trim":           {Pure: true},
	"strings.TrimFunc":       {Pure: true},
	"strings.TrimLeft":       {Pure: true},
	"strings.TrimLeftFunc":   {Pure: true},
	"strings.TrimPrefix":     {Pure: true},
	"strings.TrimRight":      {Pure: true},
	"strings.TrimRightFunc":  {Pure: true},
	"strings.TrimSpace":      {Pure: true},
	"strings.TrimSuffix":     {Pure: true},

	"(*net/http.Request).WithContext": {Pure: true},
}

type Description struct {
	// The function is known to be pure
	Pure bool
	// The function is known to be a stub
	Stub bool
	// The function is known to never return (panics notwithstanding)
	Infinite bool
	// Variable ranges
	Ranges vrp.Ranges
	Loops  []Loop
}

type descriptionEntry struct {
	ready  chan struct{}
	result Description
}

type Descriptions struct {
	CallGraph *callgraph.Graph
	mu        sync.Mutex
	cache     map[*ssa.Function]*descriptionEntry
}

func NewDescriptions(prog *ssa.Program) *Descriptions {
	return &Descriptions{
		CallGraph: static.CallGraph(prog),
		cache:     map[*ssa.Function]*descriptionEntry{},
	}
}

func (d *Descriptions) Get(fn *ssa.Function) Description {
	d.mu.Lock()
	fd := d.cache[fn]
	if fd == nil {
		fd = &descriptionEntry{
			ready: make(chan struct{}),
		}
		d.cache[fn] = fd
		d.mu.Unlock()

		{
			fd.result = stdlibDescs[fn.RelString(nil)]
			fd.result.Pure = fd.result.Pure || d.IsPure(fn)
			fd.result.Stub = fd.result.Stub || d.IsStub(fn)
			fd.result.Infinite = fd.result.Infinite || !terminates(fn)
			fd.result.Ranges = vrp.BuildGraph(fn).Solve()
			fd.result.Loops = findLoops(fn)
		}

		close(fd.ready)
	} else {
		d.mu.Unlock()
		<-fd.ready
	}
	return fd.result
}
