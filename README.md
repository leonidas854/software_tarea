# Justificación de la Práctica 1 - Segundo Parcial

Este proyecto implementa todos los requerimientos solicitados en la Práctica 1 del Segundo Parcial, dividiéndose en un **Frontend** moderno con SvelteKit (que engloba HTML, CSS y JS) y un **Backend** robusto desarrollado en Go.

A continuación se detalla cómo se cumple cada punto especificado:

## 1. Módulo Web de Autenticación y Control de Sesión

### a) Frontend (Estándares Web)
Se construyó una interfaz interactiva y estilizada utilizando estándares web modernos (HTML, CSS mediante TailwindCSS, y JavaScript a través de Svelte):
- **Estructura y Presentación (HTML/CSS):** El archivo [`frontend/src/routes/+page.svelte`](file:///home/leonidas/deveploment/software_tarea/frontend/src/routes/+page.svelte) contiene el diseño del formulario de inicio de sesión (`<form>`, `<input>`, estilos CSS).
- **Interactividad y Validación (JavaScript/Svelte):** En el mismo archivo, la función `handleLogin` valida mediante JavaScript del lado del cliente que los campos no estén vacíos antes de enviar la petición:
  ```javascript
  if (!username.trim() || !password.trim()) {
      error = 'Por favor ingresa usuario y contraseña.';
      return;
  }
  ```

### b) Backend (Scripting / Processing)
El servidor procesa la petición HTTP de inicio de sesión, verifica las credenciales y gestiona la sesión del usuario:
- En [`backend/main.go`](file:///home/leonidas/deveploment/software_tarea/backend/main.go), la función `loginHandler` procesa las peticiones `POST` a `/api/login`. 
- Si las credenciales son correctas (usando `auth.Authenticate`), se genera un identificador único (ID de Sesión).
- Este **ID de Sesión** se devuelve al cliente y se almacena utilizando una cookie segura (`http.SetCookie`), quedando registrada en la memoria del servidor.

### c) Manejo de Estado
Para proteger los recursos y permitir/restringir el acceso al contenido dinámico:
- Se implementó el middleware `authMiddleware` en [`backend/main.go`](file:///home/leonidas/deveploment/software_tarea/backend/main.go). 
- Este middleware intercepta las peticiones a rutas protegidas (como `/api/products` y `/api/cart`), extrae el ID de sesión de la cookie y valida su existencia y vigencia utilizando `auth.IsValidSession(cookie.Value)`. Si no es válido, restringe el acceso retornando un error `401 Unauthorized`.

---

## 2. Carrito de Compras Web

### a) Interfaz Cliente (HTML/CSS/JS)
El carrito de compras permite interactuar con los productos de forma dinámica:
- En [`frontend/src/routes/cart/+page.svelte`](file:///home/leonidas/deveploment/software_tarea/frontend/src/routes/cart/+page.svelte), se renderiza la lista de productos y los controles para modificar las cantidades.
- **Validación del Cliente:** La función `updateQuantity(productId, qty)` valida en JavaScript que la cantidad no sea menor a 1 antes de enviar la actualización al servidor. Si es menor a 1, dispara la lógica de eliminación (`removeItem`) para evitar enviar valores inválidos o negativos al backend.

### b) Procesamiento Server-Side (Programación Orientada a Objetos)
El backend procesa la lógica del carrito mediante un diseño orientado a objetos:
- Se implementó la clase/estructura `GestorCarrito` en [`backend/cart/cart.go`](file:///home/leonidas/deveploment/software_tarea/backend/cart/cart.go). 
- Contiene los métodos requeridos para operar sobre el carrito, tales como `AgregarProducto` (procesa la solicitud POST) y `VerCarrito` (procesa la solicitud GET), validando también en el backend que la cantidad sea un número entero mayor a 0:
  ```go
  if req.Quantity <= 0 {
      return errors.New("la cantidad debe ser un número entero mayor a 0")
  }
  ```

### c) Manejo de Sesiones (Estado del Carrito)
El carrito es temporal y está asociado a la sesión de un cliente en específico:
- El archivo [`backend/cart/cart.go`](file:///home/leonidas/deveploment/software_tarea/backend/cart/cart.go) cuenta con un almacén global en memoria (`var carts = make(map[string]*GestorCarrito)`) que asocia de forma única el ID de Sesión de cada cliente con su instancia particular del carrito.
- La función `GetCartForSession(sessionID)` asegura que, al recibir peticiones a través de múltiples llamadas HTTP, el servidor recupere el carrito correcto y mantenga el estado de los productos seleccionados para ese usuario.

### d) Respuesta Dinámica
El sistema genera dinámicamente la vista del carrito basada en la información persistida en la sesión:
- El método `VerCarrito` del backend devuelve la lista acumulada de ítems y calcula el `subtotal` exacto.
- Al responder a las peticiones (GET, POST, PUT, DELETE), el `cartHandler` en [`backend/main.go`](file:///home/leonidas/deveploment/software_tarea/backend/main.go) devuelve un JSON con el carrito actualizado.
- SvelteKit ([`frontend/src/routes/cart/+page.svelte`](file:///home/leonidas/deveploment/software_tarea/frontend/src/routes/cart/+page.svelte)) toma esta respuesta y genera de forma completamente dinámica la página web en HTML/CSS, mostrando los productos en el carrito (`{#each cartItems as item}`), calculando descuentos, envío aplicable, y mostrando el costo total en tiempo real.
