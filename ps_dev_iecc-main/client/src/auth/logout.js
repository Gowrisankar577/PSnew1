import React, { useEffect } from "react";
import "./style.css";
import { useAuth } from "./AuthContext";

const AuthLogout = () => {
  const { logout } = useAuth();

  useEffect(() => {
    logout();
  }, []);

  return (
    <div className="w-screen h-screen fixed bg-background flex items-center justify-center"></div>
  );
};

export default AuthLogout;
