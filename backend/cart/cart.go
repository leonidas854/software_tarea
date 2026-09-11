package cart

import (
	"errors"
	"sync"
)

// Product represents an item in the cart
type Product struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

// GestorCarrito manages the shopping cart using OOP principles
type GestorCarrito struct {
	mu       sync.RWMutex
	products map[string]Product
}

// NewGestorCarrito creates a new shopping cart
func NewGestorCarrito() *GestorCarrito {
	return &GestorCarrito{
		products: make(map[string]Product),
	}
}

// AgregarProducto adds a product to the cart or increments its quantity
func (gc *GestorCarrito) AgregarProducto(p Product) error {
	if p.Quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}

	gc.mu.Lock()
	defer gc.mu.Unlock()

	if existingProduct, ok := gc.products[p.ID]; ok {
		existingProduct.Quantity += p.Quantity
		gc.products[p.ID] = existingProduct
	} else {
		gc.products[p.ID] = p
	}

	return nil
}

// VerCarrito returns all products in the cart and the subtotal
func (gc *GestorCarrito) VerCarrito() ([]Product, float64) {
	gc.mu.RLock()
	defer gc.mu.RUnlock()

	var items []Product
	var subtotal float64

	for _, p := range gc.products {
		items = append(items, p)
		subtotal += p.Price * float64(p.Quantity)
	}

	return items, subtotal
}

// Global store for carts associated with session IDs
var (
	carts = make(map[string]*GestorCarrito)
	mu    sync.RWMutex
)

// GetCartForSession returns the cart for a session, creating it if it doesn't exist
func GetCartForSession(sessionID string) *GestorCarrito {
	mu.Lock()
	defer mu.Unlock()

	if c, exists := carts[sessionID]; exists {
		return c
	}
	newCart := NewGestorCarrito()
	carts[sessionID] = newCart
	return newCart
}
