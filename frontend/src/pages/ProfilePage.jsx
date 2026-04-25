import { useAuth } from "../context/AuthContext";

export default function ProfilePage() {
  const { user } = useAuth();

  return (
    <section className="card">
      <h2>My Profile</h2>
      <div className="profile-grid">
        <p className="profile-row">
          <strong>User ID:</strong> {user?.userId}
        </p>
        <p className="profile-row">
          <strong>Username:</strong> {user?.username}
        </p>
        <p className="profile-row">
          <strong>Email:</strong> {user?.email}
        </p>
        <p className="profile-row">
          <strong>Role:</strong> {user?.role}
        </p>
        <p className="profile-row">
          <strong>Verification:</strong> {user?.isVerified ? "Verified" : "Pending verification"}
        </p>
      </div>
    </section>
  );
}
