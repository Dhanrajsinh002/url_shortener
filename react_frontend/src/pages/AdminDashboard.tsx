import { useEffect, useState } from "react";
import apiClient from "../api/client";
import usePageMeta from "../hooks/usePageMeta";

type UrlRecord = {
    short_code: string;
    long_url: string;
    click_count: number;
    created_at: string;
}

function AdminDashboard() {
    const [urls, setUrls] = useState<UrlRecord[]>([]);
    const [error, setError] = useState("");

    usePageMeta("Admin Dashboard — URL Shortener", "/favicon-dashboard.svg");

    useEffect(() => {
        apiClient
            .get<UrlRecord[]>("/admin/urls")
            .then((res) => setUrls(res.data))
            .catch(() => setError("Failed to load Urls. Please log in again."))
    }, []);

    return (
        <main id="center">
            <h1>Admin Dashboard</h1>
            {error && <p className="error">{error}</p>}
            <table>
                <thead>
                    <tr>
                        <th>Short Code</th>
                        <th>Long Url</th>
                        <th>Clicks</th>
                        <th>Created On</th>
                    </tr>
                </thead>
                <body>
                    {urls.map((u) => (
                        <tr key={u.short_code}>
                            <td>{u.short_code}</td>
                            <td>{u.long_url}</td>
                            <td>{u.click_count}</td>
                            <td>{new Date(u.created_at).toLocaleString()}</td>
                        </tr>
                    ))}
                </body>
            </table>
        </main>
    );
}

export default AdminDashboard;