SDK_VERSION = v0.13.0
GO_BUILD_CMD = go build
GO_ENV = GOOS=linux CGO_ENABLED=0
MACHINE = $(shell uname -m)
OPERATOR_IMAGE ?= quay.io/dotmesh/dotscience-operator:dev

OPERATOR_VERSION = 0.1.0
METADATA_FILE = dotscience-olm-metadata.zip

.PHONY: image
image:
	# operator-sdk build quay.io/dotmesh/dotscience-operator
	# docker build 
	docker build -t $(OPERATOR_IMAGE) -f build/Dockerfile .

.PHONY: image-push
image-push: image
	docker push quay.io/dotmesh/dotscience-operator:latest

.PHONY: generate
generate:
	operator-sdk generate k8s

operator-sdk:
	# Download sdk only if it's not available.
	@if [ ! -f build/operator-sdk ]; then \
		curl -Lo build/operator-sdk https://github.com/operator-framework/operator-sdk/releases/download/$(SDK_VERSION)/operator-sdk-$(SDK_VERSION)-$(MACHINE)-linux-gnu && \
		chmod +x build/operator-sdk; \
	fi

operator-image:
	@echo "Building dotscience-operator"
	# build/operator-sdk build $(OPERATOR_IMAGE)
	$(GO_ENV) $(GO_BUILD_CMD) \
		-o ./build/_output/bin/dotscience-operator \
		./cmd/manager

# Install operator on a host. Might fail on containers that don't have sudo.
install-operator-sdk: operator-sdk
	sudo cp build/operator-sdk /usr/local/bin/

gen-csv:
	operator-sdk olm-catalog gen-csv --csv-version $(OPERATOR_VERSION)

metadata-zip:
	# Remove any existing metadata bundle.
	rm -f build/_output/$(METADATA_FILE)
	# Ensure the target path exists.
	mkdir -p build/_output/
	# -j strips the parent directories and adds the files at the root. 
	zip -j build/_output/$(METADATA_FILE) \
		deploy/olm-catalog/dotscience-operator/0.1.0/* \
		deploy/olm-catalog/dotscience-operator/dotscience-operator.package.yaml \
		deploy/crds/deployer.dotscience.com_deployerservices_crd.yaml

metadata-bundle-lint: metadata-zip
	docker run -it --rm -v $(PWD)/build/_output/:/metadata \
		-w /home/test/ \
		python:3 bash -c "pip install operator-courier && unzip /metadata/$(METADATA_FILE) -d out && operator-courier --verbose verify --ui_validate_io out/"


# metadata-lint:
# 	docker run -it --rm -v $(PWD)/build/_output/dotscience-olm-metadata/:/dotscience \
#     python:3 bash -c "pip install operator-courier && operator-courier verify --ui_validate_io /dotscience"