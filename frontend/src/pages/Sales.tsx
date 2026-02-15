import { useEffect, useState } from 'react';
import { productsApi, categoriesApi, salesApi, api } from '../services/api';
import type { Product, Category } from '../types';
import styles from './Sales.module.css';

interface CartItem {
  product: Product;
  quantity: number;
  toppings: { name: string; price: number; quantity: number }[];
}

export default function Sales() {
  const [products, setProducts] = useState<Product[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [cart, setCart] = useState<CartItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState('');
  const [lastSaleId, setLastSaleId] = useState<string | null>(null);
  const [filterCat, setFilterCat] = useState('');

  const load = async () => {
    try {
      setLoading(true);
      setError('');
      const [prodsRes, catsRes] = await Promise.all([
        productsApi.list({ active: 'true' }).then((r) => r.data),
        categoriesApi.list().then((r) => r.data),
      ]);
      setProducts(Array.isArray(prodsRes) ? prodsRes : []);
      setCategories(Array.isArray(catsRes) ? catsRes : []);
    } catch (e) {
      const msg = e && typeof e === 'object' && 'response' in e
        ? (e as { response?: { data?: { error?: string } } }).response?.data?.error
        : null;
      setError(msg || 'Error al cargar. ¿El backend está corriendo en http://localhost:8081?');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    load();
  }, []);

  const filteredProducts = filterCat
    ? products.filter((p) => p.category_id === filterCat)
    : products;

  const addToCart = (product: Product) => {
    const existing = cart.find((c) => c.product.id === product.id && c.toppings.length === 0);
    if (existing) {
      setCart(cart.map((c) =>
        c === existing ? { ...c, quantity: c.quantity + 1 } : c
      ));
    } else {
      setCart([...cart, { product, quantity: 1, toppings: [] }]);
    }
  };

  const updateQty = (index: number, delta: number) => {
    const item = cart[index];
    const newQty = Math.max(0, item.quantity + delta);
    if (newQty === 0) {
      setCart(cart.filter((_, i) => i !== index));
    } else {
      setCart(cart.map((c, i) => (i === index ? { ...c, quantity: newQty } : c)));
    }
  };

  const removeFromCart = (index: number) => {
    setCart(cart.filter((_, i) => i !== index));
  };

  const total = cart.reduce(
    (sum, item) =>
      sum + item.product.price * item.quantity +
      item.toppings.reduce((t, tp) => t + tp.price * tp.quantity, 0),
    0
  );

  const handleCompleteSale = async () => {
    if (cart.length === 0) {
      setError('Agrega productos a la venta');
      return;
    }
    if (total <= 0) {
      setError('El total debe ser mayor a 0');
      return;
    }

    setError('');
    setSaving(true);
    try {
      const items = cart.map((item) => ({
        product_id: item.product.id,
        quantity: item.quantity,
        toppings: item.toppings.map((t) => ({
          name: t.name,
          price: t.price,
          quantity: t.quantity,
        })),
      }));
      const payments = [{ method: 'cash', amount: total }];
      const { data } = await salesApi.create({ items, payments });
      setLastSaleId(data.id);
      setCart([]);
    } catch (e: unknown) {
      const msg = e && typeof e === 'object' && 'response' in e
        ? (e as { response?: { data?: { error?: string } } }).response?.data?.error
        : 'Error al registrar venta';
      setError(msg || 'Error');
    } finally {
      setSaving(false);
    }
  };

  const downloadPdf = async () => {
    if (!lastSaleId) return;
    try {
      const token = localStorage.getItem('token');
      const res = await api.get(`/sales/${lastSaleId}/pdf`, {
        responseType: 'blob',
        headers: token ? { Authorization: `Bearer ${token}` } : {},
      });
      const url = window.URL.createObjectURL(new Blob([res.data]));
      const a = document.createElement('a');
      a.href = url;
      a.download = `factura-${lastSaleId}.pdf`;
      a.click();
      window.URL.revokeObjectURL(url);
    } catch (e) {
      setError('Error al descargar PDF');
    }
  };

  return (
    <div className={styles.page}>
      <h1>Nueva Venta</h1>
      {error && <div className={styles.error}>{error}</div>}
      {lastSaleId && (
        <div className={styles.success}>
          Venta registrada. <button onClick={downloadPdf} type="button">Descargar factura PDF</button>
          <button onClick={() => setLastSaleId(null)} type="button">Nueva venta</button>
        </div>
      )}

      <div className={styles.layout}>
        <div className={styles.products}>
          <div className={styles.filters}>
            <select value={filterCat} onChange={(e) => setFilterCat(e.target.value)}>
              <option value="">Todas las categorías</option>
              {categories.map((c) => (
                <option key={c.id} value={c.id}>{c.name}</option>
              ))}
            </select>
          </div>
          {loading ? (
            <p>Cargando productos...</p>
          ) : (
            <div className={styles.productGrid}>
              {filteredProducts.map((p) => (
                <button
                  key={p.id}
                  type="button"
                  className={styles.productBtn}
                  onClick={() => addToCart(p)}
                >
                  <span className={styles.productName}>{p.name}</span>
                  <span className={styles.productPrice}>${Number(p.price).toFixed(2)}</span>
                </button>
              ))}
            </div>
          )}
        </div>

        <div className={styles.cart}>
          <h2>Carrito</h2>
          {cart.length === 0 ? (
            <p className={styles.empty}>Vacío</p>
          ) : (
            <>
              <ul className={styles.cartList}>
                {cart.map((item, i) => (
                  <li key={i} className={styles.cartItem}>
                    <div>
                      <strong>{item.product.name}</strong> x {item.quantity} = $
                      {(item.product.price * item.quantity).toFixed(2)}
                    </div>
                    <div className={styles.cartActions}>
                      <button onClick={() => updateQty(i, -1)} type="button">-</button>
                      <span>{item.quantity}</span>
                      <button onClick={() => updateQty(i, 1)} type="button">+</button>
                      <button onClick={() => removeFromCart(i)} type="button" className={styles.remove}>✕</button>
                    </div>
                  </li>
                ))}
              </ul>
              <div className={styles.total}>
                <strong>Total: ${total.toFixed(2)}</strong>
              </div>
              <button
                onClick={handleCompleteSale}
                disabled={saving || total <= 0}
                className={styles.completeBtn}
                type="button"
              >
                {saving ? 'Procesando...' : 'Completar venta'}
              </button>
            </>
          )}
        </div>
      </div>
    </div>
  );
}
