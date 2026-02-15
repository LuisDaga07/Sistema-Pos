import { createContext, useContext, useState, useEffect } from 'react';
import type { ReactNode } from 'react';
import { authApi } from '../services/api';

interface User {
  id: string;
  email: string;
  role: string;
}

interface Restaurant {
  id: string;
  name: string;
  email: string;
}

interface AuthContextType {
  user: User | null;
  restaurant: Restaurant | null;
  token: string | null;
  login: (email: string, password: string) => Promise<void>;
  register: (data: RegisterData) => Promise<void>;
  logout: () => void;
  isAuthenticated: boolean;
}

interface RegisterData {
  restaurant_name: string;
  email: string;
  password: string;
  phone?: string;
  address?: string;
  tax_id?: string;
}

const AuthContext = createContext<AuthContextType | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [restaurant, setRestaurant] = useState<Restaurant | null>(null);
  const [token, setToken] = useState<string | null>(localStorage.getItem('token'));

  useEffect(() => {
    const stored = localStorage.getItem('user');
    const rest = localStorage.getItem('restaurant');
    if (stored && rest) {
      try {
        setUser(JSON.parse(stored));
        setRestaurant(JSON.parse(rest));
      } catch {
        localStorage.removeItem('user');
        localStorage.removeItem('restaurant');
        localStorage.removeItem('token');
      }
    }
  }, []);

  const login = async (email: string, password: string) => {
    const { data } = await authApi.login({ email, password });
    localStorage.setItem('token', data.token);
    localStorage.setItem('user', JSON.stringify(data.user));
    localStorage.setItem('restaurant', JSON.stringify(data.restaurant));
    setToken(data.token);
    setUser(data.user);
    setRestaurant(data.restaurant);
  };

  const register = async (data: RegisterData) => {
    const { data: res } = await authApi.register(data);
    localStorage.setItem('token', res.token);
    localStorage.setItem('user', JSON.stringify(res.user));
    localStorage.setItem('restaurant', JSON.stringify(res.restaurant));
    setToken(res.token);
    setUser(res.user);
    setRestaurant(res.restaurant);
  };

  const logout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    localStorage.removeItem('restaurant');
    setToken(null);
    setUser(null);
    setRestaurant(null);
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        restaurant,
        token,
        login,
        register,
        logout,
        isAuthenticated: !!token,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error('useAuth must be used within AuthProvider');
  return ctx;
}
