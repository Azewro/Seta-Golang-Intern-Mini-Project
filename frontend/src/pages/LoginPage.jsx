import { useState } from "react";
import { Link, useLocation, useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import { useToast } from "../context/ToastContext";
import { validateEmail, validatePassword } from "../utils/validators";

export default function LoginPage() {
  const navigate = useNavigate();
  const location = useLocation();
  const { login } = useAuth();
  const { pushToast } = useToast();

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [touched, setTouched] = useState({ email: false, password: false });
  const [error, setError] = useState("");
  const [submitting, setSubmitting] = useState(false);

  const from = location.state?.from?.pathname;
  const emailError = validateEmail(email);
  const passwordError = validatePassword(password);
  const hasFieldErrors = Boolean(emailError || passwordError);

  const onSubmit = async (event) => {
    event.preventDefault();
    setTouched({ email: true, password: true });
    if (hasFieldErrors) {
      setError("Please fix form errors before submitting");
      return;
    }

    setError("");
    setSubmitting(true);

    try {
      const loginResult = await login({ email, password });
      pushToast("Login successful", "success");
      const defaultPath = loginResult?.user?.role === "manager" ? "/users" : "/profile";
      const nextPath = from && from !== "/login" ? from : defaultPath;
      navigate(nextPath, { replace: true });
    } catch (err) {
      setError(err.message);
      pushToast(err.message, "error");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <section className="card">
      <h2>Login</h2>
      <form className="form-grid" onSubmit={onSubmit}>
        <input
          type="email"
          placeholder="Email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          onBlur={() => setTouched((prev) => ({ ...prev, email: true }))}
          className={touched.email && emailError ? "input-error" : ""}
          required
        />
        {touched.email && emailError && <p className="field-error">{emailError}</p>}
        <input
          type="password"
          placeholder="Password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          onBlur={() => setTouched((prev) => ({ ...prev, password: true }))}
          className={touched.password && passwordError ? "input-error" : ""}
          required
        />
        {touched.password && passwordError && <p className="field-error">{passwordError}</p>}
        <button className="primary" type="submit" disabled={submitting || hasFieldErrors}>
          {submitting ? "Logging in..." : "Login"}
        </button>
      </form>
      {error && <p className="error">{error}</p>}
      <p>
        No account yet? <Link to="/register">Register here</Link>
      </p>
    </section>
  );
}




