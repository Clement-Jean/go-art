UNAME_S := $(shell uname -s)
SED :=
ifeq ($(UNAME_S),Linux)
	SED = sed
endif
ifeq ($(UNAME_S),Darwin)
	SED = gsed
endif

all: bench analyze

bench:
	go test -bench=BenchmarkGoART -benchmem -count=10 > art.txt
	go test -bench=BenchmarkGoAdaptiveRadixTree -benchmem -count=10 > gart.txt
	go test -bench=BenchmarkOriginalGoART -benchmem -count=10 > oart.txt

analyze:
	@${SED} -i 's/BenchmarkGoART/Benchmark/g' art.txt
	@${SED} -i 's/BenchmarkGoAdaptiveRadixTree/Benchmark/g' gart.txt
	@${SED} -i 's/BenchmarkOriginalGoART/Benchmark/g' oart.txt
	benchstat oart.txt gart.txt art.txt
