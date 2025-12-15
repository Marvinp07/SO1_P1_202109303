# Manual Técnico – Proyecto 1

## 1. Estructura del Módulo

### 1.1 Organización de archivos y directorios
- `/Bash/`
  - `cronjob.sh` → script que genera 10 contenedores cada minuto.
  - `load_modules.sh` → script para compilar y cargar los módulos de kernel.
  - Otros scripts auxiliares para automatización.

- `/Dashboard/`
  - `docker-compose.yml` → orquesta Grafana y SQLite.
  - Configuración de volúmenes (`./data`, `./grafana-data`).
  - Archivos de configuración de Grafana (si aplica).

- `/go-daemon/`
  - `main.go` → código principal del daemon.
  - `db.go` → funciones para guardar métricas en SQLite.
  - `analyzer.go` → lógica de análisis de métricas y eliminación de contenedores.
  - `proc_reader.go` → funciones para leer y parsear `/proc`.

- `/Kernel/`
  - `Makefile` → reglas de compilación del módulo.
  - `continfo.c` → módulo para procesos de contenedores.
  - `sysinfo.c` → módulo para procesos generales del sistema.
  - Archivos generados (`*.o`, `*.ko`) → no se versionan, se crean al compilar.

---

### 1.2 Funciones principales y su propósito
- **Bash**
  - `cronjob.sh` → genera contenedores de prueba.
  - `load_modules.sh` → compila y carga módulos automáticamente.

- **Dashboard**
  - Grafana consume datos de SQLite.
  - Paneles: RAM total, RAM libre, contenedores eliminados, Top 5 RAM, Top 5 CPU.

- **Daemon en Go**
  - `InitDaemon()` → inicializa cronjob y carga módulos.
  - `ReadProcFiles()` → lee `/proc/continfo_so1_#CARNET` y `/proc/sysinfo_so1_#CARNET`.
  - `ParseJSON()` → convierte la salida en estructuras Go.
  - `AnalyzeMetrics()` → decide qué contenedores eliminar según umbrales.
  - `SaveToSQLite()` → guarda métricas en `metrics.db`.
  - `Loop()` → ejecuta cada 20 segundos.

- **Kernel (`continfo.c`, `sysinfo.c`)**
  - `init_module()` → inicializa el módulo, crea la entrada en `/proc`.
  - `cleanup_module()` → elimina la entrada en `/proc`.
  - Recorrido de `task_struct` para obtener PID, nombre, RAM, CPU.
  - Salida en formato JSON.

---

### 1.3 Dependencias externas
- Kernel headers (`linux/module.h`, `linux/proc_fs.h`, `linux/sched.h`).
- Herramientas de compilación: `make`, `gcc`.
- Go runtime.
- Docker y Docker Compose.
- Grafana y SQLite.

---

## 2. Compilación del Módulo
```bash
cd modulo-kernel
make
```

## 3. Carga y Descarga del Módulo
- Cargar

    ```bash
    sudo insmod continfo.ko
    sudo insmod sysinfo.ko
    ```

- Descargar

    ```bash
    sudo rmmod continfo
    sudo rmmod sysinfo
    ```

- Verificación

     ```bash
    lsmod | grep continfo
    dmesg | tail
    ```

## 4. Pruebas y Verificación
- Contenedores

    ```bash
    cat /proc/continfo_so1_#CARNET
    ```

- Procesos del sistema

    ```bash
    cat /proc/sysinfo_so1_#CARNET
    ```

## 5. Decisiones de Diseño y Problema

- Uso de /proc por simplicidad.

- JSON para facilitar parseo en Go.

- Separación en dos módulos para claridad.

- Problemas: permisos, precisión de CPU → solucionados con proc_create() y normalización de tiempos.

## 6. Estructura del Daemon en Go

### 6.1 Funciones principales

- InitDaemon() → inicializa cronjob y carga módulos.

- ReadProcFiles() → lee /proc/continfo_so1_#CARNET y /proc/sysinfo_so1_#CARNET.

- ParseJSON() → convierte la salida en estructuras Go.

- AnalyzeMetrics() → decide qué contenedores eliminar según umbrales.

- SaveToSQLite() → guarda métricas en metrics.db.

- Loop() → ejecuta cada 20 segundos.

### 6.2 Archivos de Automatización

- crear_contenedores.sh → Genera 10 contenedores cada 5 minuto.

- limpiar_contenedores.sh → Todos los contenedores cada 6 minuto.

- load_modules.sh → Carga los módulos de kernel.

- docker-compose.yml → Levanta Grafana, SQLite y el main.go.