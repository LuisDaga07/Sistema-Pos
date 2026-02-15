import axios from 'axios';

const API_BASE = import.meta.env.VITE_API_URL || '/api/v1';

export const api = axios.create({
  baseURL: API_BASE,
  headers: { 'Content-Type': 'application/json' },
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

api.interceptors.response.use(
  (res) => res,
  (err) => {
    if (err.response?.status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    return Promise.reject(err);
  }
);

// Auth
export const authApi = {
  register: (data: { restaurant_name: string; email: string; password: string; phone?: string; address?: string; tax_id?: string }) =>
    api.post('/auth/register', data),
  login: (data: { email: string; password: string }) => api.post('/auth/login', data),
};

// Categories
export const categoriesApi = {
  list: () => api.get('/categories'),
  create: (data: { name: string; description?: string; sort_order?: number }) =>
    api.post('/categories', data),
};

// Products
export const productsApi = {
  list: (params?: { category_id?: string; active?: string }) =>
    api.get('/products', { params }),
  get: (id: string) => api.get(`/products/${id}`),
  create: (data: { category_id?: string; name: string; description?: string; price: number; image_url?: string; active?: boolean }) =>
    api.post('/products', data),
  update: (id: string, data: Partial<{ category_id: string; name: string; description: string; price: number; image_url: string; active: boolean }>) =>
    api.put(`/products/${id}`, data),
  delete: (id: string) => api.delete(`/products/${id}`),
};

// Sales
export const salesApi = {
  create: (data: {
    items: Array<{
      product_id: string;
      quantity: number;
      notes?: string;
      toppings?: Array<{ name: string; price: number; quantity: number }>;
    }>;
    payments: Array<{ method: string; amount: number; reference?: string }>;
  }) => api.post('/sales', data),
  get: (id: string) => api.get(`/sales/${id}`),
};
