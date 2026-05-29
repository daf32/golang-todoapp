const BASE = '/api/v1';

function getToken() { return localStorage.getItem('access_token'); }
function getRefresh() { return localStorage.getItem('refresh_token'); }

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
    throw Object.assign(new Error(err.error || err.message || res.statusText), { status: res.status });
  }
  if (res.status === 204) return null;
  return res.json();
}

export const api = {
  /* Auth */
  login: (email, password) => request('POST', '/auth/login', { email, password }, false),
  register: (full_name, email, password) => request('POST', '/auth/register', { full_name, email, password }, false),
  logout: () => request('POST', '/auth/logout', { refresh_token: getRefresh() }),
  resendConfirmation: (email) => request('POST', '/auth/resend-confirmation', { email }, false),

  /* Users */
  getUser: (id) => request('GET', `/users/${id}`),
  patchUser: (id, patch) => request('PATCH', `/users/${id}`, patch),
  changePassword: (id, current_password, new_password) =>
    request('POST', `/users/${id}/change-password`, { current_password, new_password }),
  deleteUser: (id) => request('DELETE', `/users/${id}`),

  /* Tasks */
  getTasks: (params = {}) => {
    const q = new URLSearchParams();
    if (params.limit != null) q.set('limit', params.limit);
    if (params.offset != null) q.set('offset', params.offset);
    const qs = q.toString();
    return request('GET', `/tasks${qs ? '?' + qs : ''}`);
  },
  createTask: (data) => request('POST', '/tasks', data),
  patchTask: (id, patch) => request('PATCH', `/tasks/${id}`, patch),
  deleteTask: (id) => request('DELETE', `/tasks/${id}`),

  /* Statistics */
  getStatistics: (params = {}) => {
    const q = new URLSearchParams();
    if (params.from_days != null) q.set('from_days', params.from_days);
    if (params.to_days != null) q.set('to_days', params.to_days);
    const qs = q.toString();
    return request('GET', `/statistics${qs ? '?' + qs : ''}`);
  },
};
