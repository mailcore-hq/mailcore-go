#!/bin/sh
set -eu
raiz=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
mkdir -p "$raiz/docs" "$raiz/dist"
for perfil in 1.26 1.27; do
  case "$perfil" in
    1.26) base=golang:1.26.11-bookworm@sha256:dc40a13b61304647ed7e5d581d246b86413dfc29d63c6d018310feb93cd98fc0 ;;
    1.27) base=golang:1.27.4-bookworm@sha256:4fa7c968ba93f1635670522dbb01a2c78150443e1385746924566079f36570ea ;;
  esac
  docker build --target testes --build-arg "GO_IMAGE=$base" --tag "mailcore/go-testes:$perfil" "$raiz"
  imagem=$(docker image inspect "mailcore/go-testes:$perfil" --format '{{.Id}}')
  docker run --rm --network none --read-only --cap-drop ALL \
    --security-opt no-new-privileges:true --memory 2g --pids-limit 256 \
    --tmpfs /tmp:rw,nosuid,exec,size=1g,mode=1777 --user "$(id -u):$(id -g)" \
    --env GOMAXPROCS=2 --env GOFLAGS=-p=2 --env "MAILCORE_IMAGEM_TESTE=$imagem" \
    --mount "type=bind,src=$raiz/docs,dst=/relatorios" \
    "$imagem" python3 scripts/executar_testes.py
done
docker run --rm --network none --read-only --tmpfs /tmp:rw,nosuid,exec,size=1g \
  --user "$(id -u):$(id -g)" --env GOMAXPROCS=2 --env GOFLAGS=-p=2 \
  mailcore/go-testes:1.27 sh -c \
  'test -z "$(gofmt -l .)" && go vet ./... && python3 scripts/gerar_recursos.py --verificar'
docker run --rm --network none --read-only --tmpfs /tmp:rw,nosuid,exec,size=1g \
  --user "$(id -u):$(id -g)" --env GOMAXPROCS=2 --env GOFLAGS=-p=2 \
  --mount "type=bind,src=$raiz/docs,dst=/relatorios" \
  mailcore/go-testes:1.27 python3 scripts/verificar_plataformas.py --relatorio=/relatorios/plataformas.json
docker build --target distribuicao --tag mailcore/go-distribuicao:local "$raiz"
container=$(docker create mailcore/go-distribuicao:local)
trap 'docker rm --force "$container" >/dev/null 2>&1 || true' EXIT HUP INT TERM
docker cp "$container:/pacotes/." "$raiz/dist/"
docker rm "$container" >/dev/null
trap - EXIT HUP INT TERM
for perfil in 1.26 1.27; do
  docker run --rm --network none --read-only --tmpfs /tmp:rw,nosuid,exec,size=1g,mode=1777 \
    --user "$(id -u):$(id -g)" --env GOMAXPROCS=2 --env GOFLAGS=-p=2 \
    --mount "type=bind,src=$raiz/dist,dst=/pacotes,readonly" \
    --mount "type=bind,src=$raiz/docs,dst=/relatorios" "mailcore/go-testes:$perfil" \
    python3 scripts/instalar_consumidor.py --artefatos=/pacotes \
    "--relatorio=/relatorios/instalacao-go-$perfil.json" "--perfil=$perfil"
done
docker run --rm --network none --read-only --tmpfs /tmp:rw,size=64m \
  --user "$(id -u):$(id -g)" --mount "type=bind,src=$raiz/docs,dst=/relatorios" \
  mailcore/go-testes:1.27 python3 scripts/instalar_consumidor.py --agregar=/relatorios
