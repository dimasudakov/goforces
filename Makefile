x ?= A

start_tests_saver:
	go run ./service/task_manager.go

all: merge build run_tests

test: build run_tests

merge:
	@go run ./service/merger.go $(x)

build:
	@go build -o task_$(x) solutions/task_$(x).go

run_tests:
	@go run ./service/tester.go task_$(x) $(x)
	@rm -f task_$(x)

clear:
	@rm -f task_*.go
	@rm -f ./tests/*
	@rm -f ./solutions/*

