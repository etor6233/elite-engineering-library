export function parseHelpQuery(url: string) {
  if (url.length > 2048) return null;
  const p = new URL(url).searchParams;
  if ([...p.keys()].some(key => !["q", "article", "version"].includes(key) || p.getAll(key).length !== 1)) return null;
  const q = p.get("q") ?? "", article = p.get("article"), version = p.get("version");
  if (q.length > 160 || /[\u0000-\u001f\u007f]/.test(q)) return null;
  if ((article === null) !== (version === null) || (article !== null && (p.has("q") || !/^[a-z][a-z0-9-]{0,63}$/.test(article) || !/^\d{1,6}\.\d{1,6}\.\d{1,6}$/.test(version!)))) return null;
  return { q: q.trim().normalize("NFC").toLocaleLowerCase("es"), article, version };
}
