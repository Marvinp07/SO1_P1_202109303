# Bash - Proyecto SO1

## Imágenes de Docker
- **low**: Contenedor de bajo consumo.
- **ram_high**: Contenedor que consume RAM creando listas grandes en Python.
- **cpu_high**: Contenedor que consume CPU con un bucle infinito.

## Permisos a los archivos .sh
```bash
chmod +x /home/marvinpaz/Documentos/SOPES1/Proyecto/Bash/crear_contenedores.sh
chmod +x /home/marvinpaz/Documentos/SOPES1/Proyecto/Bash/limpiar_contenedores.sh
```

## Construcción de imágenes
```bash
docker build -t low -f Dockerfile.low .
docker build -t ram_high -f Dockerfile.ram_high .
docker build -t cpu_high -f Dockerfile.cpu_high .

