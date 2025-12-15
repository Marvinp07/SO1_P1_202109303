#!/bin/bash

# Lista de imágenes disponibles
IMAGENES=("low" "ram_high" "cpu_high")

LOGFILE="/var/log/contenedores.log"

# Crear log si no existe
touch "$LOGFILE"

echo "==== $(date) ====" >> "$LOGFILE"

# Limitar a máximo 30 contenedores activos
MAX_CONT=30
ACTIVOS=$(docker ps -aq --filter "name=cont_" | wc -l)

if [ "$ACTIVOS" -ge "$MAX_CONT" ]; then
    echo "Límite de contenedores alcanzado ($ACTIVOS). No se crean nuevos." >> "$LOGFILE"
    exit 0
fi

# Crear 10 contenedores aleatorios
for i in {1..10}; do
    IMG=${IMAGENES[$RANDOM % ${#IMAGENES[@]}]}
    NAME="cont_${i}_$(date +%s)"

    # Eliminar si ya existe
    docker rm -f "$NAME" >/dev/null 2>&1

    # Crear contenedor con etiqueta para limpieza
    if docker run -d --name "$NAME" --label test_container "$IMG"; then
        echo "Contenedor $NAME creado con imagen $IMG" >> "$LOGFILE"
    else
        echo "Error al crear contenedor $NAME con imagen $IMG" >> "$LOGFILE"
    fi
done
