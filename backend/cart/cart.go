package cart

import (
	"errors"
	"sync"
)

// Product represents a product in the catalog
type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Image       string  `json:"image"`
	Category    string  `json:"category"`
	Subtitle    string  `json:"subtitle"` // Added for the design (e.g. "Arena Natural • 350ml")
}

// CartItem represents a product added to a user's cart
type CartItem struct {
	ProductID string  `json:"product_id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
	Image     string  `json:"image"`
	Subtitle  string  `json:"subtitle"`
}

// AddToCartRequest is the payload the client sends
type AddToCartRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

// Catalog for "La Tiendita"
var Catalog = []Product{
	{
		ID:          "TAZA-001", 
		Name:        "Taza Artesanal de Gres", 
		Description: "Modelada en torno manual por maestras alfareras en los Valles Centrales de Oaxaca. Cada pieza se cuece a 1,220 °C.", 
		Price:       24.00, 
		Image:       "https://lh3.googleusercontent.com/aida-public/AB6AXuBPDrrKB2jFj0aLTSjxUhThClvTF08OUrCnavMB0u81UbGi14gy_nFGUVN5OTHd97miczUg9H44wGboi8cJpFOmQa7UtTyVIJxBDqbC5CgldK0zqqghEba2-6OjU9AqfbBiSEtXTY8WGcNJTMBbdNSmLwOxLoT7c7-Yw7KTfjCZPjExn0d80pCWzEATDLurTi-YKfevwO_8qKdGwcKtyfIL9ywK3IQTryKsw4722-pi2njoJRosXBjJ", 
		Category:    "Gres Natural",
		Subtitle:    "Arena Natural • 350ml",
	},
	{
		ID:          "VELA-001", 
		Name:        "Vela Botánica de Soja", 
		Description: "Vela aromática de cera de soja con pabilo de madera. Notas de cedro y vainilla para un ambiente de calma.", 
		Price:       18.50, 
		Image:       "https://lh3.googleusercontent.com/aida-public/AB6AXuCaTiqOiL-VyTazm0Y1t2wqshcQj8jYvBXHnpVeyuJdwR3eqAbgWlCyAgpBA_GdneYBC2brxLWldtGLUXv0c7ZrO_nRvmtHOjwB6dMo3yH5fTFcJx_rla8nQYTCsoexQWBscuzdglPgLbchr6TLOmR_W2p-WCg-eVz3AACL6Y95JstrK3K_QkEwabgR-vSYx4VbIRskfMme-Q6ASg1cdSwjlUiT_qV6OX_jLDdUsfYPWUvXzBOVKzdT", 
		Category:    "Cera Soja",
		Subtitle:    "Cedro & Vainilla • 220g",
	},
	{
		ID:          "LINO-001", 
		Name:        "Camisa Lino Lavado", 
		Description: "Camisa unisex de lino lavado, textura táctil orgánica, botones de cuerno natural, estilo nórdico minimalista.", 
		Price:       65.00, 
		Image:       "https://lh3.googleusercontent.com/aida-public/AB6AXuD5isRFcLE-iXTThg75Upa9SKQ8-Bir3J5mhJ-nDjVxJj2frtPlZQUvdP9Nrd2reLKSy2bn1tsQN9T2Se3QQvEAkZjdWg8WXAXzJrKLQttBds9rZKJCO4ZTvQgG4XjB5EXa1BG4A5vXHcCOTJLGRDaP0raYhWrTLXJhhtVr04bFhnjT7xmluEM9WvIkZLINHLKl2_7fycIR1DMNPv9L24QSxCfOtfTpQljZDZmlDGjf8FGXVwZFj2kH", 
		Category:    "100% Lino",
		Subtitle:    "Arena • Talla M",
	},
	{
		ID:          "PRENSA-001", 
		Name:        "Prensa Francesa Cobre", 
		Description: "Elegante cafetera de prensa francesa en cobre y vidrio resistente al calor. Ideal para infusiones y café de especialidad.", 
		Price:       42.00, 
		Image:       "https://lh3.googleusercontent.com/aida-public/AB6AXuBXUcFN0xFiTnGVBONRpy_YXsd1Tmjh-evTN22-oUJZY-JyZ1tylOjfBjJyGLZdVxdPycKyC3JRtsL3MyMnCtUmRqrXxvDesJ5JjF72LPYVD23VH4CJs3m1qslNh4iSgMl1ahyE9RgmPagH4rm7fmM0lNmMOZnbIRUBeRwT74B9kL93SVsRfM3XGwQuRxV2T08PM4b6XgjMS01RJ-ihRPOMLIJnvuDYIHeqs-KNGdtqg4rekPtNqODT", 
		Category:    "Infusión & Café",
		Subtitle:    "Cobre • 800ml",
	},
}

// catalogMap for fast lookup by ID
var catalogMap map[string]Product

func init() {
	catalogMap = make(map[string]Product)
	for _, p := range Catalog {
		catalogMap[p.ID] = p
	}
}

// GetProductByID looks up a product in the catalog
func GetProductByID(id string) (Product, bool) {
	p, ok := catalogMap[id]
	return p, ok
}

// GestorCarrito manages the shopping cart using OOP principles
type GestorCarrito struct {
	mu    sync.RWMutex
	items map[string]CartItem
}

// NewGestorCarrito creates a new shopping cart instance
func NewGestorCarrito() *GestorCarrito {
	return &GestorCarrito{
		items: make(map[string]CartItem),
	}
}

// AgregarProducto adds a product to the cart
func (gc *GestorCarrito) AgregarProducto(req AddToCartRequest) error {
	if req.Quantity <= 0 {
		return errors.New("la cantidad debe ser un número entero mayor a 0")
	}

	product, exists := GetProductByID(req.ProductID)
	if !exists {
		return errors.New("producto no encontrado en el catálogo")
	}

	gc.mu.Lock()
	defer gc.mu.Unlock()

	if existing, ok := gc.items[req.ProductID]; ok {
		existing.Quantity += req.Quantity
		gc.items[req.ProductID] = existing
	} else {
		gc.items[req.ProductID] = CartItem{
			ProductID: product.ID,
			Name:      product.Name,
			Price:     product.Price,
			Quantity:  req.Quantity,
			Image:     product.Image,
			Subtitle:  product.Subtitle,
		}
	}

	return nil
}

// EliminarProducto removes a product from the cart entirely
func (gc *GestorCarrito) EliminarProducto(productID string) {
	gc.mu.Lock()
	defer gc.mu.Unlock()
	delete(gc.items, productID)
}

// UpdateProductQuantity updates the exact quantity of a product in the cart
func (gc *GestorCarrito) UpdateProductQuantity(productID string, qty int) error {
	if qty <= 0 {
		gc.EliminarProducto(productID)
		return nil
	}
	gc.mu.Lock()
	defer gc.mu.Unlock()

	if existing, ok := gc.items[productID]; ok {
		existing.Quantity = qty
		gc.items[productID] = existing
		return nil
	}
	return errors.New("producto no encontrado en el carrito")
}

// VerCarrito returns all items in the cart and the calculated subtotal
func (gc *GestorCarrito) VerCarrito() ([]CartItem, float64) {
	gc.mu.RLock()
	defer gc.mu.RUnlock()

	var items []CartItem
	var subtotal float64

	for _, item := range gc.items {
		items = append(items, item)
		subtotal += item.Price * float64(item.Quantity)
	}

	return items, subtotal
}

// Global store: carts associated with session IDs
var (
	carts = make(map[string]*GestorCarrito)
	mu    sync.RWMutex
)

// GetCartForSession returns the cart for a session, creating it if needed
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
