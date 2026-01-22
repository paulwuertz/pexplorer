module github.com/paulwuertz/pexplorer/selfperf

go 1.24.0

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/go-delve/delve v1.26.0
	github.com/ianlancetaylor/demangle v0.0.0-20251118225945-96ee0021ea0f
)

require (
	github.com/knightsc/gapstone v4.0.1+incompatible // indirect
)

//replace github.com/paulwuertz/pexplorer/selfperf/symbolextraction => ./symbolextraction
