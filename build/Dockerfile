FROM golang:1.12.9 AS build-env
WORKDIR /usr/local/go/src/github.com/dotmesh-io/dotscience-operator
COPY . /usr/local/go/src/github.com/dotmesh-io/dotscience-operator

ENV GO111MODULE=off

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh curl make

RUN make operator-sdk
RUN make operator-image OPERATOR_IMAGE=registry.connect.redhat.com/storageos/cluster-operator:1.1.0

FROM registry.access.redhat.com/ubi8/ubi-minimal:latest

LABEL name="Dotscience Operator" \
      vendor="Dotscience." \
      version="v0.1.0" \
      release="1" \
      distribution-scope="public" \
      architecture="x86_64" \
      url="https://docs.dotscience.com" \
      io.k8s.description="Dotscience operators runs within a Openshift (or Kubernetes) cluster, and provides a means to deploy and manage ML models, configuring deployments, services and ingresses." \
      io.k8s.display-name="Dotscience Operator" \
      io.openshift.tags="dotscience,ml,operator,scm,git" \
      summary="Dotscience Deployment Operator" \
      description="Dotscience Deployment operators runs within a Openshift (or Kubernetes) cluster, and provides a means to deploy and manage ML models, configuring deployments, services and ingresses."

ENV OPERATOR=/usr/local/bin/dotscience-operator \
    USER_UID=1001 \
    USER_NAME=dotscience-operator

# get licenses
COPY --from=build-env /usr/local/go/src/github.com/dotmesh-io/dotscience-operator/licenses /licenses
# install operator binary
# COPY build/_output/bin/dotscience-operator ${OPERATOR}
# COPY --from=build-env /usr/local/go/src/github.com/dotmesh-io/dotscience-operator ${OPERATOR}
# COPY --from=build-env /usr/local/go/src/github.com/dotmesh-io/dotscience-operator ${OPERATOR}

COPY --from=build-env /go/src/github.com/dotmesh-io/dotscience-operator/build/_output/bin/dotscience-operator ${OPERATOR}

COPY build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
