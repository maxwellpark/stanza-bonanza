import { useEffect } from 'react';

function setMeta(selector: string, attr: string, value: string) {
  var el = document.head.querySelector<HTMLMetaElement>(selector);
  if (!el) {
    el = document.createElement('meta');
    var [name, key] = attr.split('=');
    el.setAttribute(name, key.replace(/"/g, ''));
    document.head.appendChild(el);
  }
  el.setAttribute('content', value);
}

// Sets the document title + description/OG meta for the current page. Helps the
// browser tab and JS-running unfurlers; crawler cards still need SSR/an edge
// function. Restores the site defaults on unmount.
export function usePageMeta(title: string, description?: string) {
  useEffect(() => {
    var previousTitle = document.title;
    document.title = title;

    if (description) {
      setMeta('meta[name="description"]', 'name="description"', description);
      setMeta('meta[property="og:description"]', 'property="og:description"', description);
      setMeta('meta[name="twitter:description"]', 'name="twitter:description"', description);
    }
    setMeta('meta[property="og:title"]', 'property="og:title"', title);
    setMeta('meta[name="twitter:title"]', 'name="twitter:title"', title);

    return () => {
      document.title = previousTitle;
    };
  }, [title, description]);
}
