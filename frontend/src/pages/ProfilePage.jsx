import { useAuth } from "../context/AuthContext";

export default function ProfilePage() {
  const { user } = useAuth();

  return (
    <section className="card page-panel profile-panel">
      <div style={{ display: 'flex', alignItems: 'center', gap: '1.5rem', marginBottom: '2.5rem' }}>
        <div style={{
          width: '76px', height: '76px', borderRadius: '50%',
          background: 'linear-gradient(135deg, var(--accent), var(--accent-alt))',
          display: 'flex', alignItems: 'center', justifyContent: 'center',
          fontSize: '2rem', fontWeight: 'bold', color: '#030712'
        }}>
          {user?.username?.[0]?.toUpperCase() || 'U'}
        </div>
        <div>
          <h2 style={{ marginBottom: '0.2rem' }}>My Profile</h2>
          <p className="muted" style={{ margin: 0 }}>Manage your personal information and account security.</p>
        </div>
      </div>

      <div className="profile-grid">
        <div className="profile-row">
          <strong>Username</strong>
          <span>{user?.username}</span>
        </div>
        <div className="profile-row">
          <strong>Email</strong>
          <span>{user?.email}</span>
        </div>
        <div className="profile-row">
          <strong>Workspace Role</strong>
          <span className={`badge badge-${user?.role}`}>{user?.role}</span>
        </div>
      </div>
    </section>
  );
}
