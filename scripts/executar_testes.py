"""Executa a suíte Go em Docker e registra somente resultados e versões."""

import datetime
import json
import os
from pathlib import Path
import subprocess
import sys

versao = subprocess.check_output(["go", "env", "GOVERSION"], text=True).strip()
comando = ["go", "test", "-race", "-count=1", "-json", "./..."]
processo = subprocess.Popen(comando, stdout=subprocess.PIPE, stderr=subprocess.STDOUT, text=True)
passaram, ignorados, falharam = set(), set(), set()
pacotes = set()
for linha in processo.stdout:
    try:
        evento = json.loads(linha)
    except ValueError:
        print(linha, end="", flush=True)
        continue
    acao, teste, pacote = evento.get("Action"), evento.get("Test"), evento.get("Package", "")
    if acao == "output":
        # São saídas de testes sintéticos. Nunca executar esse runner com dados reais.
        print(evento.get("Output", ""), end="", flush=True)
    if acao == "pass" and teste:
        passaram.add((pacote, teste))
    if acao == "pass" and not teste:
        pacotes.add(pacote)
    if acao == "skip" and teste:
        ignorados.add((pacote, teste))
    if acao == "fail":
        falharam.add((pacote, teste or "pacote"))
codigo = processo.wait()
folhas = [item for item in passaram if not any(outro[0] == item[0] and outro[1].startswith(item[1] + "/") for outro in passaram)]
resultado = {
    "data_utc": datetime.datetime.now(datetime.timezone.utc).isoformat(),
    "go": versao,
    "imagem": os.environ.get("MAILCORE_IMAGEM_TESTE", "nao_informada"),
    "resultado": "aprovado" if codigo == 0 and not falharam and not ignorados and folhas else "reprovado",
    "testes_folha": len(folhas),
    "testes_incluindo_grupos": len(passaram),
    "pacotes_aprovados": len(pacotes),
    "falhas": len(falharam),
    "ignorados": len(ignorados),
    "detector_corridas": True,
    "rede_externa": False,
    "instalacao_host": False,
}
perfil = ".".join(versao.removeprefix("go").split(".")[:2])
Path(f"/relatorios/go-{perfil}.json").write_text(json.dumps(resultado, indent=2, ensure_ascii=False) + "\n")
print(json.dumps(resultado, ensure_ascii=False), flush=True)
sys.exit(0 if resultado["resultado"] == "aprovado" else 1)
