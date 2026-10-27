export function parseRoute(hash = location.hash) {
  const value = hash.replace(/^#/, '') || '/';
  const [path, query = ''] = value.split('?');
  return { path, params: new URLSearchParams(query) };
}

export function routeLink(path, params = {}) {
  const query = new URLSearchParams();
  for (const [key, value] of Object.entries(params)) {
    for (const item of Array.isArray(value) ? value : [value]) {
      if (item !== '' && item !== undefined && item !== null) query.append(key, String(item));
    }
  }
  return `#${path}${query.size ? `?${query}` : ''}`;
}

export function navigate(path) {
  location.hash = path;
}
