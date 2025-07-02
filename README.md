# Sistema de Gestión de Información de Clientes Potenciales

## 🎯 Descripción

Sistema centralizado de integración y análisis de datos que automatiza la recolección, normalización y almacenamiento de datos de clientes, ventas y stock para una empresa de marketing digital. Implementado en **Go** con arquitectura dirigida por eventos y patrones de diseño SOLID.

## ✨ Características Principales

- **🔄 Arquitectura Event-Driven**: Sistema de eventos asíncrono para comunicación entre componentes
- **🏗️ Patrones de Diseño**: Singleton, Builder, Observer, Repository, Strategy
- **📊 API RESTful**: Endpoints documentados con Swagger
- **🗄️ Persistencia**: Base de datos MySQL con GORM
- **📈 Pipeline de Procesamiento**: Recolección → Procesamiento → Persistencia → Reportes
- **🔧 Principios SOLID**: Código limpio y mantenible

## 🚀 Inicio Rápido

### Prerrequisitos
- Go 1.21+
- MySQL 8.0+
- Git

### Instalación

1. **Clonar el repositorio**
```bash
git clone <url-del-repositorio>
cd sistema-gestion-informacion
```

2. **Configurar base de datos**
```sql
CREATE DATABASE sistema_gestion_informacion;
```

3. **Configurar variables de entorno**
```bash
cp env.example .env
# Editar .env con tus credenciales de MySQL
```

4. **Ejecutar el sistema**
```bash
go mod download
go run cmd/main.go
```

5. **Verificar funcionamiento**
```bash
curl http://localhost:8080/health
```

## 📚 Documentación

- **[Guía de Ejecución](docs/guia_ejecucion.md)** - Instrucciones detalladas de instalación y uso
- **[Documentación de la API](docs/api_documentation.md)** - Endpoints disponibles y ejemplos
- **[Patrones de Diseño](docs/patrones_diseno.md)** - Explicación de los patrones implementados

## 🏛️ Arquitectura

### Estructura del Proyecto
```
sistema-gestion-informacion/
├── cmd/main.go                    # 🚀 Punto de entrada
├── internal/
│   ├── domain/entities/           # 📦 Entidades de dominio
│   │   ├── cliente_potencial.go
│   │   ├── producto.go
│   │   ├── venta.go
│   │   └── sucursal.go
│   ├── infrastructure/            # 🔧 Infraestructura
│   │   ├── database/singleton.go  # Patrón Singleton
│   │   ├── events/event_bus.go    # Sistema de eventos
│   │   └── builders/              # Patrón Builder
│   └── interfaces/                # 🌐 Interfaces
│       ├── handlers/              # Controladores HTTP
│       └── routes/                # Configuración de rutas
├── docs/                          # 📖 Documentación
├── go.mod                         # 📦 Dependencias
└── env.example                    # ⚙️ Configuración
```

### Patrones de Diseño Implementados

#### 1. Singleton
- **Base de datos**: Una única conexión compartida
- **Bus de eventos**: Una única instancia global

#### 2. Builder
- **Reportes**: Construcción flexible de reportes
- **Configuraciones**: Construcción de configuraciones complejas

#### 3. Observer (Event-Driven)
- **Sistema de eventos**: Comunicación desacoplada entre componentes
- **Pipeline de procesamiento**: Activación automática de pasos

## 🔌 API Endpoints

### Clientes Potenciales
- `GET /api/clientes` - Obtener todos los clientes
- `POST /api/clientes` - Crear nuevo cliente
- `GET /api/clientes/{id}` - Obtener cliente por ID

### Procesamiento
- `POST /api/procesar` - Ejecutar pipeline completo

### Sistema
- `GET /` - Información del sistema
- `GET /health` - Estado de salud
- `GET /swagger/` - Documentación Swagger

## 💡 Ejemplos de Uso

### Crear Cliente Potencial
```bash
curl -X POST http://localhost:8080/api/clientes \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Juan Pérez",
    "email": "juan.perez@email.com",
    "telefono": "+1234567890",
    "fuente": "website",
    "interes": "producto_a",
    "sucursal_id": 1
  }'
```

### Ejecutar Procesamiento
```bash
curl -X POST http://localhost:8080/api/procesar
```

## 🔄 Pipeline de Procesamiento

El sistema implementa un pipeline completo de procesamiento de datos:

1. **📥 Recolección**: Simulación de recolección desde fuentes externas
2. **⚙️ Procesamiento**: Validación y normalización de datos
3. **💾 Persistencia**: Almacenamiento en base de datos
4. **📊 Reportes**: Generación automática de reportes

Cada paso dispara eventos que pueden ser escuchados por otros componentes del sistema.

## 🛠️ Tecnologías Utilizadas

- **Go 1.21** - Lenguaje principal
- **GORM** - ORM para MySQL
- **MySQL** - Base de datos
- **Swagger** - Documentación de API
- **HTTP** - Servidor web nativo de Go

## 📊 Estado del Proyecto

### ✅ Implementado
- [x] Arquitectura base con patrones de diseño
- [x] Sistema de eventos (Event-Driven)
- [x] API REST básica
- [x] Persistencia con GORM
- [x] Documentación Swagger
- [x] Pipeline de procesamiento simulado
- [x] Patrones Singleton y Builder

### 🚧 En Desarrollo
- [ ] Repositorios concretos
- [ ] Autenticación JWT
- [ ] Generación de reportes Excel
- [ ] Tests unitarios
- [ ] Tareas programadas

### 📋 Pendiente
- [ ] Monitoreo y métricas
- [ ] Rate limiting
- [ ] Logging estructurado
- [ ] Backup automático

## 🤝 Contribución

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## 📝 Licencia

Este proyecto está bajo la Licencia MIT. Ver el archivo `LICENSE` para más detalles.

## 📞 Soporte

Para soporte técnico o consultas:
- Revisar la [documentación](docs/)
- Verificar los [logs del sistema](docs/guia_ejecucion.md#verificación-del-sistema)
- Crear un issue en el repositorio

---
