import { useState, useEffect } from 'react';
import { Check, Circle, Inbox, BadgeCheck, Pencil, Sun, Moon, Monitor, Key, Bell, LogOut, Trash, ChevronRight } from '../components/Icons.jsx';
import { StatCard } from './Tasks.jsx';
import { PasswordField } from './Auth.jsx';
import { initials, fmtDayMonth } from '../lib/format.js';

export function ProfileView({ user, tasks, onSave }) {
  const [editing, setEditing] = useState(false);
  const [form, setForm] = useState({ full_name: user.full_name || '', phone_number: user.phone_number || '' });
  const [saved, setSaved] = useState(false);
  const [err, setErr] = useState('');
  const [busy, setBusy] = useState(false);
  const done = tasks.filter((t) => t.completed).length;

  useEffect(() => { setForm({ full_name: user.full_name || '', phone_number: user.phone_number || '' }); }, [user]);

  async function save(e) {
    e.preventDefault();
    setErr('');
    const name = form.full_name.trim();
    const phone = form.phone_number.trim();
    if (name.length < 3 || name.length > 100) { setErr('Name must be between 3 and 100 characters'); return; }
    if (phone && (!phone.startsWith('+') || phone.length < 10 || phone.length > 15)) {
      setErr('Phone must start with "+" and be 10–15 characters');
      return;
    }
    setBusy(true);
    try {
      await onSave({ full_name: name, phone_number: phone || null });
      setEditing(false); setSaved(true); setTimeout(() => setSaved(false), 2000);
    } catch (ex) {
      setErr(ex.message || 'Failed to save');
    } finally {
      setBusy(false);
    }
  }
  return (
    <div>
      <div className="profile-hero">
        <div className="profile-avatar">{initials(user.full_name)}</div>
        <div style={{ minWidth: 0 }}>
          <div className="profile-name">{user.full_name}</div>
          <div className="profile-meta">
            <span className="pill outline" style={{ textTransform: 'capitalize' }}>{user.role}</span>
            {user.email_verified && <span className="pill" style={{ color: 'var(--green)' }}><BadgeCheck size={14} />Verified</span>}
          </div>
        </div>
      </div>

      <div className="stat-grid" style={{ gridTemplateColumns: 'repeat(3, 1fr)', marginTop: 28 }}>
        <StatCard icon={<Inbox size={18} />} value={tasks.length} label="Total tasks" delay="0ms" />
        <StatCard icon={<Check size={18} stroke={2.4} />} value={done} label="Completed" delay="60ms" />
        <StatCard icon={<Circle size={18} />} value={tasks.length - done} label="Active" delay="120ms" />
      </div>

      <div className="settings-block">
        <div className="settings-block-label section-label">Account details</div>
        <div className="card info-card" style={{ marginTop: 0 }}>
          {!editing ? (
            <>
              <div className="info-row"><span className="info-key">Full name</span><span className="info-val">{user.full_name}</span></div>
              <div className="info-row"><span className="info-key">Email</span><span className="info-val">{user.email}</span></div>
              <div className="info-row"><span className="info-key">Phone</span><span className="info-val">{user.phone_number || '—'}</span></div>
              {user.created_at && <div className="info-row"><span className="info-key">Member since</span><span className="info-val">{fmtDayMonth(user.created_at)} {new Date(user.created_at).getFullYear()}</span></div>}
              <div style={{ padding: 16 }}>
                <button className="btn btn-pearl btn-block" onClick={() => setEditing(true)}><Pencil size={16} />{saved ? 'Saved ✓' : 'Edit profile'}</button>
              </div>
            </>
          ) : (
            <form onSubmit={save} style={{ padding: 20, display: 'flex', flexDirection: 'column', gap: 16 }}>
              <div className="field">
                <label className="field-label">Full name</label>
                <input className="input" value={form.full_name} minLength={3} maxLength={100} required
                  onChange={(e) => setForm((f) => ({ ...f, full_name: e.target.value }))} />
              </div>
              <div className="field">
                <label className="field-label">Phone number</label>
                <input className="input" value={form.phone_number} placeholder="+1 555 555 1234"
                  onChange={(e) => setForm((f) => ({ ...f, phone_number: e.target.value }))} />
              </div>
              <div className="field">
                <label className="field-label">Email</label>
                <input className="input" value={user.email} disabled style={{ opacity: 0.6 }} />
                <span className="t-fine fg-4">Email can&#39;t be changed here.</span>
              </div>
              {err && <div className="field-err">{err}</div>}
              <div style={{ display: 'flex', gap: 10 }}>
                <button type="button" className="btn btn-pearl btn-block" onClick={() => { setEditing(false); setErr(''); }} disabled={busy}>Cancel</button>
                <button type="submit" className="btn btn-primary btn-block" disabled={busy}>{busy ? 'Saving…' : 'Save changes'}</button>
              </div>
            </form>
          )}
        </div>
      </div>
    </div>
  );
}

const THEMES = [
  { value: 'light', name: 'Light', icon: <Sun size={20} /> },
  { value: 'dark', name: 'Dark', icon: <Moon size={20} /> },
  { value: 'system', name: 'System', icon: <Monitor size={20} /> },
];

export function SettingsView({ user, theme, setTheme, onChangePassword, onLogout, onDelete }) {
  return (
    <div>
      <div className="settings-block" style={{ marginTop: 0 }}>
        <div className="settings-block-label section-label">Appearance</div>
        <div className="card settings-card">
          <div className="theme-picker">
            {THEMES.map((t) => (
              <button key={t.value} className={`theme-opt ${theme === t.value ? 'active' : ''}`} onClick={() => setTheme(t.value)}>
                {t.icon}<span className="theme-opt-name">{t.name}</span>
              </button>
            ))}
          </div>
        </div>
      </div>

      <div className="settings-block">
        <div className="settings-block-label section-label">Account</div>
        <div className="card settings-card">
          <button className="settings-row" onClick={onChangePassword}>
            <span className="settings-row-icon"><Key size={18} /></span>
            <span className="settings-row-text">
              <div className="settings-row-title">Change password</div>
              <div className="settings-row-sub">Update your sign-in credentials</div>
            </span>
            <ChevronRight size={18} className="fg-4" />
          </button>
          <div className="settings-row">
            <span className="settings-row-icon"><Bell size={18} /></span>
            <span className="settings-row-text">
              <div className="settings-row-title">Notifications</div>
              <div className="settings-row-sub">Reminders for upcoming tasks</div>
            </span>
            <Toggle />
          </div>
        </div>
      </div>

      <div className="settings-block">
        <div className="settings-block-label section-label">Session</div>
        <div className="card settings-card">
          <button className="settings-row" onClick={onLogout}>
            <span className="settings-row-icon"><LogOut size={18} /></span>
            <span className="settings-row-text"><div className="settings-row-title">Log out</div></span>
            <ChevronRight size={18} className="fg-4" />
          </button>
          <button className="settings-row" onClick={onDelete}>
            <span className="settings-row-icon" style={{ color: 'var(--red)' }}><Trash size={18} /></span>
            <span className="settings-row-text"><div className="settings-row-title" style={{ color: 'var(--red)' }}>Delete account</div><div className="settings-row-sub">Permanently remove your data</div></span>
            <ChevronRight size={18} className="fg-4" />
          </button>
        </div>
      </div>

      <div style={{ textAlign: 'center', marginTop: 32 }} className="t-fine fg-4">
        todo · signed in as {user.email}
      </div>
    </div>
  );
}

export function Toggle({ initial = true }) {
  const [on, setOn] = useState(initial);
  return (
    <button onClick={() => setOn((o) => !o)} aria-label="Toggle"
      style={{ width: 46, height: 28, borderRadius: 99, background: on ? 'var(--accent)' : 'var(--surface-muted)', padding: 3, transition: 'background 160ms var(--ease)', flexShrink: 0 }}>
      <span style={{ display: 'block', width: 22, height: 22, borderRadius: 99, background: on ? 'var(--on-accent)' : '#fff', boxShadow: '0 1px 3px rgba(0,0,0,0.3)', transform: on ? 'translateX(18px)' : 'translateX(0)', transition: 'transform 160ms var(--ease)' }} />
    </button>
  );
}

export function ChangePasswordForm({ onClose, onSubmit }) {
  const [cur, setCur] = useState('');
  const [next, setNext] = useState('');
  const [done, setDone] = useState(false);
  const [err, setErr] = useState('');
  const [busy, setBusy] = useState(false);
  async function submit(e) {
    e.preventDefault();
    if (next.length < 8) { setErr('New password must be at least 8 characters'); return; }
    setBusy(true);
    try {
      if (onSubmit) await onSubmit(cur, next);
      setDone(true);
      setTimeout(onClose, 1100);
    } catch (ex) {
      setErr(ex.message || 'Failed to update');
      setBusy(false);
    }
  }
  if (done) return <div className="banner ok" style={{ justifyContent: 'center' }}><Check size={15} />Password updated</div>;
  return (
    <form className="modal-body" onSubmit={submit}>
      <div className="field">
        <label className="field-label">Current password</label>
        <PasswordField value={cur} onChange={setCur} placeholder="Current password" />
      </div>
      <div className="field">
        <label className="field-label">New password</label>
        <PasswordField value={next} onChange={(v) => { setNext(v); setErr(''); }} placeholder="At least 8 characters" />
        {err && <span className="field-err">{err}</span>}
      </div>
      <div className="modal-foot">
        <button type="button" className="btn btn-pearl btn-block" onClick={onClose} disabled={busy}>Cancel</button>
        <button type="submit" className="btn btn-primary btn-block" disabled={busy}>{busy ? 'Updating…' : 'Update'}</button>
      </div>
    </form>
  );
}
