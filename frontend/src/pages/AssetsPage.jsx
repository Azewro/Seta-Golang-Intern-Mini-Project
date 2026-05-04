import { useEffect, useMemo, useState } from "react";
import { useAuth } from "../context/AuthContext";
import { useToast } from "../context/ToastContext";
import {
  createFolderApi,
  createNoteApi,
  deleteFolderApi,
  deleteNoteApi,
  listFoldersApi,
  listGrantedSharesApi,
  listNotesByFolderApi,
  listReceivedSharesApi,
  revokeShareApi,
  shareAssetApi,
  updateFolderApi,
  updateNoteApi,
} from "../api/assetApi";

export default function AssetsPage() {
  const { token } = useAuth();
  const { pushToast } = useToast();

  const [folders, setFolders] = useState([]);
  const [foldersLoading, setFoldersLoading] = useState(true);
  const [selectedFolderId, setSelectedFolderId] = useState(null);
  const [notes, setNotes] = useState([]);
  const [notesLoading, setNotesLoading] = useState(false);
  const [receivedShares, setReceivedShares] = useState([]);
  const [grantedShares, setGrantedShares] = useState([]);

  const [newFolderName, setNewFolderName] = useState("");
  const [newNoteTitle, setNewNoteTitle] = useState("");
  const [newNoteContent, setNewNoteContent] = useState("");

  const [shareForm, setShareForm] = useState({
    assetType: "folder",
    assetId: "",
    sharedWithUserId: "",
    accessLevel: "read",
  });

  const selectedFolder = useMemo(
    () => folders.find((f) => f.folderId === selectedFolderId) || null,
    [folders, selectedFolderId]
  );

  const loadFolders = async () => {
    try {
      setFoldersLoading(true);
      const data = await listFoldersApi(token);
      const list = data?.data || [];
      setFolders(list);
      if (list.length > 0 && !selectedFolderId) {
        setSelectedFolderId(list[0].folderId);
      }
      if (list.length === 0) {
        setSelectedFolderId(null);
      }
    } catch (err) {
      pushToast(err.message || "Failed to load folders", "error");
    } finally {
      setFoldersLoading(false);
    }
  };

  const loadNotes = async (folderId) => {
    if (!folderId) {
      setNotes([]);
      return;
    }
    try {
      setNotesLoading(true);
      const data = await listNotesByFolderApi(folderId, token);
      setNotes(data?.data || []);
    } catch (err) {
      setNotes([]);
      pushToast(err.message || "Failed to load notes", "error");
    } finally {
      setNotesLoading(false);
    }
  };

  const loadShares = async () => {
    try {
      const [received, granted] = await Promise.all([
        listReceivedSharesApi(token),
        listGrantedSharesApi(token),
      ]);
      setReceivedShares(received?.data || []);
      setGrantedShares(granted?.data || []);
    } catch (err) {
      pushToast(err.message || "Failed to load shares", "error");
    }
  };

  useEffect(() => {
    if (!token) return;
    loadFolders();
    loadShares();
  }, [token]);

  useEffect(() => {
    if (!token) return;
    loadNotes(selectedFolderId);
  }, [selectedFolderId, token]);

  const handleCreateFolder = async (event) => {
    event.preventDefault();
    if (!newFolderName.trim()) return;
    try {
      const folder = await createFolderApi({ name: newFolderName }, token);
      setNewFolderName("");
      pushToast("Folder created", "success");
      await loadFolders();
      if (folder?.folderId) {
        setSelectedFolderId(folder.folderId);
      }
    } catch (err) {
      pushToast(err.message, "error");
    }
  };

  const handleCreateNote = async (event) => {
    event.preventDefault();
    if (!selectedFolderId || !newNoteTitle.trim()) return;
    try {
      await createNoteApi(
        selectedFolderId,
        { title: newNoteTitle, content: newNoteContent },
        token
      );
      setNewNoteTitle("");
      setNewNoteContent("");
      pushToast("Note created", "success");
      await loadNotes(selectedFolderId);
    } catch (err) {
      pushToast(err.message, "error");
    }
  };

  const handleRenameFolder = async (folderId, currentName) => {
    const nextName = window.prompt("Rename folder", currentName);
    if (!nextName || nextName.trim() === currentName) return;
    try {
      await updateFolderApi(folderId, { name: nextName.trim() }, token);
      pushToast("Folder updated", "success");
      await loadFolders();
    } catch (err) {
      pushToast(err.message, "error");
    }
  };

  const handleDeleteFolder = async (folderId) => {
    if (!window.confirm("Delete this folder and its notes?")) return;
    try {
      await deleteFolderApi(folderId, token);
      pushToast("Folder deleted", "success");
      if (selectedFolderId === folderId) {
        setSelectedFolderId(null);
      }
      await loadFolders();
    } catch (err) {
      pushToast(err.message, "error");
    }
  };

  const handleEditNote = async (note) => {
    const nextTitle = window.prompt("Edit note title", note.title);
    if (!nextTitle) return;
    const nextContent = window.prompt("Edit note content", note.content || "");
    if (nextContent === null) return;
    try {
      await updateNoteApi(
        note.noteId,
        { title: nextTitle.trim(), content: nextContent },
        token
      );
      pushToast("Note updated", "success");
      await loadNotes(selectedFolderId);
    } catch (err) {
      pushToast(err.message, "error");
    }
  };

  const handleDeleteNote = async (noteId) => {
    if (!window.confirm("Delete this note?")) return;
    try {
      await deleteNoteApi(noteId, token);
      pushToast("Note deleted", "success");
      await loadNotes(selectedFolderId);
    } catch (err) {
      pushToast(err.message, "error");
    }
  };

  const handleShareSubmit = async (event) => {
    event.preventDefault();
    try {
      await shareAssetApi(
        {
          assetType: shareForm.assetType,
          assetId: Number(shareForm.assetId),
          sharedWithUserId: Number(shareForm.sharedWithUserId),
          accessLevel: shareForm.accessLevel,
        },
        token
      );
      pushToast("Share updated", "success");
      setShareForm((prev) => ({ ...prev, assetId: "", sharedWithUserId: "" }));
      await loadShares();
    } catch (err) {
      pushToast(err.message, "error");
    }
  };

  const handleRevokeShare = async (shareId) => {
    if (!window.confirm("Revoke this share?")) return;
    try {
      await revokeShareApi(shareId, token);
      pushToast("Share revoked", "success");
      await loadShares();
    } catch (err) {
      pushToast(err.message, "error");
    }
  };

  return (
    <section className="page-panel">
      <div className="card" style={{ width: "100%" }}>
        <h2>Assets</h2>
        <p className="muted">
          Manage folders, notes, and sharing permissions (read/write).
        </p>

        <form
          className="form-group"
          onSubmit={handleCreateFolder}
          style={{ display: "flex", gap: "12px", marginBottom: "20px" }}
        >
          <input
            type="text"
            placeholder="New folder name"
            value={newFolderName}
            onChange={(e) => setNewFolderName(e.target.value)}
          />
          <button type="submit" className="primary" style={{ margin: 0, whiteSpace: "nowrap" }}>
            Create Folder
          </button>
        </form>

        <div style={{ display: "grid", gap: "16px", gridTemplateColumns: "minmax(280px, 1fr) minmax(0, 2fr)" }}>
          <div className="table-wrap">
            <table className="table">
              <thead>
                <tr>
                  <th>Folder</th>
                  <th>Access</th>
                </tr>
              </thead>
              <tbody>
                {foldersLoading && (
                  <tr>
                    <td colSpan={2}>Loading folders...</td>
                  </tr>
                )}
                {!foldersLoading && folders.length === 0 && (
                  <tr>
                    <td colSpan={2}>No folders yet.</td>
                  </tr>
                )}
                {!foldersLoading &&
                  folders.map((folder) => (
                    <tr
                      key={folder.folderId}
                      onClick={() => setSelectedFolderId(folder.folderId)}
                      style={{
                        cursor: "pointer",
                        background:
                          selectedFolderId === folder.folderId
                            ? "rgba(45, 212, 191, 0.08)"
                            : "transparent",
                      }}
                    >
                      <td>{folder.name}</td>
                      <td>
                        <span className="badge badge-member">{folder.accessLevel}</span>
                      </td>
                    </tr>
                  ))}
              </tbody>
            </table>
          </div>

          <div>
            {!selectedFolder && <p className="muted">Select a folder to view notes.</p>}
            {selectedFolder && (
              <div style={{ display: "grid", gap: "12px" }}>
                <div
                  style={{
                    border: "1px solid var(--surface-border)",
                    borderRadius: "16px",
                    padding: "14px",
                    background: "rgba(15, 23, 42, 0.4)",
                  }}
                >
                  <div style={{ display: "flex", alignItems: "center", justifyContent: "space-between", gap: "8px" }}>
                    <div>
                      <h3 style={{ margin: 0 }}>{selectedFolder.name}</h3>
                      <p className="muted" style={{ margin: "4px 0 0 0" }}>
                        Access: {selectedFolder.accessLevel}
                      </p>
                    </div>
                    {selectedFolder.canWrite && (
                      <div style={{ display: "flex", gap: "8px" }}>
                        <button
                          className="primary"
                          style={{ margin: 0, backgroundColor: "#475569", color: "#fff" }}
                          onClick={() => handleRenameFolder(selectedFolder.folderId, selectedFolder.name)}
                        >
                          Rename
                        </button>
                        <button
                          className="primary"
                          style={{ margin: 0, backgroundColor: "#881337", color: "#fff" }}
                          onClick={() => handleDeleteFolder(selectedFolder.folderId)}
                        >
                          Delete
                        </button>
                      </div>
                    )}
                  </div>
                </div>

                {selectedFolder.canWrite && (
                  <form className="form-grid" onSubmit={handleCreateNote}>
                    <input
                      type="text"
                      placeholder="Note title"
                      value={newNoteTitle}
                      onChange={(e) => setNewNoteTitle(e.target.value)}
                    />
                    <input
                      type="text"
                      placeholder="Note content"
                      value={newNoteContent}
                      onChange={(e) => setNewNoteContent(e.target.value)}
                    />
                    <button className="primary" type="submit">
                      Create Note
                    </button>
                  </form>
                )}

                <div className="table-wrap">
                  <table className="table">
                    <thead>
                      <tr>
                        <th>Title</th>
                        <th>Access</th>
                        <th>Actions</th>
                      </tr>
                    </thead>
                    <tbody>
                      {notesLoading && (
                        <tr>
                          <td colSpan={3}>Loading notes...</td>
                        </tr>
                      )}
                      {!notesLoading && notes.length === 0 && (
                        <tr>
                          <td colSpan={3}>No notes in this folder.</td>
                        </tr>
                      )}
                      {!notesLoading &&
                        notes.map((note) => (
                          <tr key={note.noteId}>
                            <td>
                              <div>{note.title}</div>
                              <div className="muted" style={{ fontSize: "0.85rem" }}>
                                {note.content}
                              </div>
                            </td>
                            <td>
                              <span className="badge badge-member">{note.accessLevel}</span>
                            </td>
                            <td>
                              {note.canWrite ? (
                                <div style={{ display: "flex", gap: "8px" }}>
                                  <button
                                    className="primary"
                                    style={{ margin: 0, backgroundColor: "#475569", color: "#fff" }}
                                    onClick={() => handleEditNote(note)}
                                  >
                                    Edit
                                  </button>
                                  <button
                                    className="primary"
                                    style={{ margin: 0, backgroundColor: "#881337", color: "#fff" }}
                                    onClick={() => handleDeleteNote(note.noteId)}
                                  >
                                    Delete
                                  </button>
                                </div>
                              ) : (
                                <span className="muted">Read-only</span>
                              )}
                            </td>
                          </tr>
                        ))}
                    </tbody>
                  </table>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>

      <div className="card" style={{ width: "100%", marginTop: "20px" }}>
        <h3>Share Asset</h3>
        <form className="form-grid" onSubmit={handleShareSubmit}>
          <select
            value={shareForm.assetType}
            onChange={(e) => setShareForm((prev) => ({ ...prev, assetType: e.target.value }))}
          >
            <option value="folder">Folder</option>
            <option value="note">Note</option>
          </select>
          <input
            type="number"
            placeholder="Asset ID"
            value={shareForm.assetId}
            onChange={(e) => setShareForm((prev) => ({ ...prev, assetId: e.target.value }))}
            required
          />
          <input
            type="number"
            placeholder="Shared with user ID"
            value={shareForm.sharedWithUserId}
            onChange={(e) => setShareForm((prev) => ({ ...prev, sharedWithUserId: e.target.value }))}
            required
          />
          <select
            value={shareForm.accessLevel}
            onChange={(e) => setShareForm((prev) => ({ ...prev, accessLevel: e.target.value }))}
          >
            <option value="read">Read</option>
            <option value="write">Write</option>
          </select>
          <button className="primary" type="submit">
            Share
          </button>
        </form>

        <div style={{ display: "grid", gap: "16px", marginTop: "16px", gridTemplateColumns: "1fr 1fr" }}>
          <div className="table-wrap">
            <table className="table">
              <thead>
                <tr>
                  <th colSpan={3}>Granted Shares</th>
                </tr>
              </thead>
              <tbody>
                {grantedShares.length === 0 && (
                  <tr>
                    <td colSpan={3}>No granted shares.</td>
                  </tr>
                )}
                {grantedShares.map((share) => (
                  <tr key={`g-${share.shareId}`}>
                    <td>{share.assetType} #{share.assetId}</td>
                    <td>U#{share.sharedWithUserId} ({share.accessLevel})</td>
                    <td>
                      <button
                        className="primary"
                        style={{ margin: 0, backgroundColor: "#881337", color: "#fff" }}
                        onClick={() => handleRevokeShare(share.shareId)}
                      >
                        Revoke
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          <div className="table-wrap">
            <table className="table">
              <thead>
                <tr>
                  <th colSpan={2}>Received Shares</th>
                </tr>
              </thead>
              <tbody>
                {receivedShares.length === 0 && (
                  <tr>
                    <td colSpan={2}>No received shares.</td>
                  </tr>
                )}
                {receivedShares.map((share) => (
                  <tr key={`r-${share.shareId}`}>
                    <td>{share.assetType} #{share.assetId}</td>
                    <td>From U#{share.sharedByUserId} ({share.accessLevel})</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </section>
  );
}
