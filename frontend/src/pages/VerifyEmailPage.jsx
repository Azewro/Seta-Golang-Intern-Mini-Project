import { useEffect, useState } from "react";
import { Link, useSearchParams } from "react-router-dom";
import { verifyEmailApi } from "../api/authApi";

export default function VerifyEmailPage() {
  const [searchParams] = useSearchParams();
  const token = searchParams.get("token") || "";

  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  useEffect(() => {
    let mounted = true;

    async function verify() {
      if (!token) {
        if (mounted) {
          setError("Missing verification token");
          setLoading(false);
        }
        return;
      }

      try {
        const response = await verifyEmailApi(token);
        if (mounted) {
          setSuccess(response?.message || "Email verified successfully");
        }
      } catch (err) {
        if (mounted) {
          setError(err.message || "Failed to verify email");
        }
      } finally {
        if (mounted) {
          setLoading(false);
        }
      }
    }

    verify();
    return () => {
      mounted = false;
    };
  }, [token]);

  return (
    <section className="auth-layout">
      <aside className="auth-hero">
        <p className="auth-kicker">Email Security</p>
        <h2>Verify Account</h2>
        <p>
          We validate your token before granting login access to keep account
          authentication secure.
        </p>
      </aside>

      <section className="card auth-card">
        <h2>Email Verification</h2>
        <p className="muted page-lead">
          We are confirming your verification link.
        </p>
        {loading && <div className="state-card">Verifying your account...</div>}
        {!loading && success && <div className="state-card">{success}</div>}
        {!loading && error && <p className="error">{error}</p>}
        <p className="muted auth-switch">
          Continue to <Link to="/login">Login</Link>
        </p>
      </section>
    </section>
  );
}
