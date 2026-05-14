import { Navigate, Route, Routes } from "react-router-dom";
import Navbar from "./components/Navbar";
import LoginPage from "./pages/LoginPage";
import RegisterPage from "./pages/RegisterPage";
import ProfilePage from "./pages/ProfilePage";
import UsersPage from "./pages/UsersPage";
import TeamsPage from "./pages/TeamsPage";
import VerifyEmailPage from "./pages/VerifyEmailPage";
import NotFoundPage from "./pages/NotFoundPage";
import ProtectedRoute from "./routes/ProtectedRoute";
import ManagerRoute from "./routes/ManagerRoute";
import ToastContainer from "./components/ToastContainer";
import { useAuth } from "./context/AuthContext";

function HomeRedirect() {
  const { user } = useAuth();
  if (!user) {
    return <Navigate to="/login" replace />;
  }
  const destination = user.role === "manager" ? "/users" : "/profile";
  return <Navigate to={destination} replace />;
}

function GuestRoute({ children }) {
  const { user } = useAuth();
  if (user) {
    const destination = user.role === "manager" ? "/users" : "/profile";
    return <Navigate to={destination} replace />;
  }
  return children;
}

export default function App() {
  return (
    <div className="app-shell">
      <Navbar />
      <ToastContainer />
      <main className="container">
        <Routes>
          <Route path="/" element={<HomeRedirect />} />
          <Route
            path="/login"
            element={
              <GuestRoute>
                <LoginPage />
              </GuestRoute>
            }
          />
          <Route
            path="/register"
            element={
              <GuestRoute>
                <RegisterPage />
              </GuestRoute>
            }
          />
          <Route path="/verify-email" element={<VerifyEmailPage />} />

          <Route
            path="/profile"
            element={
              <ProtectedRoute>
                <ProfilePage />
              </ProtectedRoute>
            }
          />

          <Route
            path="/users"
            element={
              <ManagerRoute>
                <UsersPage />
              </ManagerRoute>
            }
          />

          <Route
            path="/teams"
            element={
              <ManagerRoute>
                <TeamsPage />
              </ManagerRoute>
            }
          />

          <Route path="*" element={<NotFoundPage />} />
        </Routes>
      </main>
    </div>
  );
}
