import { useEffect, useState } from "react";
import apiClient from "../api/client";
import usePageMeta from "../hooks/usePageMeta";
import Button from "../components/Button";

type UrlRecord = {
    short_code: string;
    long_url: string;
    click_count: number;
    created_at: string;
}

type CreateShortUrlResponse = {
    message: string;
    short_url: string;
}

function AdminDashboard() {
    const [urls, setUrls] = useState<UrlRecord[]>([]);
    const [error, setError] = useState("");

    const [longUrl, setLongUrl] = useState("");
    const [createError, setCreateError] = useState("");
    const [createdUrl, setCreatedUrl] = useState("");
    const [creating, setCreating] = useState(false);

    usePageMeta("Admin Dashboard — URL Shortener", "/favicon-dashboard.svg");

    const fetchUrls = () => {
        apiClient
            .get<UrlRecord[]>("/admin/urls")
            .then((res) => setUrls(res.data))
            .catch(() => setError("Failed to load URLs. Please Login again."))
    };

    useEffect(() => {
        fetchUrls
    }, []);

    const handleCreate = () => {
        setCreateError("");
        setCreatedUrl("");

        if(!longUrl.trim()) {
            setCreateError("Please enter URL.")
            return;
        }

        setCreating(true);

        apiClient
            .post<CreateShortUrlResponse>("/admin/create-short-url", {
                long_url: longUrl
            })
            .then((res) => {
                setCreatedUrl(res.data.short_url)
                setLongUrl("");
                fetchUrls();
            })
            .catch((err) => {
                const msg = err?.response?.data?.error ?? "Could not create short URL.";
                setCreateError(msg);
            })
            .finally(() => setCreating(false));
    };

    return (
        <main id="center">
            <h1>User Dashboard</h1>

            <div className="auth-form">
                <input 
                    className="auth-input"
                    placeholder="https://example.com/some/very/long/url"
                    value={longUrl}
                    onChange={(e) => setLongUrl(e.target.value)}
                    onKeyDown={(e) => e.key === "Enter" && handleCreate()} />

                {createError && <p className="error">{createError}</p>}
                {createdUrl && (
                    <p className="result">
                        Created:{" "}
                        <a href="{createdUrl}" target="_blank" rel="noreferrer">
                            {createdUrl}
                        </a>
                    </p>
                )}

                <Button 
                    label={ creating ? "Creating..." : "Create Short URL"}
                    click={handleCreate}
                />
            </div>

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