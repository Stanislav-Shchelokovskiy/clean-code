bin-deps:
	go install github.com/vektra/mockery/v3@v3.5.4

gen-mocks:
	mockery
	$(info Tidying module requirements...)
	go mod tidy

.PHONY: \
	gen-mocks \
	bin-deps
