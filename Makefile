

.PHONY: image
image:
	operator-sdk build quay.io/dotmesh/dotscience-operator

.PHONY: image-push
image-push: image
	docker push quay.io/dotmesh/dotscience-operator:latest

.PHONY: generate
generate:
	operator-sdk generate k8s