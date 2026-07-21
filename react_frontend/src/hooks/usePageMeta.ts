import { useEffect } from "react";

/**
 * Updates the browser tab's title and favicon for the current page.
 * SPAs don't get this for free per-route — each page calls this hook
 * with its own title/icon, and it resets nothing on unmount since the
 * next page mounting will just overwrite it again.
 */
function usePageMeta(title: string, faviconHref?: string) {
  useEffect(() => {
    document.title = title;

    if (faviconHref) {
      let link = document.querySelector<HTMLLinkElement>("link[rel='icon']");
      if (!link) {
        link = document.createElement("link");
        link.rel = "icon";
        document.head.appendChild(link);
      }
      link.href = faviconHref;
    }
  }, [title, faviconHref]);
}

export default usePageMeta;