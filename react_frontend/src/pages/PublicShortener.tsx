import React, { useState } from "react";
import axios from "axios";
import Button from "../components/Button";
import usePageMeta from "../hooks/usePageMeta";
import "../App.css";

const API_BASE = "http://localhost:8000";

type CreateShortUrlResponse = {
  message: string;
  short_url: string;
  click_count: string;
};

function PublicShortener() {
  const [longUrl, setLongUrl] = useState("");
  const [shortUrl, setShortUrl] = useState("");
  const [clickCount, setClickCount] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const [copied, setCopied] = useState(false);

  usePageMeta("URL Shortener — Shorten a Link", "/favicon.svg");

  const handleShorten = () => {
    const trimmed = longUrl.trim();
    if (!trimmed) {
      setError("Please enter a URL.");
      return;
    }

    setLoading(true);
    setError("");
    setShortUrl("");
    setClickCount("");
    setCopied(false);

    axios
      .post<CreateShortUrlResponse>(`${API_BASE}/create-short-url`, {
        long_url: trimmed,
      })
      .then((response) => {
        setShortUrl(response.data.short_url);
        setClickCount(response.data.click_count);
      })
      .catch((err) => {
        const message =
          err?.response?.data?.error ??
          "Something went wrong. Is the backend running?";
        setError(message);
      })
      .finally(() => setLoading(false));
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter") {
      handleShorten();
    }
  };

  const handleCopy = () => {
    navigator.clipboard.writeText(shortUrl).then(() => {
      setCopied(true);
      setTimeout(() => setCopied(false), 1500);
    });
  };

  return (
    <>
      <main id="center">
        <h1>URL Shortener</h1>
        <p>Paste a long link below and get a short one back.</p>

        <div className="shortener">
          <input
            className="url-input"
            type="text"
            placeholder="https://example.com/some/very/long/link"
            value={longUrl}
            onChange={(e) => setLongUrl(e.target.value)}
            onKeyDown={handleKeyDown}
          />
          <Button label={loading ? "Shortening..." : "Shorten"} click={handleShorten} />
        </div>

        {error && <p className="error">{error}</p>}

        {shortUrl && (
          <div className="result">
            <a href={shortUrl} target="_blank" rel="noreferrer">
              {shortUrl}
            </a>
            <Button label={copied ? "Copied!" : "Copy"} click={handleCopy} />
          </div>
        )}

        {clickCount && (
          <div className="result">
            Total Clicks: {clickCount}
          </div>
        )}
      </main>
    </>
  );
}

export default PublicShortener;
