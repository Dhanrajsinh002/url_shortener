function Footer() {
    const year = new Date().getFullYear();

    return (
        <footer className="site-footer">
            <p>
                &copy; {year} URL Shortener — built with Go, React, PostgreSQL &amp;
                Redis.
            </p>
        </footer>
    );
}

export default Footer;