x ?= A

start_tests_saver:
	go run test_saver.go

all: merge build run_tests

test: build run_tests

merge:
	@go run merger.go $(x)

build:
	@go build -o task_$(x) solutions/task_$(x).go

run_tests:
	@go run tester.go task_$(x) $(x)
	@rm -f task_$(x)
