import { useEffect, useState } from "react";
import { listUsersApi } from "../api/userApi";
import { useAuth } from "../context/AuthContext";

export default function UsersPage() {
  const { token } = useAuth();
  const [users, setUsers] = useState([]);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(true);
  const [reloadKey, setReloadKey] = useState(0);

  useEffect(() => {
    let active = true;

    async function loadUsers() {
      setLoading(true);
      setError("");
      try {
        const data = await listUsersApi(token, 1, 50);
        if (active) {
          setUsers(data.data || []);
        }
      } catch (err) {
        if (active) {
          setError(err.message);
        }
      } finally {
        if (active) {
          setLoading(false);
        }
      }
    }

    loadUsers();
    return () => {
      active = false;
    };
  }, [token, reloadKey]);

  if (loading) {
    return (
      <section className="card state-card">
        <h2>Users (Manager only)</h2>
        <p className="muted">Loading users...</p>
      </section>
    );
  }

  if (error) {
    return (
      <section className="card state-card">
        <h2>Users (Manager only)</h2>
        <p className="error">{error}</p>
        <button type="button" className="primary" onClick={() => setReloadKey((prev) => prev + 1)}>
          Retry
        </button>
      </section>
    );
  }

  if (users.length === 0) {
    return (
      <section className="card state-card">
        <h2>Users (Manager only)</h2>
        <p className="muted">No users found.</p>
      </section>
    );
  }

  return (
    <section className="card">
      <h2>Users (Manager only)</h2>
      <table className="table">
        <thead>
          <tr>
            <th>ID</th>
            <th>Username</th>
            <th>Email</th>
            <th>Role</th>
          </tr>
        </thead>
        <tbody>
          {users.map((u) => (
            <tr key={u.userId}>
              <td>{u.userId}</td>
              <td>{u.username}</td>
              <td>{u.email}</td>
              <td>{u.role}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </section>
  );
}

