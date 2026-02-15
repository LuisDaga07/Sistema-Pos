# Guía de uso - POS Restaurante

## Cómo funciona el software

1. **Registro** → Creas tu restaurante y usuario admin
2. **Login** → Inicias sesión con tu email y contraseña
3. **Categorías** → Crea categorías para organizar productos (ej: Bebidas, Hamburguesas)
4. **Productos** → Agrega productos con nombre, precio y categoría
5. **Ventas** → Registra ventas seleccionando productos y completando el pago

## Orden recomendado para empezar

### Paso 1: Categorías (opcional pero útil)

- Menú superior → **Categorías**
- Clic en **+ Nueva categoría**
- Ejemplos: "Bebidas", "Combos", "Postres"

### Paso 2: Productos

- Menú → **Productos**
- Clic en **+ Nuevo producto**
- Completa: Nombre, Precio, Categoría (si creaste alguna)
- Clic en **Guardar**

### Paso 3: Registrar una venta

- Menú → **Nueva Venta**
- Haz clic en los productos para agregarlos al carrito
- Ajusta cantidades con + y -
- Clic en **Completar venta**
- Descarga la factura en PDF si lo necesitas

## Si la pantalla se queda en blanco

1. **Backend debe estar corriendo** (en otra terminal):
   ```powershell
   cd backend
   go run ./cmd/api
   ```
   Debe mostrar: `Server starting on :8081`

2. **PostgreSQL debe estar activo** (Docker o local):
   ```powershell
   docker-compose up -d
   ```

3. **Revisa la consola del navegador** (F12 → pestaña Consola):
   - Errores de red = backend apagado o puerto equivocado
   - Error 401 = sesión expirada, vuelve a iniciar sesión

4. **Proxy del frontend**: El frontend espera el backend en `http://localhost:8081`. Si usas otro puerto, crea `frontend/.env` con:
   ```
   VITE_API_URL=http://localhost:TU_PUERTO/api/v1
   ```
