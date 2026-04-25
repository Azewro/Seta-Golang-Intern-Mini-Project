import { NavLink, useNavigate } from "react-router-dom";
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

  const navClassName = ({ isActive }) => (isActive ? "nav-link nav-link-active" : "nav-link");

  return (
    <header className="navbar">
      <div className="nav-content">
        <div className="nav-brand">
          <strong>Seta Auth Portal</strong>
          <span>Stage 1 • Identity and Team Foundation</span>
        </div>

        <nav className="nav-links">
          {!isAuthenticated && (
            <>
              <NavLink className={navClassName} to="/login">
                Login
              </NavLink>
              <NavLink className={navClassName} to="/register">
                Register
              </NavLink>
            </>
          )}

          {isAuthenticated && (
            <>
              <NavLink className={navClassName} to="/profile">
                My Profile
              </NavLink>
              {user?.role === "manager" && (
                <NavLink className={navClassName} to="/users">
                  Users
                </NavLink>
              )}
              <button type="button" className="nav-btn" onClick={onLogout}>
                Logout
              </button>
            </>
          )}
        </nav>
      </div>
    </header>
  );
}
