#docker run -v $GOPATH:/go -e "GOPATH=/go" tinygo/tinygo:0.40.1 GOOS=js GOARCH=wasm tinygo build -o /go/src/github.com/myuser/myrepo/wasm.wasm --no-debug /go/src/github.com/myuser/myrepo/wasm-main.go

.PHONY: ../static/sELFperf.wasm

../static/sELFperf.wasm: main.go
	GOOS=js GOARCH=wasm tinygo build -o ../static/sELFperf.wasm ./wasm/jsonreport.go
#PATH="/home/paul/opt/go24/bin/:$PATH" make
