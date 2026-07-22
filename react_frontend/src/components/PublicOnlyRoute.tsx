import { Navigate } from "react-router-dom";
import type { ReactNode } from "react";

// Inverse of ProtectedRoute — for pages that should only be visible to
// logged-OUT users (the public shortener, login, register). A logged-in
// user hitting any of these gets sent straight to the dashboard instead.
function PublicOnlyRoute({ children }: { children: ReactNode }) {
  const token = localStorage.getItem("admin_token");
  if (token) {
    return <Navigate to="/admin" replace />;
  }
  return <>{children}</>;
}

export default PublicOnlyRoute;