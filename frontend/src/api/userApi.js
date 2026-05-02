import request from "./apiClient";

export function getMeApi(token) {
  return request("/users/me", { method: "GET" }, token);
}

export function listUsersApi(token, page = 1, limit = 20) {
  return request(`/users?page=${page}&limit=${limit}`, { method: "GET" }, token);
}

export function bulkGetUsersApi(userIds, token) {
  return request(`/users/bulk`, {
    method: "POST",
    body: JSON.stringify({ userIds })
  }, token);
}
