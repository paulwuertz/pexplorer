.PHONY: ../static/sELFperf.wasm test

../static/sELFperf.wasm: main.go
	GOOS=js GOARCH=wasm go build -o ../static/sELFperf.wasm ./wasm/jsonreport.go

test:
	go test -bench=. ./cmd/jsonreport/

testdataReportOutput:
	mkdir -p tmp
	go run ./cmd/jsonreport/jsonreport.go -i ../testdata/elf_testdata/zswatch_nrf5340_07.elf -o tmp/zswatch_nrf5340_07.json -p
	go run ./cmd/jsonreport/jsonreport.go -i ../testdata/elf_testdata/Pinecilv2_EN_v2_23.elf -o tmp/Pinecilv2_EN_v2_23.json -p
	go run ./cmd/jsonreport/jsonreport.go -i ../testdata/elf_testdata/Pinecilv1_EN_v2_23.elf -o tmp/Pinecilv1_EN_v2_23.json -p
	go run ./cmd/jsonreport/jsonreport.go -i ../testdata/elf_testdata/zephyr_cannectivity_12_llvm_lpc55s16.elf -o tmp/zephyr_cannectivity_12_llvm_lpc55s16.json -p
	go run ./cmd/jsonreport/jsonreport.go -i ../testdata/elf_testdata/zephyr_cannectivity_13_llvm_lpc55s16.elf -o tmp/zephyr_cannectivity_13_llvm_lpc55s16.json -p
	go run ./cmd/jsonreport/jsonreport.go -i ../testdata/elf_testdata/zephyr_cannectivity_13_gcc_lpc55s16.elf -o tmp/zephyr_cannectivity_13_gcc_lpc55s16.json -p
	go run ./cmd/jsonreport/jsonreport.go -i ../testdata/elf_testdata/zephyr_cannectivity_12_gcc_lpc55s16.elf -o tmp/zephyr_cannectivity_12_gcc_lpc55s16.json -p
	go run ./cmd/jsonreport/jsonreport.go -i ../testdata/elf_testdata/prusa_buddy_boot_64.elf -o tmp/prusa_buddy_boot_64.json -p

jsonvalidate:
	jsonschema validate --verbose  schema.json tmp/*.json
