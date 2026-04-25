import { useState } from "react";
import { Link, useLocation, useNavigate } from "react-router-dom";
import { resendVerificationApi } from "../api/authApi";
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
  const [resending, setResending] = useState(false);

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

  const onResendVerification = async () => {
    setTouched((prev) => ({ ...prev, email: true }));
    if (emailError) {
      setError("Enter a valid email before resending verification");
      return;
    }

    setResending(true);
    try {
      const response = await resendVerificationApi({ email });
      pushToast(response?.message || "Verification email sent", "success");
    } catch (err) {
      pushToast(err.message, "error");
    } finally {
      setResending(false);
    }
  };

  const shouldShowResend = error.toLowerCase().includes("not verified");

  return (
    <section className="auth-layout">
      <aside className="auth-hero">
        <p className="auth-kicker">Identity Access</p>
        <h2>Welcome Back</h2>
        <p>
          Sign in to access your profile, verify your account status, and continue team management workflows.
        </p>
      </aside>

      <section className="card auth-card">
        <h2>Login</h2>
        <p className="muted page-lead">Use your registered email and password.</p>
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
        {shouldShowResend && (
          <button className="primary" type="button" onClick={onResendVerification} disabled={resending}>
            {resending ? "Sending..." : "Resend verification email"}
          </button>
        )}
        <p className="muted auth-switch">
          No account yet? <Link to="/register">Register here</Link>
        </p>
      </section>
    </section>
  );
}
