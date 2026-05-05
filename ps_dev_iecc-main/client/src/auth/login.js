import React, { useState } from "react";
import { GoogleLogin } from "@react-oauth/google";
import { useMsal } from "@azure/msal-react";
import Logo from "../assets/img/logo.png";
import "./style.css";
import InputBox from "../components/input";
import CustomButton from "../components/button";
import { useAuth } from "./AuthContext";
import { apiPostRequest } from "../utils/api";
import { appBase, portalName } from "../utils/settings";
import { encryptCourseId } from "../utils/crypto";
import MicrosoftLogo from "../assets/img/microsoft.png";

const AuthLogin = () => {
  const { login } = useAuth();
  const { instance } = useMsal();
  const [loginError, setLoginError] = useState("");
  const isMicrosoftEnabled =
    process.env.REACT_APP_IS_MICROSOFT_ENABLED === "true";

  const finishLogin = (data) => {
    login();
    localStorage.setItem(
      "PS1",
      encryptCourseId(`${data.name}-**-**-${data.profile}`),
    );
    window.location = appBase + "/dashboard";
  };

  // ---- Google ----
  const handleGoogleSuccess = async (response) => {
    const token = response.credential;
    try {
      const res = await apiPostRequest("/auth/GLogin", { id_token: token });
      if (res.success) finishLogin(res.data);
      else setLoginError(res.error || "Login failed, please try again.");
    } catch (error) {
      console.error("Error during Google authentication:", error);
      setLoginError("Google login failed. Please try again.");
    }
  };

  const handleGoogleError = (error) => {
    console.error("Google Login failed:", error);
    setLoginError("Google login failed. Please try again.");
  };

  // ---- Microsoft (MSAL) ----
  const handleMicrosoftLogin = async () => {
    setLoginError("");
    try {
      const result = await instance.loginPopup({
        authority: "https://login.microsoftonline.com/organizations",
        scopes: ["openid", "profile", "email"], // add API scopes if needed
      });
      const idToken = result.idToken; // send ID token to your backend
      const res = await apiPostRequest("/auth/MSLogin", { id_token: idToken });
      if (res.success) finishLogin(res.data);
      else setLoginError(res.error || "Login failed, please try again.");
    } catch (e) {
      console.error("Microsoft Login failed:", e);
      setLoginError("Microsoft login failed. Please try again.");
    }
  };

  return (
    <div className="w-screen h-screen fixed bg-background flex items-center justify-center">
      <div className="sm:bg-white sm:p-4 p-0 sm:px-10 px-5 w-100 rounded-xl sm:shadow-sm flex flex-col items-center">
        <div className="flex gap-2 items-center">
          <img width={40} height={40} src={Logo} alt="logo" />
          <h2 className="text-xl font-medium">{portalName}</h2>
        </div>

        <h2 className="mt-4 text-xl font-semibold text-primary">
          Hi, Welcome Back!
        </h2>

        {loginError && (
          <div
            className="bg-red text-white px-4 py-1 rounded mt-3 w-full"
            style={{ maxWidth: "500px" }}
          >
            <h1 className="text-md">{loginError}</h1>
          </div>
        )}

        <InputBox
          margin="mt-3"
          label="Username"
          placeholder="Enter your username"
          type="text"
          onChange={() => {}}
        />
        <InputBox
          margin="mt-5"
          label="Password"
          placeholder="Enter your password"
          type="password"
          onChange={() => {}}
        />

        <CustomButton margin="mt-5" label="Login" />

        <h2 className="m-3 text-sm">Or</h2>

        <div className="flex flex-row items-center justify-center w-full gap-2">
          {/* Google (Provider should be at root) */}
          <GoogleLogin
            onSuccess={handleGoogleSuccess}
            onError={handleGoogleError}
          />

          {/* Microsoft (MSAL) */}
          {isMicrosoftEnabled && (
            <div className="w-full" style={{ maxWidth: 320 }}>
              <button
                onClick={handleMicrosoftLogin}
                className="w-full flex items-center justify-center gap-2 border rounded-sm py-1.5"
              >
                <img
                  className="w-6 h-6"
                  src={MicrosoftLogo}
                  alt="Microsoft Logo"
                />
                <span className="text-sm">Sign in with Microsoft</span>
              </button>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default AuthLogin;
