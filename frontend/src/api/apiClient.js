const API_PREFIX = "/api/v1";

export class ApiError extends Error {
  constructor(message, status, details) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.details = details;
  }
}

async function request(path, options = {}, token) {
  const headers = {
    "Content-Type": "application/json",
    ...(options.headers || {}),
  };

  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }

  const response = await fetch(`${API_PREFIX}${path}`, {
    ...options,
    headers,
  });

  const isJSON = response.headers.get("content-type")?.includes("application/json");
  const body = isJSON ? await response.json() : null;

  if (!response.ok) {
    const message = body?.error || body?.message || `Request failed: ${response.status}`;
    throw new ApiError(message, response.status, body);
  }

  return body;
}

export default request;

