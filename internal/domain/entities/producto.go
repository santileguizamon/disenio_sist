package entities

import (
	"time"
)

// Producto representa un producto en el sistema
type Producto struct {
	SKU                 string    `json:"sku"`
	Nombre              string    `json:"nombre"`
	Descripcion         string    `json:"descripcion"`
	Categoria           string    `json:"categoria"`
	Fabricante          string    `json:"fabricante"`
	Precio              float64   `json:"precio"`
	PrecioOferta        float64   `json:"precio_oferta"`
	StockMinimo         int       `json:"stock_minimo"`
	StockActual         int       `json:"stock_actual"`
	Estado              string    `json:"estado"`
	FechaCreacion       time.Time `json:"fecha_creacion"`
	UltimaActualizacion time.Time `json:"ultima_actualizacion"`
}

// EsValido verifica si el producto tiene los datos mínimos requeridos
func (p *Producto) EsValido() bool {
	return p.SKU != "" && p.Nombre != "" && p.Precio >= 0
}

// EnriquecerDesdeAPI enriquece los datos del producto desde una API externa
func (p *Producto) EnriquecerDesdeAPI(datosAPI map[string]interface{}) {
	if descripcion, ok := datosAPI["descripcion"].(string); ok {
		p.Descripcion = descripcion
	}
	if categoria, ok := datosAPI["categoria"].(string); ok {
		p.Categoria = categoria
	}
	if fabricante, ok := datosAPI["fabricante"].(string); ok {
		p.Fabricante = fabricante
	}
	p.UltimaActualizacion = time.Now()
}

// ActualizarStock actualiza el stock del producto
func (p *Producto) ActualizarStock(nuevoStock int) {
	p.StockActual = nuevoStock
	p.UltimaActualizacion = time.Now()
}

// TieneStockSuficiente verifica si el producto tiene stock suficiente
func (p *Producto) TieneStockSuficiente(cantidad int) bool {
	return p.StockActual >= cantidad
}

// CalcularPrecioFinal calcula el precio final considerando ofertas
func (p *Producto) CalcularPrecioFinal() float64 {
	if p.PrecioOferta > 0 && p.PrecioOferta < p.Precio {
		return p.PrecioOferta
	}
	return p.Precio
}
