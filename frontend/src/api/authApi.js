import request from "./apiClient";

export function registerApi(payload) {
  return request("/auth/register", {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

export function loginApi(payload) {
  return request("/auth/login", {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

export function logoutApi(token) {
  return request("/auth/logout", { method: "POST" }, token);
}

