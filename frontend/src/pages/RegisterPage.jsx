import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { registerApi } from "../api/authApi";
import { useToast } from "../context/ToastContext";
import { validateEmail, validatePassword, validateUsername } from "../utils/validators";

export default function RegisterPage() {
  const navigate = useNavigate();
  const { pushToast } = useToast();
  const [form, setForm] = useState({
    username: "",
    email: "",
    password: "",
  });
  const [touched, setTouched] = useState({
    username: false,
    email: false,
    password: false,
  });
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");
  const [submitting, setSubmitting] = useState(false);

  const usernameError = validateUsername(form.username);
  const emailError = validateEmail(form.email);
  const passwordError = validatePassword(form.password);
  const hasFieldErrors = Boolean(usernameError || emailError || passwordError);

  const onChange = (key, value) => {
    setForm((prev) => ({ ...prev, [key]: value }));
  };

  const onSubmit = async (event) => {
    event.preventDefault();
    setTouched({ username: true, email: true, password: true });
    if (hasFieldErrors) {
      setError("Please fix form errors before submitting");
      return;
    }

    setError("");
    setSuccess("");
    setSubmitting(true);

    try {
      await registerApi(form);
      setSuccess("Registration successful. Please verify your email within 5 minutes. Redirecting to login...");
      pushToast("Registration successful. Check your email for the verification link.", "success");
      setTimeout(() => navigate("/login"), 1000);
    } catch (err) {
      setError(err.message);
      pushToast(err.message, "error");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <section className="auth-layout">
      <aside className="auth-hero">
        <p className="auth-kicker">Member Onboarding</p>
        <h2>Create Account</h2>
        <p>
          Register as a member account, then verify your email link to activate secure login access.
        </p>
      </aside>

      <section className="card auth-card">
        <h2>Register</h2>
        <p className="muted page-lead">Your role is assigned as member at sign-up.</p>
        <form className="form-grid" onSubmit={onSubmit}>
          <input
            type="text"
            placeholder="Username"
            value={form.username}
            onChange={(e) => onChange("username", e.target.value)}
            onBlur={() => setTouched((prev) => ({ ...prev, username: true }))}
            className={touched.username && usernameError ? "input-error" : ""}
            required
          />
          {touched.username && usernameError && <p className="field-error">{usernameError}</p>}
          <input
            type="email"
            placeholder="Email"
            value={form.email}
            onChange={(e) => onChange("email", e.target.value)}
            onBlur={() => setTouched((prev) => ({ ...prev, email: true }))}
            className={touched.email && emailError ? "input-error" : ""}
            required
          />
          {touched.email && emailError && <p className="field-error">{emailError}</p>}
          <input
            type="password"
            placeholder="Password (min 6 chars)"
            value={form.password}
            onChange={(e) => onChange("password", e.target.value)}
            onBlur={() => setTouched((prev) => ({ ...prev, password: true }))}
            className={touched.password && passwordError ? "input-error" : ""}
            minLength={6}
            required
          />
          {touched.password && passwordError && <p className="field-error">{passwordError}</p>}
          <button className="primary" type="submit" disabled={submitting || hasFieldErrors}>
            {submitting ? "Creating..." : "Register"}
          </button>
        </form>
        {error && <p className="error">{error}</p>}
        {success && <div className="state-card">{success}</div>}
        <p className="muted auth-switch">
          Already have account? <Link to="/login">Login here</Link>
        </p>
      </section>
    </section>
  );
}
