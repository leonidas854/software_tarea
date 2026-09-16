package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIntegration(t *testing.T) {
	mux := SetupRouter()
	server := httptest.NewServer(mux)
	defer server.Close()

	var sessionCookie *http.Cookie

	t.Run("Login Exitoso", func(t *testing.T) {
		creds := map[string]string{"username": "admin", "password": "1234"}
		body, _ := json.Marshal(creds)
		resp, err := http.Post(server.URL+"/api/login", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Error peticion login: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Esperado 200 OK, obtenido %d", resp.StatusCode)
		}

		cookies := resp.Cookies()
		for _, c := range cookies {
			if c.Name == "session_id" {
				sessionCookie = c
				break
			}
		}

		if sessionCookie == nil || sessionCookie.Value == "" {
			t.Fatalf("Esperado cookie de sesion, no se recibio o esta vacia")
		}
	})

	t.Run("Acceso Protegido Sin Cookie", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/api/cart")
		if err != nil {
			t.Fatalf("Error peticion cart: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Esperado 401 Unauthorized, obtenido %d", resp.StatusCode)
		}
	})

	t.Run("Agregar Producto y Ver Carrito", func(t *testing.T) {
		// AGREGAR_PRODUCTO: send only product_id and quantity (price comes from catalog)
		addReq := map[string]interface{}{"product_id": "TAZA-001", "quantity": 2}
		body, _ := json.Marshal(addReq)
		
		req, _ := http.NewRequest(http.MethodPost, server.URL+"/api/cart", bytes.NewBuffer(body))
		req.AddCookie(sessionCookie)
		
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Error agregando producto: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Esperado 200 OK, obtenido %d", resp.StatusCode)
		}

		// VER_CARRITO
		reqGet, _ := http.NewRequest(http.MethodGet, server.URL+"/api/cart", nil)
		reqGet.AddCookie(sessionCookie)
		
		respGet, err := client.Do(reqGet)
		if err != nil {
			t.Fatalf("Error viendo carrito: %v", err)
		}
		defer respGet.Body.Close()

		if respGet.StatusCode != http.StatusOK {
			t.Errorf("Esperado 200 OK al ver carrito, obtenido %d", respGet.StatusCode)
		}

		var result map[string]interface{}
		json.NewDecoder(respGet.Body).Decode(&result)
		
		// TAZA-001 costs $24.00, quantity 2 => subtotal = 48.00
		subtotal, ok := result["subtotal"].(float64)
		expectedSubtotal := 24.00 * 2
		if !ok || subtotal != expectedSubtotal {
			t.Errorf("Esperado subtotal %.2f, obtenido %v", expectedSubtotal, result["subtotal"])
		}
	})
	
	t.Run("Logout", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, server.URL+"/api/logout", nil)
		req.AddCookie(sessionCookie)
		
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Error logout: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Esperado 200 OK, obtenido %d", resp.StatusCode)
		}
		
		// Verificamos que la cookie tiene MaxAge negativo o está vacía
		var found bool
		for _, c := range resp.Cookies() {
			if c.Name == "session_id" && c.Value == "" {
				found = true
			}
		}
		if !found {
			t.Errorf("Esperado que se elimine la cookie de sesion")
		}
	})
}
