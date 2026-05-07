module github.com/paulwuertz/pexplorer/testdata

go 1.25.3

require github.com/paulwuertz/pexplorer/selfperf v0.0.0-20260218094942-abdd34d37a06

require (
	github.com/go-delve/delve v1.26.0 // indirect
	github.com/ianlancetaylor/demangle v0.0.0-20251118225945-96ee0021ea0f // indirect
	github.com/knightsc/gapstone v4.0.1+incompatible // indirect
)

replace github.com/paulwuertz/pexplorer/selfperf v0.0.0-20260218094942-abdd34d37a06 => ../selfperf

replace github.com/paulwuertz/pexplorer/selfperf/symbolextraction v0.0.0-20260218094942-abdd34d37a06 => ../selfperf/symbolextraction

replace github.com/paulwuertz/pexplorer/selfperf/callgraph v0.0.0-20260218094942-abdd34d37a06 => ../selfperf/callgraph
