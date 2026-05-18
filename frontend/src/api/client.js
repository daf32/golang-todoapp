const BASE = '/api/v1';

function getToken() {
  return localStorage.getItem('access_token');
}

async function request(method, path, body, auth = true) {
  const headers = { 'Content-Type': 'application/json' };
  if (auth) {
    const token = getToken();
    if (token) headers['Authorization'] = `Bearer ${token}`;
  }

  const res = await fetch(`${BASE}${path}`, {
    method,
    headers,
    body: body !== undefined ? JSON.stringify(body) : undefined,
  });

  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }));
    throw Object.assign(
      new Error(err.error || err.message || res.statusText),
      { status: res.status }
    );
  }

  if (res.status === 204) return null;
  return res.json();
}

export const api = {
  /* ── Auth ── */
  login: (email, password) =>
    request('POST', '/auth/login', { email, password }, false),

  register: (full_name, email, password) =>
    request('POST', '/auth/register', { full_name, email, password }, false),

  logout: (refreshToken) =>
    request('POST', '/auth/logout', { refresh_token: refreshToken }),

  /* ── Users ── */
  getUser: (id) =>
    request('GET', `/users/${id}`),

  patchUser: (id, patch) =>
    request('PATCH', `/users/${id}`, patch),

  /* ── Tasks ── */
  getTasks: (params = {}) => {
    const q = new URLSearchParams();
    if (params.limit  != null) q.set('limit',  params.limit);
    if (params.offset != null) q.set('offset', params.offset);
    const qs = q.toString();
    return request('GET', `/tasks${qs ? '?' + qs : ''}`);
  },

  createTask: (title, description) =>
    request('POST', '/tasks', { title, ...(description ? { description } : {}) }),

  patchTask: (id, patch) =>
    request('PATCH', `/tasks/${id}`, patch),

  deleteTask: (id) =>
    request('DELETE', `/tasks/${id}`),
};
