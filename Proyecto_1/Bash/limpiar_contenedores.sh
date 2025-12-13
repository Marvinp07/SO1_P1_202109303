#!/bin/bash
LOGFILE="/home/marvinpaz/Documentos/SOPES1/Proyecto/logs/limpiar.log"

# Crear log si no existe
touch "$LOGFILE"

echo "==== $(date) ====" >> "$LOGFILE"
echo "Eliminando todos los contenedores..." >> "$LOGFILE"

docker rm -f $(docker ps -aq) >> "$LOGFILE" 2>&1
