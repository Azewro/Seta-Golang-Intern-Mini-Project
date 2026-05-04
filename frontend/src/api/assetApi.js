import request from "./apiClient";

export function listFoldersApi(token) {
  return request("/folders", { method: "GET" }, token);
}

export function createFolderApi(payload, token) {
  return request("/folders", {
    method: "POST",
    body: JSON.stringify(payload),
  }, token);
}

export function getFolderApi(folderId, token) {
  return request(`/folders/${folderId}`, { method: "GET" }, token);
}

export function updateFolderApi(folderId, payload, token) {
  return request(`/folders/${folderId}`, {
    method: "PATCH",
    body: JSON.stringify(payload),
  }, token);
}

export function deleteFolderApi(folderId, token) {
  return request(`/folders/${folderId}`, { method: "DELETE" }, token);
}

export function listNotesByFolderApi(folderId, token) {
  return request(`/folders/${folderId}/notes`, { method: "GET" }, token);
}

export function createNoteApi(folderId, payload, token) {
  return request(`/folders/${folderId}/notes`, {
    method: "POST",
    body: JSON.stringify(payload),
  }, token);
}

export function updateNoteApi(noteId, payload, token) {
  return request(`/notes/${noteId}`, {
    method: "PATCH",
    body: JSON.stringify(payload),
  }, token);
}

export function deleteNoteApi(noteId, token) {
  return request(`/notes/${noteId}`, { method: "DELETE" }, token);
}

export function shareAssetApi(payload, token) {
  return request("/shares", {
    method: "POST",
    body: JSON.stringify(payload),
  }, token);
}

export function revokeShareApi(shareId, token) {
  return request(`/shares/${shareId}`, { method: "DELETE" }, token);
}

export function listReceivedSharesApi(token) {
  return request("/shares/received", { method: "GET" }, token);
}

export function listGrantedSharesApi(token) {
  return request("/shares/granted", { method: "GET" }, token);
}
