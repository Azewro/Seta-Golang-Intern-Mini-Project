import { useEffect, useState } from "react";
import { listUsersApi } from "../api/userApi";
import { useAuth } from "../context/AuthContext";

export default function UsersPage() {
  const { token } = useAuth();
  const [users, setUsers] = useState([]);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let active = true;

    async function loadUsers() {
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
  }, [token]);

  if (loading) {
    return (
      <section className="card page-panel users-panel" style={{ display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
        <div style={{ color: 'var(--text-muted)', fontSize: '1.1rem' }}>
          <div className="spinner"></div> Loading workspace users...
        </div>
      </section>
    );
  }

  return (
    <section className="card page-panel users-panel">
      <div style={{ marginBottom: '2.5rem' }}>
        <h2 style={{ marginBottom: '0.5rem' }}>Users Directory</h2>
        <p className="muted" style={{ margin: 0 }}>Manager override: Viewing all workspace accounts and access levels.</p>
      </div>
      
      {error && <p className="error">{error}</p>}
      {!error && (
        <div className="table-wrap">
          <table className="table">
            <thead>
              <tr>
                <th>Member ID</th>
                <th>Username</th>
                <th>Contact Email</th>
                <th>Access Level</th>
              </tr>
            </thead>
            <tbody>
              {users.map((u) => (
                <tr key={u.userId}>
                  <td style={{ opacity: 0.6, fontFamily: 'monospace' }}>#{u.userId}</td>
                  <td style={{ fontWeight: '500', color: '#fff' }}>{u.username}</td>
                  <td className="muted">{u.email}</td>
                  <td>
                    <span className={`badge badge-${u.role}`}>
                      {u.role}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </section>
  );
}

