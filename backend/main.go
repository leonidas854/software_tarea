package main

import (
	"backend/auth"
	"backend/cart"
	"encoding/json"
	"log"
	"net/http"
)

// LoginCredentials matches the JSON payload for login
type LoginCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173") // SvelteKit dev server
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds LoginCredentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	sessionID, err := auth.Authenticate(creds.Username, creds.Password)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true, // Security best practice
		SameSite: http.SameSiteLaxMode,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged in successfully"})
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil || !auth.IsValidSession(cookie.Value) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

// productsHandler serves the product catalog as JSON
func productsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart.Catalog)
}

// cartHandler handles AGREGAR_PRODUCTO (POST) and VER_CARRITO (GET)
func cartHandler(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session_id") // Already validated by middleware
	sessionID := cookie.Value
	userCart := cart.GetCartForSession(sessionID)

	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		// VER_CARRITO
		items, subtotal := userCart.VerCarrito()
		if items == nil {
			items = []cart.CartItem{} // Return empty array instead of null
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"items":    items,
			"subtotal": subtotal,
		})

	case http.MethodPost:
		// AGREGAR_PRODUCTO
		var req cart.AddToCartRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		if err := userCart.AgregarProducto(req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		// Return updated cart
		items, subtotal := userCart.VerCarrito()
		if items == nil {
			items = []cart.CartItem{}
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":  "Producto agregado al carrito",
			"items":    items,
			"subtotal": subtotal,
		})

	case http.MethodPut:
		// UPDATE QUANTITY
		var req struct {
			ProductID string `json:"product_id"`
			Quantity  int    `json:"quantity"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		if err := userCart.UpdateProductQuantity(req.ProductID, req.Quantity); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		items, subtotal := userCart.VerCarrito()
		if items == nil {
			items = []cart.CartItem{}
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":  "Cantidad actualizada",
			"items":    items,
			"subtotal": subtotal,
		})

	case http.MethodDelete:
		// Remove product from cart
		var req struct {
			ProductID string `json:"product_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		userCart.EliminarProducto(req.ProductID)

		items, subtotal := userCart.VerCarrito()
		if items == nil {
			items = []cart.CartItem{}
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":  "Producto eliminado del carrito",
			"items":    items,
			"subtotal": subtotal,
		})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Invalidate session on the server
	if cookie, err := r.Cookie("session_id"); err == nil {
		auth.InvalidateSession(cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out"})
}

func SetupRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/login", enableCORS(loginHandler))
	mux.HandleFunc("/api/logout", enableCORS(logoutHandler))
	mux.HandleFunc("/api/products", enableCORS(authMiddleware(productsHandler)))
	mux.HandleFunc("/api/cart", enableCORS(authMiddleware(cartHandler)))

	return mux
}

func main() {
	mux := SetupRouter()
	log.Println("Backend server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
