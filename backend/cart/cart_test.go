package cart

import (
	"testing"
)

func TestAgregarProducto(t *testing.T) {
	gc := NewGestorCarrito()

	err := gc.AgregarProducto(AddToCartRequest{ProductID: "TAZA-001", Quantity: 2})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	items, subtotal := gc.VerCarrito()
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].Quantity != 2 {
		t.Fatalf("expected quantity 2, got %d", items[0].Quantity)
	}
	expectedSubtotal := 24.00 * 2
	if subtotal != expectedSubtotal {
		t.Fatalf("expected subtotal %.2f, got %.2f", expectedSubtotal, subtotal)
	}
}

func TestAgregarProductoInvalidQuantity(t *testing.T) {
	gc := NewGestorCarrito()
	err := gc.AgregarProducto(AddToCartRequest{ProductID: "TAZA-001", Quantity: 0})
	if err == nil {
		t.Fatal("expected error for quantity <= 0")
	}
}

func TestAgregarProductoNotInCatalog(t *testing.T) {
	gc := NewGestorCarrito()
	err := gc.AgregarProducto(AddToCartRequest{ProductID: "FAKE-999", Quantity: 1})
	if err == nil {
		t.Fatal("expected error for non-existent product")
	}
}

func TestAgregarProductoAccumulates(t *testing.T) {
	gc := NewGestorCarrito()
	gc.AgregarProducto(AddToCartRequest{ProductID: "VELA-001", Quantity: 1})
	gc.AgregarProducto(AddToCartRequest{ProductID: "VELA-001", Quantity: 3})

	items, _ := gc.VerCarrito()
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].Quantity != 4 {
		t.Fatalf("expected quantity 4, got %d", items[0].Quantity)
	}
}

func TestEliminarProducto(t *testing.T) {
	gc := NewGestorCarrito()
	gc.AgregarProducto(AddToCartRequest{ProductID: "PRENSA-001", Quantity: 1})
	gc.EliminarProducto("PRENSA-001")

	items, subtotal := gc.VerCarrito()
	if len(items) != 0 {
		t.Fatalf("expected 0 items, got %d", len(items))
	}
	if subtotal != 0 {
		t.Fatalf("expected subtotal 0, got %.2f", subtotal)
	}
}

func TestVerCarritoEmpty(t *testing.T) {
	gc := NewGestorCarrito()
	items, subtotal := gc.VerCarrito()
	if items != nil && len(items) != 0 {
		t.Fatalf("expected empty cart, got %d items", len(items))
	}
	if subtotal != 0 {
		t.Fatalf("expected subtotal 0, got %.2f", subtotal)
	}
}

func TestGetCartForSession(t *testing.T) {
	c1 := GetCartForSession("session-test-1")
	c2 := GetCartForSession("session-test-1")
	if c1 != c2 {
		t.Fatal("expected same cart instance for same session")
	}

	c3 := GetCartForSession("session-test-2")
	if c1 == c3 {
		t.Fatal("expected different cart for different session")
	}
}

func TestGetProductByID(t *testing.T) {
	p, ok := GetProductByID("TAZA-001")
	if !ok {
		t.Fatal("expected to find TAZA-001")
	}
	if p.Name != "Taza Artesanal de Gres" {
		t.Fatalf("wrong product name: %s", p.Name)
	}

	_, ok = GetProductByID("NONEXISTENT")
	if ok {
		t.Fatal("expected not to find NONEXISTENT")
	}
}
