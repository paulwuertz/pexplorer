.PHONY: ../static/sELFperf.wasm test testdataReportOutput jsonvalidate goprepublish gopublish stackcheck

../static/sELFperf.wasm: main.go
	GOOS=js GOARCH=wasm go build -o ../static/sELFperf.wasm ./wasm/jsonreport.go

test:
	go test -bench=. ./cmd/jsonreport/

testdataReportOutput:
	mkdir -p tmp
	go run ./cmd/jsonreport/jsonreport.go -i ../testdata/elf_testdata/zswatch_nrf5340_070.elf -o tmp/zswatch_nrf5340_070.json -f tmp/zswatch_nrf5340_070_functions.json -p
	go run ./cmd/jsonreport/jsonreport.go -i ../testdata/elf_testdata/zswatch_nrf5340_080.elf -o tmp/zswatch_nrf5340_080.json -f tmp/zswatch_nrf5340_080_functions.json -p
	go run ./cmd/jsonreport/jsonreport.go -i ../testdata/elf_testdata/zswatch_nrf5340_081.elf -o tmp/zswatch_nrf5340_081.json -f tmp/zswatch_nrf5340_081_functions.json -p
	go run ./cmd/jsonreport/jsonreport.go -i ../testdata/elf_testdata/Pinecilv2_EN_v2_23.elf -o tmp/Pinecilv2_EN_v2_23.json -p
	go run ./cmd/jsonreport/jsonreport.go -i ../testdata/elf_testdata/Pinecilv1_EN_v2_23.elf -o tmp/Pinecilv1_EN_v2_23.json -p
	go run ./cmd/jsonreport/jsonreport.go -i ../testdata/elf_testdata/zephyr_cannectivity_12_llvm_lpc55s16.elf -o tmp/zephyr_cannectivity_12_llvm_lpc55s16.json -f tmp/zephyr_cannectivity_12_llvm_lpc55s16_functions.json -p
	go run ./cmd/jsonreport/jsonreport.go -i ../testdata/elf_testdata/zephyr_cannectivity_13_llvm_lpc55s16.elf -o tmp/zephyr_cannectivity_13_llvm_lpc55s16.json -f tmp/zephyr_cannectivity_13_llvm_lpc55s16_functions.json -p
	go run ./cmd/jsonreport/jsonreport.go -i ../testdata/elf_testdata/zephyr_cannectivity_13_gcc_lpc55s16.elf -o tmp/zephyr_cannectivity_13_gcc_lpc55s16.json -f tmp/zephyr_cannectivity_13_gcc_lpc55s16_functions.json -p
	go run ./cmd/jsonreport/jsonreport.go -i ../testdata/elf_testdata/zephyr_cannectivity_12_gcc_lpc55s16.elf -o tmp/zephyr_cannectivity_12_gcc_lpc55s16.json -f tmp/zephyr_cannectivity_12_gcc_lpc55s16_functions.json -p
	go run ./cmd/jsonreport/jsonreport.go -i ../testdata/elf_testdata/cannectivity_lpcxpresso55s16_1.4.elf -o tmp/zephyr_cannectivity_14_gcc_lpc55s16.json -f tmp/zephyr_cannectivity_14_gcc_lpc55s16_functions.json -p
	go run ./cmd/jsonreport/jsonreport.go -i ../testdata/elf_testdata/prusa_buddy_boot_64.elf -f tmp/prusa_buddy_boot_64_functions.elf -o tmp/prusa_buddy_boot_64.json -p
	go run ./cmd/jsonreport/jsonreport.go -i ../testdata/elf_testdata/libre_solar_zephyr_43_b953.elf -o tmp/libre_solar_zephyr_43_b953.elf.json -p

jsonvalidate:
	jsonschema validate --verbose  schema.json tmp/*.json

goprepublish:
	go mod tidy
	go test ./...

gopublish:
	# before
	# git tag v0.1.0
	# git push origin v0.1.0
	GOPROXY=proxy.golang.org go list -m github.com/paulwuertz/pexplorer/selfperf@v0.1.0

stackcheck:
	go run ./cmd/stackcheck/stackcheck.go -i ../testdata/elf_testdata/cannectivity_lpcxpresso55s16_1.4.elf
