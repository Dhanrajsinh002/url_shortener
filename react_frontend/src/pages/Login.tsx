import { useState } from "react";
import { useNavigate } from "react-router-dom";
import axios from "axios";
import Button from "../components/Button";
import usePageMeta from "../hooks/usePageMeta";

const API_BASE = "http://localhost:8000";

function Login() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState("");
  const navigate = useNavigate();

  usePageMeta("Admin Login — URL Shortener", "/favicon-lock.svg");

  const handleLogin = () => {
    setError("");
    axios
      .post(`${API_BASE}/admin/login`, { username, password })
      .then((res) => {
        localStorage.setItem("admin_token", res.data.token);
        navigate("/admin");
      })
      .catch(() => setError("Invalid username or password."));
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter") handleLogin();
  };

  return (
    <main id="center">
      <h1>Admin Login</h1>

      <div className="auth-form">
        <input
          className="auth-input"
          placeholder="Username"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          onKeyDown={handleKeyDown}
        />

        <div className="password-field">
          <input
            className="auth-input"
            type={showPassword ? "text" : "password"}
            placeholder="Password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            onKeyDown={handleKeyDown}
          />
          <button
            type="button"
            className="password-toggle"
            onClick={() => setShowPassword((v) => !v)}
            aria-label={showPassword ? "Hide password" : "Show password"}
          >
            {showPassword ? (
              <svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M3 3l18 18M10.58 10.58a2 2 0 002.83 2.83M9.88 4.24A9.77 9.77 0 0112 4c5 0 9 4 10 8-.32 1.2-.87 2.32-1.6 3.31M6.61 6.61C4.4 8 2.83 10 2 12c1 4 5 8 10 8 1.35 0 2.63-.28 3.79-.79" />
              </svg>
            ) : (
              <svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M2 12s4-8 10-8 10 8 10 8-4 8-10 8-10-8-10-8z" />
                <circle cx="12" cy="12" r="3" />
              </svg>
            )}
          </button>
        </div>

        {error && <p className="error">{error}</p>}

        <Button label="Login" click={handleLogin} />
      </div>
    </main>
  );
}

export default Login;