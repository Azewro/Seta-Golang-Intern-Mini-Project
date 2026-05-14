import { useEffect, useState } from "react";
import { useAuth } from "../context/AuthContext";
import { useToast } from "../context/ToastContext";
import { listMyTeamsApi, createTeamApi, addMemberApi, addManagerApi, removeManagerApi, removeMemberApi, getTeamApi } from "../api/teamApi";
import { listUsersApi, bulkGetUsersApi } from "../api/userApi";

export default function TeamsPage() {
  const { user, token } = useAuth();
  const { pushToast } = useToast();

  const [teams, setTeams] = useState([]);
  const [loading, setLoading] = useState(true);
  const [newTeamName, setNewTeamName] = useState("");
  const [selectedTeam, setSelectedTeam] = useState(null);
  const [users, setUsers] = useState([]);
  const [usersLoading, setUsersLoading] = useState(false);
  const [selectedUserId, setSelectedUserId] = useState("");
  const [userSearchTerm, setUserSearchTerm] = useState("");
  const [userRoleFilter, setUserRoleFilter] = useState("");
  const [isAddUserModalOpen, setIsAddUserModalOpen] = useState(false);

  const isManager = user?.role === "manager";

  const fetchTeams = async () => {
    try {
      setLoading(true);
      const res = await listMyTeamsApi(token);
      setTeams(res?.data || []);
    } catch (err) {
      pushToast(err.message || "Failed to load teams", "error");
    } finally {
      setLoading(false);
    }
  };

  const fetchUsers = async () => {
    try {
      setUsersLoading(true);
      const data = await listUsersApi(token, 1, 200);
      setUsers(data?.data || []);
    } catch (err) {
      pushToast(err.message || "Failed to load users", "error");
    } finally {
      setUsersLoading(false);
    }
  };

  const handleCreateTeam = async (e) => {
    e.preventDefault();
    if (!newTeamName.trim()) return;
    try {
      await createTeamApi({ teamName: newTeamName }, token);
      pushToast("Team created", "success");
      setNewTeamName("");
      fetchTeams();
    } catch (err) {
      pushToast(err.message, "error");
    }
  };

  const handleSelectTeam = async (teamId) => {
    try {
      const data = await getTeamApi(teamId, token);
      setSelectedTeam(data);
    } catch (err) {
      pushToast(err.message, "error");
    }
  };

  const handleAction = async (actionFn, ...args) => {
    try {
      await actionFn(...args);
      pushToast("Action successful", "success");
      if (selectedTeam) handleSelectTeam(selectedTeam.teamId);
      setSelectedUserId("");
    } catch (err) {
      pushToast(err.message, "error");
    }
  };

  useEffect(() => {
    if (token) fetchTeams();
  }, [token]);

  useEffect(() => {
    if (token && isManager) fetchUsers();
  }, [token, isManager]);

  useEffect(() => {
    if (!isManager && selectedTeam && token) {
      const ids = [...new Set([...(selectedTeam.managers || []), ...(selectedTeam.members || []), selectedTeam.mainManagerUserId].filter(Boolean))];
      const missingIds = ids.filter(id => !users.find(u => u.userId === id));
      if (missingIds.length > 0) {
        bulkGetUsersApi(missingIds, token)
          .then(res => {
            if (res?.data) {
              setUsers(prev => {
                const newUsers = res.data.filter(nu => !prev.find(pu => pu.userId === nu.userId));
                return [...prev, ...newUsers];
              });
            }
          })
          .catch(err => {
            console.error("Failed to load team member details", err);
          });
      }
    }
  }, [selectedTeam, isManager, token, users]);

  const userMap = new Map(users.map((u) => [u.userId, u]));
  const formatUserLabel = (id) => {
    const target = userMap.get(id);
    if (!target) return `User #${id}`;
    return `${target.username} (${target.email})`;
  };

  const isTeamManager = selectedTeam?.managers?.includes(user?.userId);
  const isMainManager = selectedTeam?.mainManagerUserId === user?.userId;
  const canManageTeam = Boolean(isManager && isTeamManager);
  const canManageManagers = Boolean(isManager && isMainManager);

  const isInTeam = (userId) =>
    selectedTeam?.managers?.includes(userId) || selectedTeam?.members?.includes(userId);

  const availableUsers = selectedTeam
    ? users.filter((u) => !isInTeam(u.userId))
    : [];

  const searchedUsers = availableUsers.filter((u) => {
    const term = userSearchTerm.toLowerCase();
    const matchesTerm = u.username.toLowerCase().includes(term) || u.email.toLowerCase().includes(term);
    const matchesRole = userRoleFilter ? u.role === userRoleFilter : true;
    return matchesTerm && matchesRole;
  });

  const handleAddMemberAction = (targetId) => {
    if (!selectedTeam) return;
    const targetUser = userMap.get(targetId);
    if (!targetUser) return;
    if (isInTeam(targetId)) {
      pushToast("User already in team", "error");
      return;
    }
    handleAction(addMemberApi, selectedTeam.teamId, targetId, token);
  };

  const handleAddManagerAction = (targetId) => {
    if (!selectedTeam) return;
    const targetUser = userMap.get(targetId);
    if (!targetUser) return;
    if (targetUser.role !== "manager") {
      pushToast("User must have manager role", "error");
      return;
    }
    if (isInTeam(targetId)) {
      pushToast("User already in team", "error");
      return;
    }
    handleAction(addManagerApi, selectedTeam.teamId, targetId, token);
  };

  if (loading) {
    return (
      <section className="page-panel">
        <div className="card">
          Loading teams...
        </div>
      </section>
    );
  }

  return (
    <section className="page-panel">
      <div className="card teams-panel" style={{ width: '100%' }}>
        <h2>My Teams</h2>

        {isManager && (
          <form className="form-group" onSubmit={handleCreateTeam} style={{ display: 'flex', gap: '12px', marginBottom: '24px', alignItems: 'center' }}>
            <input
              type="text"
              placeholder="New Team Name"
              value={newTeamName}
              onChange={(e) => setNewTeamName(e.target.value)}
              style={{ flex: 1 }}
            />
            <button type="submit" className="primary" style={{ margin: 0, whiteSpace: 'nowrap' }}>Create Team</button>
          </form>
        )}

        {teams.length === 0 ? (
          <p className="muted">You are not in any teams yet.</p>
        ) : (
          <div className="table-wrap" style={{ marginBottom: '24px' }}>
            <table className="table">
              <thead>
                <tr>
                  <th style={{ width: '80px' }}>ID</th>
                  <th>Name</th>
                  <th style={{ width: '150px' }}>Action</th>
                </tr>
              </thead>
              <tbody>
                {teams.map((t) => (
                  <tr key={t.teamId}>
                    <td className="muted">#{t.teamId}</td>
                    <td><strong>{t.teamName}</strong></td>
                    <td>
                      <button
                        className="primary"
                        style={{
                          margin: 0,
                          padding: '0.6rem 1rem',
                          backgroundColor: selectedTeam?.teamId === t.teamId ? '#475569' : '',
                          color: selectedTeam?.teamId === t.teamId ? '#fff' : ''
                        }}
                        onClick={() => selectedTeam?.teamId === t.teamId ? setSelectedTeam(null) : handleSelectTeam(t.teamId)}
                      >
                        {selectedTeam?.teamId === t.teamId ? "Close Details" : "View Details"}
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}

        {selectedTeam && (
          <div className="card" style={{ marginTop: '24px', border: '1px solid rgba(148, 163, 184, 0.15)' }}>
            <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: '12px' }}>
              <h3 style={{ margin: 0 }}>{selectedTeam.teamName} Details</h3>
              <button
                className="primary"
                style={{ margin: 0, backgroundColor: '#475569', color: '#fff' }}
                onClick={() => setSelectedTeam(null)}>
                Close
              </button>
            </div>
            <p>Main manager: {formatUserLabel(selectedTeam.mainManagerUserId)}</p>

            {isManager && (
              <div style={{ marginTop: '16px', padding: '16px', borderRadius: '16px', border: '1px solid var(--surface-border)', background: 'rgba(3, 7, 18, 0.35)' }}>
                <h4 style={{ marginTop: 0, marginBottom: '12px' }}>Team Members</h4>
                <button
                  className="primary"
                  style={{ margin: 0, backgroundColor: '#2dd4bf', color: '#030712' }}
                  onClick={() => {
                    setUserSearchTerm("");
                    setUserRoleFilter("");
                    setIsAddUserModalOpen(true);
                  }}
                  disabled={!canManageTeam || usersLoading}
                >
                  {usersLoading ? "Loading users..." : "+ Add Person to Team"}
                </button>
                {!usersLoading && availableUsers.length === 0 && (
                  <p className="muted" style={{ marginTop: '12px', marginBottom: 0 }}>All users are already in this team.</p>
                )}
              </div>
            )}

            <div style={{ display: 'grid', gap: '16px', marginTop: '16px' }}>
              <div>
                <h4 style={{ marginBottom: '8px' }}>Managers</h4>
                <div style={{ display: 'grid', gap: '8px' }}>
                  {(selectedTeam.managers || []).length === 0 && (
                    <p className="muted" style={{ margin: 0 }}>No managers yet.</p>
                  )}
                  {(selectedTeam.managers || []).map((managerId) => (
                    <div
                      key={`manager-${managerId}`}
                      style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', padding: '0.6rem 0.9rem', borderRadius: '12px', border: '1px solid var(--surface-border)', background: 'rgba(15, 23, 42, 0.45)' }}
                    >
                      <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                        <span>{formatUserLabel(managerId)}</span>
                        {managerId === selectedTeam.mainManagerUserId && (
                          <span className="badge" style={{ background: 'var(--accent)', color: '#000', fontSize: '0.65rem' }}>Main Manager</span>
                        )}
                      </div>
                      {canManageManagers && (
                        <button
                          className="primary"
                          style={{ margin: 0, backgroundColor: '#475569', color: '#fff' }}
                          onClick={() => handleAction(removeManagerApi, selectedTeam.teamId, managerId, token)}
                          disabled={managerId === selectedTeam.mainManagerUserId}
                        >
                          Remove
                        </button>
                      )}
                    </div>
                  ))}
                </div>
              </div>

              <div>
                <h4 style={{ marginBottom: '8px' }}>Members</h4>
                <div style={{ display: 'grid', gap: '8px' }}>
                  {(selectedTeam.members || []).length === 0 && (
                    <p className="muted" style={{ margin: 0 }}>No members yet.</p>
                  )}
                  {(selectedTeam.members || []).map((memberId) => (
                    <div
                      key={`member-${memberId}`}
                      style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', padding: '0.6rem 0.9rem', borderRadius: '12px', border: '1px solid var(--surface-border)', background: 'rgba(15, 23, 42, 0.45)' }}
                    >
                      <span>{formatUserLabel(memberId)}</span>
                      {canManageTeam && (
                        <button
                          className="primary"
                          style={{ margin: 0, backgroundColor: '#475569', color: '#fff' }}
                          onClick={() => handleAction(removeMemberApi, selectedTeam.teamId, memberId, token)}
                        >
                          Remove
                        </button>
                      )}
                    </div>
                  ))}
                </div>
              </div>
            </div>
          </div>
        )}
      </div>

      {/* ADD USER MODAL */}
      {isAddUserModalOpen && (
        <div style={{
          position: 'fixed', top: 0, left: 0, right: 0, bottom: 0,
          backgroundColor: 'rgba(3, 7, 18, 0.8)', backdropFilter: 'blur(8px)',
          zIndex: 9999, display: 'flex', alignItems: 'center', justifyContent: 'center',
          padding: '1.5rem'
        }}>
          <div className="card" style={{
            width: '100%', maxWidth: '900px', maxHeight: '90vh',
            display: 'flex', flexDirection: 'column', overflow: 'hidden', padding: '0'
          }}>
            <div style={{ padding: '1.5rem 2rem', borderBottom: '1px solid var(--surface-border)', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <h2 style={{ margin: 0 }}>Add Person to Team</h2>
              <button
                onClick={() => setIsAddUserModalOpen(false)}
                style={{ background: 'transparent', border: 'none', color: 'var(--text-muted)', fontSize: '1.5rem', cursor: 'pointer' }}
              >
                ✕
              </button>
            </div>

            <div style={{ padding: '1.5rem 2rem', borderBottom: '1px solid var(--surface-border)', display: 'flex', gap: '1rem', background: 'rgba(15, 23, 42, 0.4)' }}>
              <input
                type="text"
                placeholder="Search username or email..."
                value={userSearchTerm}
                onChange={(e) => setUserSearchTerm(e.target.value)}
                style={{ flex: 1 }}
              />
              <select
                value={userRoleFilter}
                onChange={(e) => setUserRoleFilter(e.target.value)}
                style={{ width: '200px' }}
              >
                <option value="">All Roles</option>
                <option value="manager">Manager</option>
                <option value="member">Member</option>
              </select>
            </div>

            <div style={{ padding: '1.5rem 2rem', overflowY: 'auto', flex: 1 }}>
              {searchedUsers.length === 0 ? (
                <p className="muted text-center" style={{ padding: '2rem 0' }}>No users found matching your criteria.</p>
              ) : (
                <div className="table-wrap">
                  <table className="table">
                    <thead>
                      <tr>
                        <th style={{ width: '80px' }}>ID</th>
                        <th>User Info</th>
                        <th>Role</th>
                        <th style={{ textAlign: 'right' }}>Actions</th>
                      </tr>
                    </thead>
                    <tbody>
                      {searchedUsers.map((u) => (
                        <tr key={u.userId}>
                          <td className="muted" style={{ fontFamily: 'monospace' }}>#{u.userId}</td>
                          <td>
                            <div style={{ fontWeight: 600 }}>{u.username}</div>
                            <div className="muted" style={{ fontSize: '0.85rem' }}>{u.email}</div>
                          </td>
                          <td><span className={`badge badge-${u.role}`}>{u.role}</span></td>
                          <td style={{ textAlign: 'right' }}>
                            <div style={{ display: 'flex', gap: '8px', justifyContent: 'flex-end' }}>
                              {canManageTeam && (
                                <button
                                  className="primary"
                                  style={{ margin: 0, padding: '0.4rem 0.8rem', fontSize: '0.85rem' }}
                                  onClick={() => handleAddMemberAction(u.userId)}
                                >
                                  Add as Member
                                </button>
                              )}
                              {canManageManagers && u.role === 'manager' && (
                                <button
                                  className="primary"
                                  style={{ margin: 0, padding: '0.4rem 0.8rem', fontSize: '0.85rem', backgroundColor: '#0ea5e9' }}
                                  onClick={() => handleAddManagerAction(u.userId)}
                                >
                                  Add as Manager
                                </button>
                              )}
                            </div>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              )}
            </div>
          </div>
        </div>
      )}
    </section>
  );
}

