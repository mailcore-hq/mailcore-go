ARG GO_IMAGE=golang:1.27.4-bookworm@sha256:4fa7c968ba93f1635670522dbb01a2c78150443e1385746924566079f36570ea
FROM ${GO_IMAGE} AS ferramentas
RUN apt-get update && apt-get install -y --no-install-recommends python3 ca-certificates openssl zip unzip && rm -rf /var/lib/apt/lists/*
COPY --from=docker:29-cli@sha256:5a7af0dcb0e46c803efe9b45371dd21cb500e5b4453250394ac6f06d49b6e6f1 /usr/local/bin/docker /usr/local/bin/docker
ENV GOTOOLCHAIN=local GOCACHE=/tmp/go-cache GOPATH=/tmp/go-path GOMAXPROCS=2 GOFLAGS=-p=2
WORKDIR /sdk

FROM ferramentas AS testes
COPY . .
RUN --mount=type=cache,target=/tmp/go-cache go test -run '^$' ./... && mkdir -p /opt/mailcore && go build -trimpath -o /opt/mailcore/integracao-api ./testes/integracao
CMD ["go", "test", "-race", "./..."]

FROM ferramentas AS distribuicao
COPY . .
RUN python3 scripts/distribuir.py --destino /pacotes
CMD ["python3", "scripts/distribuir.py", "--destino", "/pacotes"]
