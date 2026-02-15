import { Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import styles from './Dashboard.module.css';

export default function Dashboard() {
  const { restaurant } = useAuth();

  return (
    <div className={styles.dashboard}>
      <h1>Bienvenido{restaurant ? `, ${restaurant.name}` : ''}</h1>
      <p>Usa el menú para gestionar tu restaurante.</p>
      <div className={styles.cards}>
        <Link to="/products" className={styles.card}>
          <h2>Productos</h2>
          <p>Gestionar menú y precios</p>
        </Link>
        <Link to="/categories" className={styles.card}>
          <h2>Categorías</h2>
          <p>Organizar productos por categoría</p>
        </Link>
        <Link to="/sales" className={styles.card}>
          <h2>Nueva Venta</h2>
          <p>Registrar una venta</p>
        </Link>
      </div>
    </div>
  );
}
