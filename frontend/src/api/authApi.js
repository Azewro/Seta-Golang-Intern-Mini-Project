import request from "./apiClient";

export function registerApi(payload) {
  return request("/auth/register", {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

export function verifyEmailApi(token) {
  const query = new URLSearchParams({ token });
  return request(`/auth/verify-email?${query.toString()}`, { method: "GET" });
}

export function resendVerificationApi(payload) {
  return request("/auth/resend-verification", {
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
