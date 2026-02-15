import { useEffect, useState } from 'react';
import { categoriesApi } from '../services/api';
import type { Category } from '../types';
import styles from './Categories.module.css';

export default function Categories() {
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showModal, setShowModal] = useState(false);
  const [form, setForm] = useState({ name: '', description: '' });

  const load = async () => {
    try {
      setLoading(true);
      setError('');
      const { data } = await categoriesApi.list();
      setCategories(Array.isArray(data) ? data : []);
    } catch (e) {
      const msg = e && typeof e === 'object' && 'response' in e
        ? (e as { response?: { data?: { error?: string } } }).response?.data?.error
        : null;
      setError(msg || 'Error al cargar. ¿El backend está en http://localhost:8081?');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    load();
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    try {
      await categoriesApi.create(form);
      setForm({ name: '', description: '' });
      setShowModal(false);
      load();
    } catch (e: unknown) {
      const msg = e && typeof e === 'object' && 'response' in e
        ? (e as { response?: { data?: { error?: string } } }).response?.data?.error
        : 'Error al guardar';
      setError(msg || 'Error');
    }
  };

  return (
    <div className={styles.page}>
      <div className={styles.header}>
        <h1>Categorías</h1>
        <button onClick={() => setShowModal(true)} type="button">+ Nueva categoría</button>
      </div>
      {error && <div className={styles.error}>{error}</div>}
      {loading ? (
        <p>Cargando...</p>
      ) : (
        <div className={styles.grid}>
          {categories.map((c) => (
            <div key={c.id} className={styles.card}>
              <h3>{c.name}</h3>
              {c.description && <p>{c.description}</p>}
            </div>
          ))}
        </div>
      )}

      {showModal && (
        <div className={styles.modal} onClick={() => setShowModal(false)}>
          <div className={styles.modalContent} onClick={(e) => e.stopPropagation()}>
            <h2>Nueva categoría</h2>
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
