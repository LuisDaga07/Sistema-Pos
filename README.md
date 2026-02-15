# POS Restaurante - Sistema SaaS

Sistema POS para restaurantes de comidas rápidas. Backend en Go (Gin) + Frontend en React.

## Arquitectura

- **Backend:** Go + Gin + PostgreSQL
- **Frontend:** React + Vite + TypeScript
- **Autenticación:** JWT
- **Multi-tenant:** `restaurant_id` en todas las tablas

## Inicio rápido

### 1. Base de datos (Docker)

```bash
docker-compose up -d
```

PostgreSQL en puerto 5433. Las migraciones se ejecutan automáticamente.

### 2. Backend

```bash
cd backend
go run ./cmd/api
```

API en `http://localhost:8081`. Usa `backend/.env` para configuración (ver `backend/.env.example`).

### 3. Frontend

```bash
cd frontend
npm install
npm run dev
```

Frontend en `http://localhost:5173`. El proxy de Vite envía las peticiones `/api` al backend.

## Uso

1. Registrar un restaurante en `/register`
2. Iniciar sesión en `/login`
3. Crear categorías y productos
4. Registrar ventas en "Nueva Venta"
5. Descargar factura PDF tras cada venta
