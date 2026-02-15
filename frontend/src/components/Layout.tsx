import { Outlet, Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import styles from './Layout.module.css';

export default function Layout() {
  const { user, restaurant, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <div className={styles.layout}>
      <header className={styles.header}>
        <div className={styles.brand}>
          <Link to="/">POS</Link>
          {restaurant && <span className={styles.restaurant}>{restaurant.name}</span>}
        </div>
        <nav className={styles.nav}>
          <Link to="/products">Productos</Link>
          <Link to="/categories">Categorías</Link>
          <Link to="/sales">Nueva Venta</Link>
        </nav>
        <div className={styles.user}>
          <span>{user?.email}</span>
          <button onClick={handleLogout} type="button">Salir</button>
        </div>
      </header>
      <main className={styles.main}>
        <Outlet />
      </main>
    </div>
  );
}
