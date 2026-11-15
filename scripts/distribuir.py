#!/usr/bin/env python3
"""Produz um módulo Go e seu proxy local, sem dependências ou publicação."""

import argparse
import hashlib
import json
import re
import sys
import zipfile
from pathlib import Path, PurePosixPath

SDK = Path(__file__).resolve().parents[1]
MODULO = "github.com/mailcore-hq/mailcore-go"
VERSAO = "v0.1.0-alpha.1"
TEMPO_VERSAO = "2026-11-14T00:00:00Z"
OBRIGATORIOS = {
    "go.mod", "LICENSE", "README.md", "SECURITY.md", "CONTRIBUTING.md", "CHANGELOG.md"
}


def exigir(condicao):
    if not condicao:
        raise ValueError("O contrato do pacote não foi satisfeito.")


def sha256(conteudo):
    return hashlib.sha256(conteudo).hexdigest()


def arquivos_publicos():
    escolhidos = set(OBRIGATORIOS)
    escolhidos.update(caminho.name for caminho in SDK.glob("*.go") if not caminho.name.endswith("_test.go"))
    for pasta, extensao in (("docs", ".md"), ("exemplos", ".go")):
        escolhidos.update(
            caminho.relative_to(SDK).as_posix()
            for caminho in (SDK / pasta).rglob("*" + extensao)
            if caminho.is_file() and not caminho.name.endswith("_test.go")
        )
    exigir(all((SDK / nome).is_file() for nome in OBRIGATORIOS))
    exigir(any(nome.endswith(".go") and "/" not in nome for nome in escolhidos))
    # Regras deliberadamente restritas à nossa árvore ASCII: prefixo seguro,
    # sem links, caminhos especiais ou colisões em sistemas sem distinção de caixa.
    nomes = {}
    conteudos = {}
    for nome in sorted(escolhidos):
        partes = PurePosixPath(nome).parts
        exigir(bool(partes) and not nome.startswith("/") and "\\" not in nome)
        exigir(re.fullmatch(r"[A-Za-z0-9_.\-/]+", nome) is not None)
        for parte in partes:
            exigir(parte not in ("", ".", "..") and not parte.endswith("."))
            exigir(parte.split(".", 1)[0].upper() not in {"CON", "PRN", "AUX", "NUL", "COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9", "LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9"})
        for comprimento in range(1, len(partes) + 1):
            trecho = "/".join(partes[:comprimento])
            equivalente = trecho.casefold()
            exigir(equivalente not in nomes or nomes[equivalente] == trecho)
            nomes[equivalente] = trecho
        caminho = SDK / nome
        exigir(not caminho.is_symlink())
        exigir(not any(pai.is_symlink() for pai in caminho.parents if pai != SDK and SDK in pai.parents))
        conteudo = caminho.read_bytes()
        exigir(len(conteudo) <= 16 * 1024 * 1024)
        conteudos[nome] = conteudo
    for nome in conteudos:
        exigir(not any(pai.as_posix() in conteudos for pai in PurePosixPath(nome).parents if pai.as_posix() != "."))
    exigir(sum(map(len, conteudos.values())) <= 50 * 1024 * 1024)
    go_mod = conteudos["go.mod"].decode("utf-8")
    exigir(re.search(r"(?m)^module\s+" + re.escape(MODULO) + r"\s*$", go_mod) is not None)
    exigir(re.search(r"(?m)^go\s+1\.26(?:\.0)?\s*$", go_mod) is not None)
    exigir(re.search(r"(?m)^\s*(require|replace|exclude|toolchain|tool)\b", go_mod) is None)
    return conteudos


def preparar(destino):
    conteudos = arquivos_publicos()
    prefixo = MODULO + "@" + VERSAO + "/"
    diretorio = destino / "proxy" / MODULO / "@v"
    diretorio.mkdir(parents=True, exist_ok=True)
    pacote = diretorio / (VERSAO + ".zip")
    with zipfile.ZipFile(pacote, "w", compression=zipfile.ZIP_DEFLATED, compresslevel=9) as arquivo:
        for nome, conteudo in conteudos.items():
            entrada = zipfile.ZipInfo(prefixo + nome, date_time=(1980, 1, 1, 0, 0, 0))
            entrada.create_system = 3
            entrada.external_attr = 0o100644 << 16
            entrada.compress_type = zipfile.ZIP_DEFLATED
            arquivo.writestr(entrada, conteudo, compress_type=zipfile.ZIP_DEFLATED, compresslevel=9)
    (diretorio / (VERSAO + ".mod")).write_bytes(conteudos["go.mod"])
    (diretorio / (VERSAO + ".info")).write_text(json.dumps({"Version": VERSAO, "Time": TEMPO_VERSAO}, separators=(",", ":")) + "\n")
    (diretorio / "list").write_text(VERSAO + "\n")
    # O manifesto acompanha o laboratório, sem integrar o próprio ZIP.
    manifesto = {
        "resultado": "preparado", "modulo": MODULO, "versao": VERSAO,
        "prefixo_zip": prefixo, "publicado": False, "dependencias_runtime": 0,
        "zip": pacote.relative_to(destino).as_posix(), "sha256_zip": sha256(pacote.read_bytes()),
        "sha256_mod": sha256(conteudos["go.mod"]),
        "arquivos": {nome: sha256(conteudo) for nome, conteudo in conteudos.items()},
    }
    (destino / "manifesto.json").write_text(json.dumps(manifesto, ensure_ascii=False, indent=2) + "\n")
    return {"resultado": "preparado", "modulo": MODULO, "versao": VERSAO, "arquivos": len(conteudos), "sha256_zip": manifesto["sha256_zip"]}


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--destino", type=Path, default=Path("/pacotes"))
    argumentos = parser.parse_args()
    try:
        resultado = preparar(argumentos.destino.resolve())
    except Exception:
        resultado = {"resultado": "falhou", "etapa": "distribuicao", "falha": "Não foi possível preparar o módulo local; detalhes omitidos."}
    print(json.dumps(resultado, ensure_ascii=False))
    return 0 if resultado["resultado"] == "preparado" else 1


if __name__ == "__main__":
    sys.exit(main())
