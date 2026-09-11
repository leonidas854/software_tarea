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
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
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
		// Pass session_id in context or headers if needed, but we can just read it in the handler
		next(w, r)
	}
}

func cartHandler(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session_id") // Already validated by middleware
	sessionID := cookie.Value
	userCart := cart.GetCartForSession(sessionID)

	if r.Method == http.MethodGet {
		items, subtotal := userCart.VerCarrito()
		response := map[string]interface{}{
			"items":    items,
			"subtotal": subtotal,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	if r.Method == http.MethodPost {
		var p cart.Product
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		if err := userCart.AgregarProducto(p); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Product added to cart"})
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out"})
}

func main() {
	mux := http.NewServeMux()
	
	mux.HandleFunc("/api/login", enableCORS(loginHandler))
	mux.HandleFunc("/api/logout", enableCORS(logoutHandler))
	mux.HandleFunc("/api/cart", enableCORS(authMiddleware(cartHandler)))

	log.Println("Backend server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
