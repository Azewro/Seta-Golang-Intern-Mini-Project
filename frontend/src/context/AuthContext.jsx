import { createContext, useContext, useEffect, useMemo, useState } from "react";
import { loginApi, logoutApi } from "../api/authApi";
import { getMeApi } from "../api/userApi";
import { clearToken, getToken, setToken } from "../auth/tokenStorage";

const AuthContext = createContext(null);

export function AuthProvider({ children }) {
  const [token, setTokenState] = useState(getToken());
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let mounted = true;

    async function bootstrap() {
      if (!token) {
        if (mounted) {
          setUser(null);
          setLoading(false);
        }
        return;
      }

      try {
        const me = await getMeApi(token);
        if (mounted) {
          setUser(me);
        }
      } catch {
        clearToken();
        if (mounted) {
          setTokenState(null);
          setUser(null);
        }
      } finally {
        if (mounted) {
          setLoading(false);
        }
      }
    }

    bootstrap();
    return () => {
      mounted = false;
    };
  }, [token]);

  const login = async (payload) => {
    const data = await loginApi(payload);
    setToken(data.accessToken);
    setTokenState(data.accessToken);
    setUser(data.user);
    return data;
  };

  const logout = async () => {
    try {
      if (token) {
        await logoutApi(token);
      }
    } finally {
      clearToken();
      setTokenState(null);
      setUser(null);
    }
  };

  const value = useMemo(
    () => ({
      token,
      user,
      loading,
      isAuthenticated: Boolean(token && user),
      login,
      logout,
      setUser,
    }),
    [token, user, loading]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within AuthProvider");
  }
  return context;
}

