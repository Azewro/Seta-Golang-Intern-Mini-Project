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
    return <p>Loading users...</p>;
  }

  return (
    <section className="card">
      <h2>Users (Manager only)</h2>
      {error && <p className="error">{error}</p>}
      {!error && (
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
      )}
    </section>
  );
}

