import { useEffect, useState } from 'react';
import { productsApi, categoriesApi } from '../services/api';
import type { Product, Category } from '../types';
import styles from './Products.module.css';

export default function Products() {
  const [products, setProducts] = useState<Product[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showModal, setShowModal] = useState(false);
  const [editing, setEditing] = useState<Product | null>(null);
  const [form, setForm] = useState({
    name: '',
    description: '',
    price: '',
    category_id: '',
    active: true,
  });

  const load = async () => {
    try {
      setLoading(true);
      setError('');
      const [prodsRes, catsRes] = await Promise.all([
        productsApi.list().then((r) => r.data),
        categoriesApi.list().then((r) => r.data),
      ]);
      setProducts(Array.isArray(prodsRes) ? prodsRes : []);
      setCategories(Array.isArray(catsRes) ? catsRes : []);
    } catch (e) {
      const msg = e && typeof e === 'object' && 'response' in e
        ? (e as { response?: { status?: number } }).response?.status === 401
          ? 'Sesión expirada. Vuelve a iniciar sesión.'
          : (e as { response?: { data?: { error?: string } } }).response?.data?.error
        : null;
      setError(msg || 'Error al cargar. ¿El backend está en http://localhost:8081?');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    load();
  }, []);

  const openCreate = () => {
    setEditing(null);
    setForm({ name: '', description: '', price: '', category_id: '', active: true });
    setShowModal(true);
  };

  const openEdit = (p: Product) => {
    setEditing(p);
    setForm({
      name: p.name,
      description: p.description || '',
      price: String(p.price),
      category_id: p.category_id || '',
      active: p.active,
    });
    setShowModal(true);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    try {
      const payload = {
        name: form.name,
        description: form.description,
        price: parseFloat(form.price) || 0,
        category_id: form.category_id || undefined,
        active: form.active,
      };
      if (editing) {
        await productsApi.update(editing.id, payload);
      } else {
        await productsApi.create(payload);
      }
      setShowModal(false);
      load();
    } catch (e: unknown) {
      const msg = e && typeof e === 'object' && 'response' in e
        ? (e as { response?: { data?: { error?: string } } }).response?.data?.error
        : 'Error al guardar';
      setError(msg || 'Error');
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm('¿Eliminar producto?')) return;
    try {
      await productsApi.delete(id);
      load();
    } catch (e) {
      setError('Error al eliminar');
    }
  };

  return (
    <div className={styles.page}>
      <div className={styles.header}>
        <h1>Productos</h1>
        <button onClick={openCreate} type="button">+ Nuevo producto</button>
      </div>
      {error && <div className={styles.error}>{error}</div>}
      {loading ? (
        <p>Cargando...</p>
      ) : (
        <div className={styles.grid}>
          {products.map((p) => (
            <div key={p.id} className={styles.card}>
              <div className={styles.cardBody}>
                <h3>{p.name}</h3>
                <p className={styles.price}>${Number(p.price).toFixed(2)}</p>
                {p.description && <p className={styles.desc}>{p.description}</p>}
              </div>
              <div className={styles.actions}>
                <button onClick={() => openEdit(p)} type="button">Editar</button>
                <button onClick={() => handleDelete(p.id)} type="button" className={styles.danger}>Eliminar</button>
              </div>
            </div>
          ))}
        </div>
      )}

      {showModal && (
        <div className={styles.modal} onClick={() => setShowModal(false)}>
          <div className={styles.modalContent} onClick={(e) => e.stopPropagation()}>
            <h2>{editing ? 'Editar producto' : 'Nuevo producto'}</h2>
            <form onSubmit={handleSubmit}>
              <input
                placeholder="Nombre"
                value={form.name}
                onChange={(e) => setForm({ ...form, name: e.target.value })}
                required
              />
              <input
                placeholder="Descripción"
                value={form.description}
                onChange={(e) => setForm({ ...form, description: e.target.value })}
              />
              <input
                type="number"
                step="0.01"
                placeholder="Precio"
                value={form.price}
                onChange={(e) => setForm({ ...form, price: e.target.value })}
                required
              />
              <select
                value={form.category_id}
                onChange={(e) => setForm({ ...form, category_id: e.target.value })}
              >
                <option value="">Sin categoría</option>
                {categories.map((c) => (
                  <option key={c.id} value={c.id}>{c.name}</option>
                ))}
              </select>
              <label>
                <input
                  type="checkbox"
                  checked={form.active}
                  onChange={(e) => setForm({ ...form, active: e.target.checked })}
                />
                Activo
              </label>
              <div className={styles.modalActions}>
                <button type="button" onClick={() => setShowModal(false)}>Cancelar</button>
                <button type="submit">Guardar</button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
