# Guía de Ejecución del Sistema

## Requisitos Previos

### Software Necesario
- **Go 1.21 o superior**
- **Git** (para clonar el repositorio)

### Verificar Instalación
```bash
# Verificar versión de Go
go version
```

## Configuración Inicial

### 1. Clonar el Repositorio
```bash
git clone <url-del-repositorio>
cd sistema-gestion-informacion
```

### 2. Configurar Variables de Entorno
```bash
# Copiar archivo de ejemplo
cp env.example .env

# Editar variables según tu configuración
nano .env
```

**Contenido mínimo del archivo .env:**
```env
# Configuración del Servidor
PORT=8080
ENVIRONMENT=development
```

## Instalación y Ejecución

### 1. Descargar Dependencias
```bash
go mod download
```

### 2. Verificar Dependencias
```bash
go mod tidy
```

### 3. Ejecutar el Sistema
```bash
go run cmd/main.go
```

### 4. Verificar que el Sistema Esté Funcionando
```bash
# Verificar endpoint de salud
curl http://localhost:8080/health

# Verificar información del sistema
curl http://localhost:8080/
```

## Endpoints Disponibles

### Endpoints Principales
- **GET /** - Información del sistema
- **GET /health** - Estado de salud del sistema
- **GET /swagger/** - Documentación Swagger

### API de Procesamiento
- **POST /api/procesar** - Procesar y depurar datos crudos
- **GET /api/datos-procesados** - Consultar datos procesados
- **GET /api/reporte** - Obtener último reporte generado

## Ejemplos de Uso

### 1. Procesar Datos Crudos
```bash
curl -X POST http://localhost:8080/api/procesar \
  -H "Content-Type: application/json" \
  -d '{
    "datos": [
      {
        "producto_id": 1,
        "nombre": "Producto A",
        "precio": 100.50,
        "sucursal_id": 1,
        "cantidad": 5
      },
      {
        "producto_id": 2,
        "nombre": "Producto B",
        "precio": 75.25,
        "sucursal_id": 1,
        "cantidad": 3
      }
    ]
  }'
```

### 2. Consultar Datos Procesados
```bash
curl http://localhost:8080/api/datos-procesados
```

### 3. Obtener Último Reporte
```bash
curl http://localhost:8080/api/reporte
```

## Verificación del Sistema

### 1. Verificar Logs
El sistema mostrará logs en la consola con emojis para facilitar la identificación:

- 🚀 **Servidor iniciando**
- ✅ **Operaciones exitosas**
- ❌ **Errores**
- 📥 **Datos recibidos**
- ⚙️ **Datos procesados**
- 📊 **Reportes generados**
- 🔔 **Eventos disparados**

### 2. Verificar Swagger
Abrir en el navegador: `http://localhost:8080/swagger/`

## Solución de Problemas

### Error de Puerto en Uso
```
❌ Error iniciando servidor: listen tcp :8080: bind: address already in use
```

**Solución:**
1. Cambiar el puerto en el archivo .env
2. O detener el proceso que está usando el puerto 8080

### Error de Dependencias
```
go: module lookup disabled by GOPROXY=off
```

**Solución:**
```bash
go env -w GOPROXY=https://proxy.golang.org,direct
go mod download
```

## Estructura de Archivos

```
sistema-gestion-informacion/
├── cmd/
│   └── main.go                 # Punto de entrada
├── internal/
│   ├── domain/entities/        # Entidades de dominio
│   │   ├── producto.go         # Entidad Producto
│   │   ├── sucursal.go         # Entidad Sucursal
│   │   └── venta.go            # Entidad Venta
│   ├── application/services/   # Servicios de aplicación
│   │   └── procesador_datos_service.go
│   ├── infrastructure/         # Infraestructura
│   │   ├── events/            # Sistema de eventos
│   │   │   └── event_bus.go   # EventBus (Singleton + Observer)
│   │   └── builders/          # Patrón Builder
│   │       └── reporte_builder.go
│   └── interfaces/            # Controladores y rutas
│       └── handlers/          # Handlers HTTP
│           └── procesamiento_handler.go
├── docs/                      # Documentación
├── go.mod                     # Dependencias Go
├── env.example                # Variables de entorno
└── README.md                  # Documentación principal
```

## Patrones de Diseño Implementados

### 1. Singleton
- **EventBus**: Una única instancia global del bus de eventos

### 2. Observer (Event-Driven Architecture)
- **Sistema de eventos**: Comunicación desacoplada entre componentes
- **Handlers de eventos**: DatosProcesadosHandler, EventLogger, EventMetrics

### 3. Builder
- **ReporteBuilder**: Construcción flexible de reportes con diferentes configuraciones

## Flujo de Procesamiento

1. **Recepción de datos**: El endpoint `POST /api/procesar` recibe datos crudos
2. **Procesamiento**: Los datos se procesan y depuran en memoria
3. **Eventos**: Se disparan eventos para notificar el procesamiento
4. **Almacenamiento**: Los datos procesados se almacenan en memoria
5. **Consulta**: Los endpoints GET permiten consultar datos y reportes

## Comandos Útiles

### Desarrollo
```bash
# Ejecutar en modo desarrollo
go run cmd/main.go

# Ejecutar con variables de entorno específicas
ENVIRONMENT=development go run cmd/main.go

# Verificar sintaxis
go vet ./...

# Formatear código
go fmt ./...
```

### Testing
```bash
# Ejecutar tests
go test ./...

# Ejecutar tests con coverage
go test -cover ./...
```

## Características del Sistema

- **Arquitectura basada en eventos**: Comunicación desacoplada entre componentes
- **Procesamiento en memoria**: Sin dependencia de base de datos
- **API REST**: Endpoints simples y claros
- **Patrones de diseño**: Singleton, Observer y Builder implementados
- **Logging estructurado**: Logs con emojis para fácil identificación
- **Documentación Swagger**: API documentada automáticamente 