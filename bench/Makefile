UNAME_S := $(shell uname -s)
SED :=
ifeq ($(UNAME_S),Linux)
	SED = sed
endif
ifeq ($(UNAME_S),Darwin)
	SED = gsed
endif

BENCH_COUNT ?= 10
BENCH_TIME ?= 1s

all: bench analyze

bench:
	go test -bench=BenchmarkGoART -benchmem -count=${BENCH_COUNT} -benchtime=${BENCH_TIME} | tee art.txt
	go test -bench=BenchmarkGoAdaptiveRadixTree -benchmem -count=${BENCH_COUNT} -benchtime=${BENCH_TIME} | tee gart.txt
	go test -bench=BenchmarkOriginalGoART -benchmem -count=${BENCH_COUNT} -benchtime=${BENCH_TIME} | tee oart.txt

analyze:
	@${SED} -i 's/BenchmarkGoART/Benchmark/g' art.txt
	@${SED} -i 's/BenchmarkGoAdaptiveRadixTree/Benchmark/g' gart.txt
	@${SED} -i 's/BenchmarkOriginalGoART/Benchmark/g' oart.txt
	benchstat oart.txt gart.txt art.txt
