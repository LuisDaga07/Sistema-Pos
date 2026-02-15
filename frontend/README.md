# POS Restaurante - Frontend

Frontend React para el sistema POS SaaS.

## Requisitos

- Node.js 18+
- Backend Go corriendo en `http://localhost:8081`

## Instalación

```bash
npm install
```

## Desarrollo

```bash
npm run dev
```

Abre [http://localhost:5173](http://localhost:5173). Las peticiones a `/api` se redirigen al backend en el puerto 8081.

## Build

```bash
npm run build
```

Los archivos se generan en `dist/`.

## Estructura

- `src/pages/` - Login, Register, Dashboard, Products, Categories, Sales
- `src/components/` - Layout, ProtectedRoute
- `src/context/` - AuthContext (JWT)
- `src/services/` - API client (axios)
