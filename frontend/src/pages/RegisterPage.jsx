import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { registerApi } from "../api/authApi";
import { useToast } from "../context/ToastContext";
import { validateEmail, validatePassword, validateRole, validateUsername } from "../utils/validators";

export default function RegisterPage() {
  const navigate = useNavigate();
  const { pushToast } = useToast();
  const [form, setForm] = useState({
    username: "",
    email: "",
    password: "",
    role: "member",
  });
  const [touched, setTouched] = useState({
    username: false,
    email: false,
    password: false,
    role: false,
  });
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");
  const [submitting, setSubmitting] = useState(false);

  const usernameError = validateUsername(form.username);
  const emailError = validateEmail(form.email);
  const passwordError = validatePassword(form.password);
  const roleError = validateRole(form.role);
  const hasFieldErrors = Boolean(usernameError || emailError || passwordError || roleError);

  const onChange = (key, value) => {
    setForm((prev) => ({ ...prev, [key]: value }));
  };

  const onSubmit = async (event) => {
    event.preventDefault();
    setTouched({ username: true, email: true, password: true, role: true });
    if (hasFieldErrors) {
      setError("Please fix form errors before submitting");
      return;
    }

    setError("");
    setSuccess("");
    setSubmitting(true);

    try {
      await registerApi(form);
      setSuccess("Register successful. Redirecting to login...");
      pushToast("Register successful", "success");
      setTimeout(() => navigate("/login"), 800);
    } catch (err) {
      setError(err.message);
      pushToast(err.message, "error");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <section className="card">
      <h2>Register</h2>
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
        <select
          value={form.role}
          onChange={(e) => onChange("role", e.target.value)}
          onBlur={() => setTouched((prev) => ({ ...prev, role: true }))}
          className={touched.role && roleError ? "input-error" : ""}
          required
        >
          <option value="member">member</option>
          <option value="manager">manager</option>
        </select>
        {touched.role && roleError && <p className="field-error">{roleError}</p>}
        <button className="primary" type="submit" disabled={submitting || hasFieldErrors}>
          {submitting ? "Creating..." : "Register"}
        </button>
      </form>
      {error && <p className="error">{error}</p>}
      {success && <p>{success}</p>}
      <p>
        Already have account? <Link to="/login">Login here</Link>
      </p>
    </section>
  );
}

