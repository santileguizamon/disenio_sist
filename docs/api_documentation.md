# Documentación de la API REST - Sistema de Procesamiento de Datos

## Descripción General

Esta API REST proporciona endpoints para procesar, depurar y consultar datos de productos, ventas y sucursales. El sistema utiliza una arquitectura basada en eventos y procesamiento en memoria, sin dependencia de base de datos.

## Base URL

```
http://localhost:8080/api
```

## Endpoints

### Procesamiento de Datos

#### Procesar Datos Crudos
- **POST** `/procesar`
- **Descripción**: Recibe datos crudos, los procesa y depura, y dispara eventos de notificación
- **Body**:
```json
{
  "datos": [
    {
      "producto_id": 1,
      "nombre": "Producto A",
      "precio": 100.50,
      "sucursal_id": 1,
      "cantidad": 5,
      "fecha_venta": "2024-01-15T10:30:00Z"
    },
    {
      "producto_id": 2,
      "nombre": "Producto B",
      "precio": 75.25,
      "sucursal_id": 1,
      "cantidad": 3,
      "fecha_venta": "2024-01-15T11:00:00Z"
    }
  ]
}
```
- **Respuesta Exitosa** (200):
```json
{
  "mensaje": "Datos procesados exitosamente",
  "registros_procesados": 2,
  "registros_depurados": 2,
  "eventos_disparados": 3
}
```

#### Consultar Datos Procesados
- **GET** `/datos-procesados`
- **Descripción**: Obtiene los datos procesados y depurados almacenados en memoria
- **Respuesta Exitosa** (200):
```json
{
  "productos": [
    {
      "ID": 1,
      "SKU": "PROD-001",
      "Nombre": "Producto A",
      "Descripcion": "Descripción del producto A",
      "Categoria": "general",
      "Fabricante": "Fabricante XYZ",
      "Precio": 100.50,
      "PrecioOferta": 90.45,
      "StockMinimo": 10,
      "StockActual": 45
    }
  ],
  "ventas": [
    {
      "ID": 1,
      "SucursalID": 1,
      "FechaVenta": "2024-01-15T10:30:00Z",
      "MetodoPago": "efectivo",
      "Total": 502.50,
      "DetallesVenta": [
        {
          "ProductoID": 1,
          "Cantidad": 5,
          "PrecioUnitario": 100.50,
          "Descuento": 0.00
        }
      ]
    }
  ],
  "sucursales": [
    {
      "ID": 1,
      "Nombre": "Sucursal Centro",
      "Direccion": "Av. Principal 123",
      "Telefono": "+1234567890",
      "Email": "centro@empresa.com",
      "Ciudad": "Ciudad Principal",
      "TipoSistema": "api",
      "ApiEndpoint": "https://api.sucursal.com",
      "ApiKey": "key123",
      "ApiSecret": "secret123"
    }
  ]
}
```

#### Obtener Último Reporte
- **GET** `/reporte`
- **Descripción**: Obtiene el último reporte generado usando el patrón Builder
- **Respuesta Exitosa** (200):
```json
{
  "ID": 1,
  "Tipo": "ventas",
  "GeneradoEn": "2024-01-15T12:00:00Z",
  "Estado": "completado",
  "FechaInicio": "2024-01-01T00:00:00Z",
  "FechaFin": "2024-01-15T23:59:59Z",
  "Datos": {
    "total_ventas": 502.50,
    "cantidad_transacciones": 1,
    "producto_mas_vendido": "Producto A",
    "sucursal_principal": "Sucursal Centro"
  },
  "Filtros": ["fecha_inicio", "fecha_fin"],
  "Resumen": "Reporte de ventas del período 2024-01-01 al 2024-01-15"
}
```

## Arquitectura de Eventos

El sistema implementa una arquitectura basada en eventos usando el patrón Observer:

### Eventos Disponibles
- `datos.procesados`: Se dispara cuando se completan el procesamiento y depuración de datos
- `reporte.generado`: Se dispara cuando se genera un nuevo reporte

### Handlers de Eventos
- **DatosProcesadosHandler**: Maneja la notificación de datos procesados
- **EventLogger**: Registra todos los eventos en logs
- **EventMetrics**: Recopila métricas de eventos

### Ejemplo de Flujo de Eventos
1. Se reciben datos crudos en `POST /api/procesar`
2. Se procesan y depuran los datos
3. Se dispara el evento `datos.procesados`
4. Los handlers suscritos procesan el evento
5. Se genera un reporte automáticamente
6. Se dispara el evento `reporte.generado`

## Códigos de Error

### Errores HTTP Comunes

- **400 Bad Request**: Datos de entrada inválidos o mal formateados
- **404 Not Found**: Endpoint no encontrado
- **422 Unprocessable Entity**: Datos válidos pero no procesables
- **500 Internal Server Error**: Error interno del servidor

### Formato de Respuesta de Error

```json
{
  "error": "Tipo de error",
  "message": "Descripción detallada del error",
  "timestamp": "2024-01-15T10:30:00Z",
  "path": "/api/procesar",
  "status": 400
}
```

## Ejemplos de Uso

### Procesar Datos de Ventas
```bash
curl -X POST http://localhost:8080/api/procesar \
  -H "Content-Type: application/json" \
  -d '{
    "datos": [
      {
        "producto_id": 1,
        "nombre": "Laptop Gaming",
        "precio": 1200.00,
        "sucursal_id": 1,
        "cantidad": 2,
        "fecha_venta": "2024-01-15T14:30:00Z"
      },
      {
        "producto_id": 2,
        "nombre": "Mouse Inalámbrico",
        "precio": 45.50,
        "sucursal_id": 1,
        "cantidad": 5,
        "fecha_venta": "2024-01-15T15:00:00Z"
      }
    ]
  }'
```

### Consultar Datos Procesados
```bash
curl -X GET http://localhost:8080/api/datos-procesados
```

### Obtener Reporte Actual
```bash
curl -X GET http://localhost:8080/api/reporte
```

## Características del Sistema

### Procesamiento en Memoria
- Todos los datos se procesan y almacenan en memoria
- No hay dependencia de base de datos
- Los datos persisten durante la ejecución del servidor

### Patrones de Diseño Implementados
- **Singleton**: EventBus con una única instancia global
- **Observer**: Sistema de eventos con handlers suscritos
- **Builder**: Construcción flexible de reportes

### Logging y Monitoreo
- Logs estructurados con emojis para fácil identificación
- Métricas de eventos procesados
- Trazabilidad completa del flujo de datos

## Notas Importantes

1. **Datos en Memoria**: Los datos se pierden al reiniciar el servidor
2. **Sin Autenticación**: La API no requiere autenticación
3. **Sin Paginación**: Los endpoints retornan todos los datos disponibles
4. **Eventos Síncronos**: Los eventos se procesan de forma síncrona
5. **Validación**: Todos los datos de entrada son validados antes del procesamiento
6. **Logs**: Todas las operaciones son registradas para auditoría 