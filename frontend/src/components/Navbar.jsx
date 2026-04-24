import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import { useToast } from "../context/ToastContext";

export default function Navbar() {
  const { user, isAuthenticated, logout } = useAuth();
  const { pushToast } = useToast();
  const navigate = useNavigate();

  const onLogout = async () => {
    try {
      await logout();
      pushToast("Logged out", "info");
      navigate("/login", { replace: true });
    } catch (err) {
      pushToast(err.message || "Logout failed", "error");
    }
  };

  return (
    <header className="navbar">
      <div className="nav-content">
        <strong>Seta Stage 1</strong>
        <nav className="nav-links">
          {!isAuthenticated && <Link className="nav-link" to="/login">Login</Link>}
          {!isAuthenticated && <Link className="nav-link" to="/register">Register</Link>}
          {isAuthenticated && <Link className="nav-link" to="/profile">My Profile</Link>}
          {isAuthenticated && user?.role === "manager" && <Link className="nav-link" to="/users">Users</Link>}
          {isAuthenticated && (
            <button type="button" className="nav-btn" onClick={onLogout}>
              Logout
            </button>
          )}
        </nav>
      </div>
    </header>
  );
}


