module github.com/paulwuertz/pexplorer/selfperf

go 1.24

require (
	github.com/go-delve/delve v1.26.0
	github.com/ianlancetaylor/demangle v0.0.0-20251118225945-96ee0021ea0f
)

require github.com/bpfsnoop/gapstone v0.0.0-20260226134052-b57d31fae271 // indirect

//replace github.com/paulwuertz/pexplorer/selfperf/symbolextraction => ./symbolextraction
