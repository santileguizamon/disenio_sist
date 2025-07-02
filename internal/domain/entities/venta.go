package entities

import (
	"time"
)

// Venta representa una venta en el sistema
type Venta struct {
	SucursalID    uint           `json:"sucursal_id"`
	ClienteID     uint           `json:"cliente_id"`
	FechaVenta    time.Time      `json:"fecha_venta"`
	Total         float64        `json:"total"`
	Subtotal      float64        `json:"subtotal"`
	Impuestos     float64        `json:"impuestos"`
	Descuento     float64        `json:"descuento"`
	Estado        string         `json:"estado"`
	MetodoPago    string         `json:"metodo_pago"`
	DetallesVenta []DetalleVenta `json:"detalles_venta"`
}

// CalcularTotal calcula el total de la venta
func (v *Venta) CalcularTotal() {
	v.Subtotal = 0
	for _, detalle := range v.DetallesVenta {
		v.Subtotal += detalle.Total
	}

	v.Impuestos = v.Subtotal * 0.21 // 21% IVA
	v.Total = v.Subtotal + v.Impuestos - v.Descuento
}

// AgregarDetalle agrega un detalle a la venta
func (v *Venta) AgregarDetalle(detalle DetalleVenta) {
	v.DetallesVenta = append(v.DetallesVenta, detalle)
	v.CalcularTotal()
}

// EsValida verifica si la venta tiene los datos mínimos requeridos
func (v *Venta) EsValida() bool {
	return v.SucursalID > 0 && !v.FechaVenta.IsZero() && len(v.DetallesVenta) > 0
}

// ActualizarEstado actualiza el estado de la venta
func (v *Venta) ActualizarEstado(nuevoEstado string) {
	v.Estado = nuevoEstado
}

// DetalleVenta representa un detalle de venta
type DetalleVenta struct {
	VentaID        uint     `json:"venta_id"`
	ProductoID     uint     `json:"producto_id"`
	Producto       Producto `json:"producto"`
	Cantidad       int      `json:"cantidad"`
	PrecioUnitario float64  `json:"precio_unitario"`
	Total          float64  `json:"total"`
	Descuento      float64  `json:"descuento"`
}

// CalcularTotal calcula el total del detalle de venta
func (d *DetalleVenta) CalcularTotal() {
	d.Total = (d.PrecioUnitario * float64(d.Cantidad)) - d.Descuento
}

// EsValido verifica si el detalle de venta tiene los datos mínimos requeridos
func (d *DetalleVenta) EsValido() bool {
	return d.VentaID > 0 && d.ProductoID > 0 && d.Cantidad > 0 && d.PrecioUnitario >= 0
}
