SHELL := /bin/bash

.PHONY: install generate fmt vet test starttestcontainer removetestcontainer spincontainer clearstates kaboom testacc testacccomplete generateresource openlocalwebapi golangcilint tfproviderlint tflint terrafmtlint importfmtlint

default: install

install:
	go mod tidy
	go install .

generate:
	go generate ./...
	go fmt ./...
	go vet ./...

fmt:
	go fmt ./...

vet:
	go vet ./...
	
test:
	go test -parallel=4 ./...

starttestcontainer:
	docker run --name pingfederate_terraform_provider_container \
		-d -p 9031:9031 \
		-d -p 9999:9999 \
		--env-file "${HOME}/.pingidentity/config" \
		pingidentity/pingfederate:2305
# Wait for the instance to become ready
	sleep 1
	duration=0
	while (( duration < 240 )) && ! docker logs pingfederate_terraform_provider_container 2>&1 | grep -q "PingFederate is up"; \
	do \
	    duration=$$((duration+1)); \
		sleep 1; \
	done
# Fail if the container didn't become ready in time
	docker logs pingfederate_terraform_provider_container 2>&1 | grep -q "PingFederate is up" || \
		{ echo "PingFederate container did not become ready in time. Logs:"; docker logs pingfederate_terraform_provider_container; exit 1; }
		
removetestcontainer:
	docker rm -f pingfederate_terraform_provider_container
	
spincontainer: removetestcontainer starttestcontainer

testacc:
	PINGFEDERATE_PROVIDER_HTTPS_HOST=https://localhost:9999 \
	PINGFEDERATE_PROVIDER_USERNAME=administrator \
	PINGFEDERATE_PROVIDER_PASSWORD=2FederateM0re \
	TF_ACC=1 go test -timeout 10m -v ./internal/... -p 4

testacccomplete: spincontainer testacc

clearstates:
	find . -name "*tfstate*" -delete
	
kaboom: clearstates spincontainer install

devcheck: golangcilint tfproviderlint tflint terrafmtlint importfmtlint install kaboom testacc

generateresource:
	PINGFEDERATE_GENERATED_ENDPOINT=oauth/authServerSettings/scopes/exclusiveScopes \
	PINGFEDERATE_RESOURCE_DEFINITION_NAME=ScopeEntry \
	PINGFEDERATE_ALLOW_REQUIRED_BYPASS=False \
	PINGFEDERATE_PUT_ONLY_RESOURCE=False \
	python3 scripts/generate_resource.py
	
openlocalwebapi:
	open "https://localhost:9999/pf-admin-api/api-docs/#/"

golangcilint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run --timeout 5m ./internal/...

tfproviderlint: 
	go run github.com/bflad/tfproviderlint/cmd/tfproviderlintx \
									-c 1 \
									-AT001.ignored-filename-suffixes=_test.go \
									-AT003=false \
									-R009=false \
									-XAT001=false \
									-XR004=false \
									-XS002=false ./internal/...

tflint:
	go run github.com/terraform-linters/tflint --recursive

terrafmtlint:
	find ./internal/acctest -type f -name '*_test.go' \
		| sort -u \
		| xargs -I {} go run github.com/katbyte/terrafmt -f fmt {} -v

importfmtlint:
	go run github.com/pavius/impi/cmd/impi --local . --scheme stdThirdPartyLocal ./internal/...