FROM docker.io/library/golang:1.24.1-bookworm AS base

COPY /LICENSE /LICENSE
COPY ./ /src

FROM base AS builder

RUN cd /src/cmd/solarwinds-otel-collector && CGO_ENABLED=0 GOEXPERIMENT=boringcrypto go build -tags full \
    -trimpath -o /src/bin/solarwinds-otel-collector "-ldflags=-s -w"

FROM builder AS tests
WORKDIR /src
# run tests for go modules of all maintained components
# image build is stopped if test failure is detected
RUN find . -name go.mod -not -path "./cmd/solarwinds-otel-collector/*" -execdir go test ./... \; | ( ! grep FAIL )

FROM gcr.io/distroless/static-debian12@sha256:95ea148e8e9edd11cc7f639dc11825f38af86a14e5c7361753c741ceadef2167

ARG USER_UID=10001
USER ${USER_UID}

COPY --from=tests /src/bin/solarwinds-otel-collector /solarwinds-otel-collector

ENTRYPOINT ["/solarwinds-otel-collector"]
CMD ["--config=/opt/default-config.yaml"]
