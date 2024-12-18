FROM docker.io/library/golang:1.23.4-bookworm AS base

COPY /LICENSE /LICENSE
COPY ./ /src

FROM base AS builder

RUN cd /src/cmd/solarwinds-otel-collector && CGO_ENABLED=0 GOEXPERIMENT=boringcrypto go build -trimpath -o /src/bin/solarwinds-otel-collector "-ldflags=-s -w"

FROM builder AS tests
WORKDIR /src
# run tests for go modules of all maintained components
# image build is stopped if test failure is detected
RUN find . -name go.mod -not -path "./cmd/solarwinds-otel-collector/*" -execdir go test ./... \; | ( ! grep FAIL )

FROM gcr.io/distroless/static-debian12@sha256:5c7e2b465ac6a2a4e5f4f7f722ce43b147dabe87cb21ac6c4007ae5178a1fa58

ARG USER_UID=10001
USER ${USER_UID}

COPY --from=tests /src/bin/solarwinds-otel-collector /solarwinds-otel-collector

ENTRYPOINT ["/solarwinds-otel-collector"]
CMD ["--config=/opt/default-config.yaml"]
