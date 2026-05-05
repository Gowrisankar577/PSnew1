import React, { createContext, useContext, useState } from "react";

const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
  const [isAuthenticated, setIsAuthenticated] = useState(() => {
    return !!localStorage.getItem("PS");
  });

  const login = () => {
    setIsAuthenticated(true);
    localStorage.setItem("PS", "true");
  };

  const logout = () => {
    setIsAuthenticated(false);
    localStorage.clear();
    window.location.href = "/auth/login";
  };

  return (
    <AuthContext.Provider value={{ isAuthenticated, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => useContext(AuthContext);
