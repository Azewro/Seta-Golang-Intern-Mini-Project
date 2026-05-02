import request from "./apiClient";

export function listMyTeamsApi(token) {
  return request("/teams/my", { method: "GET" }, token);
}

export function getTeamApi(teamId, token) {
  return request(`/teams/${teamId}`, { method: "GET" }, token);
}

export function createTeamApi(payload, token) {
  return request("/teams", {
    method: "POST",
    body: JSON.stringify(payload),
  }, token);
}

export function addMemberApi(teamId, userId, token) {
  return request(`/teams/${teamId}/members`, {
    method: "POST",
    body: JSON.stringify({ userId }),
  }, token);
}

export function removeMemberApi(teamId, userId, token) {
  return request(`/teams/${teamId}/members/${userId}`, {
    method: "DELETE",
  }, token);
}

export function addManagerApi(teamId, userId, token) {
  return request(`/teams/${teamId}/managers`, {
    method: "POST",
    body: JSON.stringify({ userId }),
  }, token);
}

export function removeManagerApi(teamId, userId, token) {
  return request(`/teams/${teamId}/managers/${userId}`, {
    method: "DELETE",
  }, token);
}
