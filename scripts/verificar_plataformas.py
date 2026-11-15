#!/usr/bin/env python3
"""Compila o núcleo Go para seis alvos; não executa binários nesses sistemas."""

import argparse
import datetime
import json
import os
from pathlib import Path
import subprocess
import sys
import time


SDK = Path(__file__).resolve().parents[1]
PERFIS = (
    ("linux", "amd64"),
    ("linux", "arm64"),
    ("darwin", "amd64"),
    ("darwin", "arm64"),
    ("windows", "amd64"),
    ("windows", "arm64"),
)


def comando(argumentos, ambiente, prazo=300):
    return subprocess.run(
        argumentos,
        cwd=SDK,
        env=ambiente,
        stdin=subprocess.DEVNULL,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        timeout=prazo,
        check=False,
    )


def verificar():
    ambiente = os.environ.copy()
    ambiente.update({
        "CGO_ENABLED": "0", "GOTOOLCHAIN": "local", "GOENV": "off",
        "GOWORK": "off", "GOPROXY": "off", "GOSUMDB": "off",
        "GOMAXPROCS": "2", "GOFLAGS": "-p=2",
        "GOCACHE": "/tmp/mailcore-go-plataformas-cache",
        "GOPATH": "/tmp/mailcore-go-plataformas-path",
    })
    # Não herdar um alvo eventual da sessão que iniciou o laboratório.
    ambiente.pop("GOOS", None)
    ambiente.pop("GOARCH", None)
    versoes = comando(["go", "env", "-json", "GOVERSION", "GOHOSTOS", "GOHOSTARCH"], ambiente, 30)
    compilador = comando(["go", "tool", "compile", "-V"], ambiente, 30)
    if versoes.returncode != 0 or compilador.returncode != 0:
        raise ValueError("Não foi possível consultar o compilador.")
    informacoes = json.loads(versoes.stdout)
    if informacoes["GOVERSION"] != "go1.27.4":
        raise ValueError("Use o compilador Go 1.27.4 deste ensaio.")
    resultados = []
    for sistema, arquitetura in PERFIS:
        alvo = {**ambiente, "GOOS": sistema, "GOARCH": arquitetura}
        inicio = time.monotonic()
        try:
            retorno = comando(["go", "build", "."], alvo)
            aprovado = retorno.returncode == 0
        except subprocess.TimeoutExpired:
            aprovado = False
        resultado = {
            "sistema": sistema,
            "arquitetura": arquitetura,
            "resultado": "aprovado" if aprovado else "reprovado",
            "duracao_segundos": round(time.monotonic() - inicio, 3),
            "executado_no_alvo": False,
        }
        if not aprovado:
            resultado["falha"] = "Compilação não concluída; saídas do compilador omitidas."
        resultados.append(resultado)
        print(json.dumps(resultado, ensure_ascii=False), flush=True)
    interfaces = Path("/sys/class/net")
    somente_loopback = interfaces.is_dir() and {item.name for item in interfaces.iterdir()} == {"lo"}
    return {
        "data_utc": datetime.datetime.now(datetime.timezone.utc).isoformat(),
        "resultado": "aprovado" if all(item["resultado"] == "aprovado" for item in resultados) else "reprovado",
        "go": informacoes["GOVERSION"],
        "compilador": compilador.stdout.decode().strip(),
        "sistema_compilador": informacoes["GOHOSTOS"],
        "arquitetura_compilador": informacoes["GOHOSTARCH"],
        "imagem": os.environ.get("MAILCORE_IMAGEM_TESTE", "nao_informada"),
        "docker": Path("/.dockerenv").is_file(),
        "interfaces_somente_loopback": somente_loopback,
        "cgo_habilitado": False,
        "comando": ["go", "build", "."],
        "escopo": "compilacao_cruzada_nucleo_sem_execucao_nos_alvos",
        "perfis": resultados,
    }


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--relatorio", type=Path, default=SDK / "docs/plataformas.json")
    opcoes = parser.parse_args()
    try:
        resultado = verificar()
    except Exception:
        resultado = {
            "resultado": "reprovado",
            "escopo": "compilacao_cruzada_nucleo_sem_execucao_nos_alvos",
            "falha": "Não foi possível concluir a verificação; detalhes omitidos.",
        }
    try:
        opcoes.relatorio.parent.mkdir(parents=True, exist_ok=True)
        with os.fdopen(os.open(opcoes.relatorio, os.O_WRONLY | os.O_CREAT | os.O_TRUNC | os.O_NOFOLLOW, 0o644), "w") as arquivo:
            json.dump(resultado, arquivo, ensure_ascii=False, indent=2)
            arquivo.write("\n")
    except Exception:
        resultado = {"resultado": "reprovado", "falha": "Não foi possível registrar a verificação."}
    print(json.dumps(resultado, ensure_ascii=False), flush=True)
    return 0 if resultado["resultado"] == "aprovado" else 1


if __name__ == "__main__":
    sys.exit(main())
