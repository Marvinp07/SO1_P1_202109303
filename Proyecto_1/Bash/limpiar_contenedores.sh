#!/bin/bash

LOGFILE="/var/log/limpiar.log"

# Crear log si no existe
touch "$LOGFILE"

echo "==== $(date) ====" >> "$LOGFILE"
echo "Eliminando contenedores de prueba..." >> "$LOGFILE"

# Solo eliminar los que tengan la etiqueta test_container
docker rm -f $(docker ps -aq --filter "label=test_container") >> "$LOGFILE" 2>&1
