import { Navigate } from "react-router-dom";
import type { ReactNode } from "react";

function ProtectedRoute({ children }: { children: ReactNode }) {
    const token = localStorage.getItem("admin_toke");
    if (!token) {
        return <Navigate to="/admin/login" replace/>;
    }
    return <>{children}</>;
}

export default ProtectedRoute;