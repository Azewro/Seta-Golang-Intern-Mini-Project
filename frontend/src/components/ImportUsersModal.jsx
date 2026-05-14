import { useState, useRef } from "react";
import { importUsersApi } from "../api/userApi";

export default function ImportUsersModal({ token, onClose, onImportSuccess }) {
  const [file, setFile] = useState(null);
  const [uploading, setUploading] = useState(false);
  const [error, setError] = useState("");
  const [result, setResult] = useState(null);
  const fileInputRef = useRef(null);

  const handleFileChange = (e) => {
    setError("");
    const selectedFile = e.target.files[0];
    if (!selectedFile) {
      setFile(null);
      return;
    }
    // Max 3MB
    if (selectedFile.size > 3 * 1024 * 1024) {
      setError("File exceeds maximum size of 3MB.");
      setFile(null);
      return;
    }
    setFile(selectedFile);
  };

  const handleUpload = async () => {
    if (!file) return;
    setUploading(true);
    setError("");
    try {
      const res = await importUsersApi(file, token);
      setResult(res);
      if (res.success > 0) {
        onImportSuccess();
      }
    } catch (err) {
      setError(err.message || "Failed to upload file");
    } finally {
      setUploading(false);
    }
  };

  return (
    <div style={overlayStyle}>
      <div className="card" style={modalStyle}>
        <button style={closeButtonStyle} onClick={onClose}>&times;</button>
        <h3 style={{ marginTop: 0 }}>Bulk Import Users</h3>

        {!result ? (
          <div>
            <p className="muted" style={{ marginBottom: "1.5rem" }}>
              Upload a .csv file with columns: <code>username, email, password, role</code>.
            </p>

            <input
              type="file"
              accept=".csv"
              ref={fileInputRef}
              onChange={handleFileChange}
              style={{ marginBottom: "1rem" }}
            />
            {error && <p className="field-error" style={{ marginBottom: "1rem" }}>{error}</p>}

            <button
              className="primary"
              onClick={handleUpload}
              disabled={!file || uploading}
              style={{ width: "100%" }}
            >
              {uploading ? (
                <span><span className="spinner" style={{ width: "16px", height: "16px" }}></span> Uploading...</span>
              ) : "Upload CSV"}
            </button>
          </div>
        ) : (
          <div>
            <div style={{ display: "flex", gap: "1rem", marginBottom: "1.5rem" }}>
              <div className="state-card" style={{ flex: 1, padding: "1rem" }}>
                <div style={{ fontSize: "2rem", fontWeight: "bold" }}>{result.success}</div>
                <div>Success</div>
              </div>
              <div className="state-card" style={{ flex: 1, padding: "1rem", background: "rgba(251, 113, 133, 0.1)", borderColor: "rgba(251, 113, 133, 0.2)", color: "var(--danger)" }}>
                <div style={{ fontSize: "2rem", fontWeight: "bold" }}>{result.failed}</div>
                <div>Failed</div>
              </div>
            </div>

            {result.errors && result.errors.length > 0 && (
              <div style={{ marginBottom: "1rem" }}>
                <h4 style={{ margin: "0 0 0.5rem 0" }}>Errors:</h4>
                <div className="table-wrap" style={{ maxHeight: "200px", overflowY: "auto" }}>
                  <table className="table" style={{ fontSize: "0.85rem" }}>
                    <thead>
                      <tr>
                        <th style={{ padding: "0.5rem" }}>Row</th>
                        <th style={{ padding: "0.5rem" }}>Email</th>
                        <th style={{ padding: "0.5rem" }}>Message</th>
                      </tr>
                    </thead>
                    <tbody>
                      {result.errors.map((e, idx) => (
                        <tr key={idx}>
                          <td style={{ padding: "0.5rem" }}>{e.rowNumber}</td>
                          <td style={{ padding: "0.5rem" }}>{e.email || "-"}</td>
                          <td style={{ padding: "0.5rem", color: "var(--danger)" }}>{e.message}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
                {result.errorsTruncated && (
                  <p className="muted" style={{ fontSize: "0.85rem", marginTop: "0.5rem" }}>* Only showing first 50 errors.</p>
                )}
              </div>
            )}

            <button className="primary" onClick={onClose} style={{ width: "100%" }}>Done</button>
          </div>
        )}
      </div>
    </div>
  );
}

const overlayStyle = {
  position: "fixed",
  top: 0, left: 0, right: 0, bottom: 0,
  backgroundColor: "rgba(3, 7, 18, 0.8)",
  backdropFilter: "blur(4px)",
  display: "flex",
  alignItems: "center",
  justifyContent: "center",
  zIndex: 1000,
  padding: "1rem"
};

const modalStyle = {
  position: "relative",
  width: "100%",
  maxWidth: "500px",
  margin: "0 auto",
  padding: "2rem",
};

const closeButtonStyle = {
  position: "absolute",
  top: "1rem",
  right: "1.5rem",
  background: "none",
  border: "none",
  color: "var(--text-muted)",
  fontSize: "2rem",
  cursor: "pointer",
  lineHeight: 1
};

