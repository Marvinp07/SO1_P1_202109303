#!/bin/bash

# Lista de imágenes disponibles
IMAGENES=("low" "ram_high" "cpu_high")

LOGFILE="/home/marvinpaz/Documentos/SOPES1/Proyecto/Bash/contenedores.log"

# Crear log si no existe
touch "$LOGFILE"

echo "==== $(date) ====" >> "$LOGFILE"

# Crear 10 contenedores aleatorios
for i in {1..10}; do
    IMG=${IMAGENES[$RANDOM % ${#IMAGENES[@]}]}
    NAME="cont_${i}_$(date +%s)"
    
    # Eliminar si ya existe
    docker rm -f "$NAME" >/dev/null 2>&1
    
    # Crear contenedor
    if docker run -d --name "$NAME" "$IMG"; then
        echo "Contenedor $NAME creado con imagen $IMG" >> "$LOGFILE"
    else
        echo "Error al crear contenedor $NAME con imagen $IMG" >> "$LOGFILE"
    fi
done