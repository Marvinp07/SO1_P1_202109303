#!/usr/bin/env bash
set -euo pipefail

CARNET="202109303"

cd "$(dirname "$0")"
make

echo "[+] Cargando módulos..."
sudo insmod sysinfo.ko || true
sudo insmod continfo.ko || true

echo "[+] Verificando /proc entradas..."
SYS_PROC="/proc/sysinfo_so1_${CARNET}"
CONT_PROC="/proc/continfo_so1_${CARNET}"

if [ -r "$SYS_PROC" ]; then
  echo "== $SYS_PROC =="
  head -n 20 "$SYS_PROC"
fi

if [ -r "$CONT_PROC" ]; then
  echo "== $CONT_PROC =="
  head -n 20 "$CONT_PROC"
fi

echo "[+] Listo."
