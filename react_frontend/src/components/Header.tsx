import { Link, useNavigate } from "react-router-dom";

function Header() {
  const navigate = useNavigate();
  const isLoggedIn = !!localStorage.getItem("admin_token");

  const handleLogout = () => {
    localStorage.removeItem("admin_token");
    navigate("/");
  };

  return (
    <header className="site-header">
      <div className="site-header-inner">
        <Link to="/" className="brand">
          🔗 URL Shortener
        </Link>
        <nav className="site-nav">
          <a
            href="https://github.com/Dhanrajsinh002/url_shortener"
            target="_blank"
            rel="noreferrer"
          >
            GitHub
          </a>
          {isLoggedIn ? (
            <>
              <Link to="/admin">Dashboard</Link>
              <button className="nav-link-button" onClick={handleLogout}>
                Logout
              </button>
            </>
          ) : (
            <>
              <Link to="/admin/login">Login</Link>
              <Link to="/admin/register">Register</Link>
            </>
          )}
        </nav>
      </div>
    </header>
  );
}

export default Header;