import { useAuth } from "../context/AuthContext";

export default function ProfilePage() {
  const { user, loading } = useAuth();

  if (loading) {
    return (
      <section className="card state-card">
        <h2>My Profile</h2>
        <p className="muted">Loading profile...</p>
      </section>
    );
  }

  if (!user) {
    return (
      <section className="card state-card">
        <h2>My Profile</h2>
        <p className="muted">No profile data available. Please log in again.</p>
      </section>
    );
  }

  return (
    <section className="card">
      <h2>My Profile</h2>
      <p><strong>User ID:</strong> {user?.userId}</p>
      <p><strong>Username:</strong> {user?.username}</p>
      <p><strong>Email:</strong> {user?.email}</p>
      <p><strong>Role:</strong> {user?.role}</p>
    </section>
  );
}

