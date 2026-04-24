import { createContext, useCallback, useContext, useEffect, useMemo, useState } from "react";
import { setUnauthorizedHandler } from "../api/apiClient";
import { loginApi, logoutApi } from "../api/authApi";
import { getMeApi } from "../api/userApi";
import { TOKEN_STORAGE_KEY, clearToken, getToken, setToken } from "../auth/tokenStorage";

const AuthContext = createContext(null);

export function AuthProvider({ children }) {
  const [token, setTokenState] = useState(getToken());
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  const resetSession = useCallback(() => {
    clearToken();
    setTokenState(null);
    setUser(null);
    setLoading(false);
  }, []);

  useEffect(() => {
    setUnauthorizedHandler(() => {
      resetSession();
    });

    return () => {
      setUnauthorizedHandler(null);
    };
  }, [resetSession]);

  useEffect(() => {
    function onStorage(event) {
      if (event.key !== TOKEN_STORAGE_KEY) {
        return;
      }

      if (!event.newValue) {
        setTokenState(null);
        setUser(null);
        return;
      }

      setTokenState(event.newValue);
    }

    window.addEventListener("storage", onStorage);
    return () => window.removeEventListener("storage", onStorage);
  }, []);

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
        if (mounted) {
          resetSession();
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
  }, [token, resetSession]);

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
      resetSession();
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

