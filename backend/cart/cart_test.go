package cart_test

import (
	"backend/cart"
	"testing"
)

func TestGestorCarrito(t *testing.T) {
	gc := cart.NewGestorCarrito()

	t.Run("Agregar Producto Nuevo", func(t *testing.T) {
		err := gc.AgregarProducto(cart.Product{ID: "p1", Name: "Laptop", Price: 1000.0, Quantity: 1})
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		
		items, subtotal := gc.VerCarrito()
		if len(items) != 1 {
			t.Errorf("Expected 1 item, got %d", len(items))
		}
		if subtotal != 1000.0 {
			t.Errorf("Expected subtotal 1000.0, got %f", subtotal)
		}
	})

	t.Run("Agregar Producto Existente", func(t *testing.T) {
		err := gc.AgregarProducto(cart.Product{ID: "p1", Name: "Laptop", Price: 1000.0, Quantity: 2})
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		items, subtotal := gc.VerCarrito()
		if len(items) != 1 {
			t.Errorf("Expected 1 item, got %d", len(items))
		}
		if items[0].Quantity != 3 {
			t.Errorf("Expected quantity 3, got %d", items[0].Quantity)
		}
		if subtotal != 3000.0 {
			t.Errorf("Expected subtotal 3000.0, got %f", subtotal)
		}
	})

	t.Run("Cantidad Invalida", func(t *testing.T) {
		err := gc.AgregarProducto(cart.Product{ID: "p2", Name: "Mouse", Price: 50.0, Quantity: 0})
		if err == nil {
			t.Errorf("Expected error for quantity 0, got nil")
		}
	})
}

func TestSessionCarts(t *testing.T) {
	sessionID := "session-123"
	
	c1 := cart.GetCartForSession(sessionID)
	if c1 == nil {
		t.Errorf("Expected cart to be created")
	}

	c2 := cart.GetCartForSession(sessionID)
	if c1 != c2 {
		t.Errorf("Expected same cart instance for same session")
	}
}
