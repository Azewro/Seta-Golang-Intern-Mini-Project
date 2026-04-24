export function validateEmail(email) {
  const value = String(email || "").trim();
  if (!value) {
    return "Email is required";
  }

  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  if (!emailRegex.test(value)) {
    return "Email format is invalid";
  }

  return "";
}

export function validatePassword(password) {
  const value = String(password || "");
  if (!value) {
    return "Password is required";
  }
  if (value.length < 6) {
    return "Password must be at least 6 characters";
  }
  return "";
}

export function validateUsername(username) {
  const value = String(username || "").trim();
  if (!value) {
    return "Username is required";
  }
  if (value.length < 3) {
    return "Username must be at least 3 characters";
  }
  return "";
}

export function validateRole(role) {
  if (role !== "member" && role !== "manager") {
    return "Role must be member or manager";
  }
  return "";
}

