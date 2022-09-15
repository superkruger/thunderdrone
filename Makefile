backendTest = go test ./... -v -count=1
frontendTest = cd web && npm i && npm test -- --watchAll=false

.PHONY: test
test: test-backend test-frontend
	@echo All tests pass!

.PHONY: test-backend
test-backend:
	$(backendTest)

.PHONY: test-frontend
test-frontend:
	$(frontendTest)
