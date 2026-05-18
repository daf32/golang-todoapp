import React, { useState, useEffect } from 'react';
import { useAuth } from '../App';
import { api } from '../api/client';
import { getUserIdFromToken } from '../utils/jwt';
import AuthLayout from '../components/AuthLayout';

export default function ProfilePage() {
  const { token, logout }       = useAuth();
  const [user,    setUser]      = useState(null);
  const [tasks,   setTasks]     = useState([]);
  const [loading, setLoading]   = useState(true);
  const [error,   setError]     = useState('');
  const [editing, setEditing]   = useState(false);
  const [form,    setForm]      = useState({ full_name: '', phone_number: '' });
  const [saving,  setSaving]    = useState(false);
  const [saved,   setSaved]     = useState(false);

  const userId = getUserIdFromToken(token);

  useEffect(() => {
    if (!userId) { logout(); return; }
    (async () => {
      setLoading(true);
      try {
        const [u, t] = await Promise.all([
          api.getUser(userId),
          api.getTasks().catch(() => []),
        ]);
        setUser(u);
        setForm({ full_name: u.full_name ?? '', phone_number: u.phone_number ?? '' });
        setTasks(Array.isArray(t) ? t : []);
      } catch (err) {
        if (err.status === 401) { logout(); return; }
        setError(err.message || 'Failed to load profile');
      } finally {
        setLoading(false);
      }
    })();
  }, [userId, logout]);

  async function handleSave(e) {
    e.preventDefault();
    setSaving(true);
    try {
      const patch = {};
      if (form.full_name   !== user.full_name)   patch.full_name   = form.full_name;
      if (form.phone_number !== (user.phone_number ?? '')) {
        patch.phone_number = form.phone_number || null;
      }
      const updated = await api.patchUser(userId, patch);
      setUser(updated);
      setEditing(false);
      setSaved(true);
      setTimeout(() => setSaved(false), 2000);
    } catch (err) {
      setError(err.message || 'Failed to save');
    } finally {
      setSaving(false);
    }
  }

  const initials = user?.full_name
    ? user.full_name.split(' ').map((w) => w[0]).join('').slice(0, 2).toUpperCase()
    : '?';

  const doneTasks  = tasks.filter((t) => t.completed).length;
  const totalTasks = tasks.length;
  const role       = user?.role ?? 'user';

  return (
    <AuthLayout title="Profile" date="Your account">
      {loading ? (
        <div className="spinner" style={{ marginTop: 60 }} />
      ) : (
        <>
          {error && <div className="error-msg">{error}</div>}

          {/* Hero */}
          <div className="profile-hero">
            <div className="profile-avatar">{initials}</div>
            <div className="profile-name">{user?.full_name ?? '—'}</div>
            <div className="profile-role">{role}</div>
          </div>

          {/* Stats */}
          <div style={{ display: 'flex', gap: 12, padding: '32px 20px 4px' }}>
            <div className="stat-card">
              <div className="stat-number">{totalTasks}</div>
              <div className="stat-label">Total</div>
            </div>
            <div className="stat-card" style={{ animationDelay: '60ms' }}>
              <div className="stat-number">{doneTasks}</div>
              <div className="stat-label">Done</div>
            </div>
            <div className="stat-card" style={{ animationDelay: '120ms' }}>
              <div className="stat-number">{totalTasks - doneTasks}</div>
              <div className="stat-label">Active</div>
            </div>
          </div>

          {/* Info */}
          <div className="settings-section" style={{ margin: '16px 20px 0', animationDelay: '100ms' }}>
            <div className="settings-section-label">Account info</div>

            {!editing ? (
              <>
                <InfoRow label="Full name"    value={user?.full_name ?? '—'} />
                <InfoRow label="Email"        value={user?.email ?? '—'} />
                <InfoRow label="Phone"        value={user?.phone_number ?? '—'} />
                <div style={{ padding: '12px 16px' }}>
                  <button
                    className="btn btn-primary"
                    style={{ width: '100%' }}
                    onClick={() => setEditing(true)}
                  >
                    Edit Profile
                  </button>
                </div>
              </>
            ) : (
              <form onSubmit={handleSave} style={{ padding: '12px 16px', display: 'flex', flexDirection: 'column', gap: 12 }}>
                <div>
                  <label style={{ fontSize: 12, fontWeight: 600, color: 'var(--text-3)', display: 'block', marginBottom: 6 }}>
                    Full name
                  </label>
                  <input
                    className="input input-no-icon"
                    value={form.full_name}
                    onChange={(e) => setForm((f) => ({ ...f, full_name: e.target.value }))}
                    minLength={3} maxLength={100} required
                  />
                </div>
                <div>
                  <label style={{ fontSize: 12, fontWeight: 600, color: 'var(--text-3)', display: 'block', marginBottom: 6 }}>
                    Phone (optional)
                  </label>
                  <input
                    className="input input-no-icon"
                    placeholder="+1234567890"
                    value={form.phone_number}
                    onChange={(e) => setForm((f) => ({ ...f, phone_number: e.target.value }))}
                    maxLength={15}
                  />
                </div>
                <div style={{ display: 'flex', gap: 10 }}>
                  <button type="button" className="btn btn-outline" style={{ flex: 1 }} onClick={() => setEditing(false)}>
                    Cancel
                  </button>
                  <button type="submit" className="btn btn-primary" style={{ flex: 1 }} disabled={saving}>
                    {saving ? 'Saving…' : saved ? '✓ Saved' : 'Save'}
                  </button>
                </div>
              </form>
            )}
          </div>

          <div style={{ height: 32 }} />
        </>
      )}
    </AuthLayout>
  );
}

function InfoRow({ label, value }) {
  return (
    <div className="settings-row" style={{ cursor: 'default' }}>
      <span className="settings-row-text">
        <span className="settings-row-sub">{label}</span>
        <div className="settings-row-title">{value}</div>
      </span>
    </div>
  );
}
