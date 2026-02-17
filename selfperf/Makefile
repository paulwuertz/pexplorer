.PHONY: ../static/sELFperf.wasm test

../static/sELFperf.wasm: main.go
	GOOS=js GOARCH=wasm go build -o ../static/sELFperf.wasm ./wasm/jsonreport.go

test:
	go test -bench=. ./cmd/jsonreport/
