#!/usr/bin/env python3
"""Instala o ZIP por GOPROXY local e verifica um consumidor sem checkout da biblioteca."""

import argparse
import hashlib
import json
import os
import re
import shutil
import subprocess
import sys
import tempfile
from datetime import datetime, timezone
from pathlib import Path

SDK = Path(__file__).resolve().parents[1]
MODULO = "github.com/mailcore-hq/mailcore-go"
VERSAO = "v0.1.0-alpha.1"
PERFIS = {"1.26", "1.27"}


def exigir(condicao):
    if not condicao:
        raise ValueError("O contrato de instalação não foi satisfeito.")


def executar(argumentos, ambiente, pasta):
    retorno = subprocess.run(argumentos, cwd=pasta, env=ambiente, stdout=subprocess.PIPE, stderr=subprocess.PIPE, timeout=180, check=False)
    exigir(retorno.returncode == 0 and len(retorno.stdout) <= 1024 * 1024)
    return retorno.stdout


def salvar(destino, documento):
    destino.parent.mkdir(parents=True, exist_ok=True)
    with os.fdopen(os.open(destino, os.O_WRONLY | os.O_CREAT | os.O_TRUNC | os.O_NOFOLLOW, 0o644), "w") as arquivo:
        json.dump(documento, arquivo, ensure_ascii=False, indent=2)
        arquivo.write("\n")


def validar_fontes_instalados(pasta, manifesto):
    arquivos = manifesto["arquivos"]
    exigir(set(("go.mod", "LICENSE", "README.md", "SECURITY.md", "CONTRIBUTING.md", "CHANGELOG.md")) <= set(arquivos))
    presentes = {item.relative_to(pasta).as_posix() for item in pasta.rglob("*") if item.is_file()}
    exigir(presentes == set(arquivos))
    for nome, esperado in arquivos.items():
        conteudo = (pasta / nome).read_bytes()
        exigir(hashlib.sha256(conteudo).hexdigest() == esperado)
    return True


def instalar(artefatos, perfil):
    manifesto = json.loads((artefatos / "manifesto.json").read_bytes())
    exigir(manifesto["modulo"] == MODULO and manifesto["versao"] == VERSAO and manifesto["dependencias_runtime"] == 0)
    pacote = (artefatos / manifesto["zip"]).resolve()
    exigir(pacote.is_relative_to(artefatos) and pacote.is_file())
    exigir(hashlib.sha256(pacote.read_bytes()).hexdigest() == manifesto["sha256_zip"])
    with tempfile.TemporaryDirectory(prefix="mailcore-go-consumidor-") as temporario:
        pasta = Path(temporario)
        ambiente = os.environ.copy()
        ambiente.update({
            "GOPROXY": (artefatos / "proxy").as_uri(), "GOSUMDB": "off", "GONOSUMDB": "*",
            "GONOPROXY": "none", "GOPRIVATE": "", "GOTOOLCHAIN": "local", "GOWORK": "off",
            "GOENV": "off", "GOFLAGS": "", "GOMODCACHE": str(pasta / "modulos"),
            "GOCACHE": str(pasta / "compilacao"), "GOPATH": str(pasta / "gopath"),
        })
        versao_go = executar(["go", "version"], ambiente, pasta).decode().strip()
        encontrado = re.search(r"\bgo(1\.\d+)(?:\.\d+)?\b", versao_go)
        exigir(encontrado is not None and encontrado.group(1) in PERFIS)
        perfil_real = encontrado.group(1)
        exigir(perfil is None or perfil == perfil_real)
        (pasta / "go.mod").write_text("module exemplo.test/consumidor\n\ngo 1.26.0\n\nrequire " + MODULO + " " + VERSAO + "\n")
        shutil.copyfile(SDK / "scripts/verificar_instalacao/main.go", pasta / "main.go")
        baixado = json.loads(executar(["go", "mod", "download", "-json", MODULO + "@" + VERSAO], ambiente, pasta))
        exigir(baixado.get("Version") == VERSAO and not baixado.get("Error"))
        modulo_pasta = Path(baixado["Dir"]).resolve()
        exigir(modulo_pasta.is_relative_to(pasta / "modulos") and not modulo_pasta.is_relative_to(SDK))
        exigir(validar_fontes_instalados(modulo_pasta, manifesto))
        executar(["go", "mod", "tidy"], ambiente, pasta)
        grafo = executar(["go", "list", "-m", "all"], ambiente, pasta).decode().splitlines()
        exigir(grafo == ["exemplo.test/consumidor", MODULO + " " + VERSAO])
        resultado = json.loads(executar(["go", "run", "-mod=readonly", "."], ambiente, pasta))
        exigir(resultado.get("resultado") == "aprovado")
        exigir(resultado.get("recursos") == 30 and resultado.get("metodos") == 151 and resultado.get("callbacks_fluxo") == 3)
        resultado.update({
            "data_utc": datetime.now(timezone.utc).isoformat(), "perfil": perfil_real,
            "modulo": MODULO, "versao": VERSAO, "sha256_zip": manifesto["sha256_zip"],
            "dependencias_runtime": 0, "instalacao_host": False, "rede_utilizada": False,
            "proxy_arquivo_local": True, "modulo_fora_checkout": True,
            "fontes_licenca_guias_preservados": True, "publicado": False,
        })
        return resultado


def agregar(pasta):
    resultados = [json.loads(caminho.read_bytes()) for caminho in sorted(pasta.glob("instalacao-go-*.json"))]
    exigir(len(resultados) == len(PERFIS))
    exigir({item.get("perfil") for item in resultados} == PERFIS)
    exigir(all(item.get("resultado") == "aprovado" and item.get("recursos") == 30 and item.get("metodos") == 151 and item.get("callbacks_fluxo") == 3 for item in resultados))
    exigir(len({item.get("sha256_zip") for item in resultados}) == 1)
    return {"resultado": "aprovado", "data_utc": datetime.now(timezone.utc).isoformat(), "modulo": MODULO, "versao": VERSAO, "publicado": False, "perfis": resultados}


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--artefatos", type=Path, default=Path("/pacotes"))
    parser.add_argument("--relatorio", type=Path)
    parser.add_argument("--perfil", choices=sorted(PERFIS))
    parser.add_argument("--agregar", type=Path)
    opcoes = parser.parse_args()
    destino = opcoes.relatorio
    try:
        if opcoes.agregar:
            resultado = agregar(opcoes.agregar.resolve())
            destino = destino or opcoes.agregar / "instalacao-local.json"
        else:
            exigir(destino is not None)
            resultado = instalar(opcoes.artefatos.resolve(), opcoes.perfil)
    except Exception:
        resultado = {"resultado": "falhou", "etapa": "instalacao", "perfil": opcoes.perfil, "falha": "O consumidor instalado não satisfez o contrato; detalhes omitidos."}
    if destino is not None:
        try:
            salvar(destino, resultado)
        except Exception:
            resultado = {"resultado": "falhou", "etapa": "relatorio", "falha": "Não foi possível registrar o ensaio local."}
    print(json.dumps(resultado, ensure_ascii=False))
    return 0 if resultado["resultado"] == "aprovado" else 1


if __name__ == "__main__":
    sys.exit(main())
