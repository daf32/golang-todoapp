import React, { useState } from 'react';
import { useNav, useAuth } from '../App';
import { api } from '../api/client';
import { CheckIcon, MailIcon, EyeIcon, EyeOffIcon } from '../components/Icons';

function AuthDeco() {
  return (
    <div className="auth-deco">
      <div className="auth-deco-logo">
        <CheckIcon size={28} />
      </div>
      <h1 className="auth-deco-title">Welcome<br />back!</h1>
      <p className="auth-deco-sub">Pick up right where you left off. Your tasks are waiting.</p>
      <div className="auth-deco-dots">
        <span /><span /><span />
      </div>
    </div>
  );
}

export default function LoginPage() {
  const { navigate }     = useNav();
  const { login }        = useAuth();
  const [form, setForm]  = useState({ email: '', password: '' });
  const [showPw, setShowPw] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError]   = useState('');

  const set = (field) => (e) => setForm((f) => ({ ...f, [field]: e.target.value }));

  async function handleSubmit(e) {
    e.preventDefault();
    setLoading(true); setError('');
    try {
      const { access_token, refresh_token } = await api.login(form.email, form.password);
      login(access_token, refresh_token);
    } catch (err) {
      setError(err.message || 'Invalid email or password');
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="auth-page-shell">
      <AuthDeco />

      <div className="auth-form-panel">
        <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', marginBottom: 28 }}>
          <div className="logo-mark" style={{ marginBottom: 16 }}>
            <CheckIcon size={22} />
          </div>
          <h1 style={{ fontSize: 26, fontWeight: 800, textAlign: 'center', color: '#1E1B2E' }}>
            Welcome back!
          </h1>
          <p style={{ fontSize: 14, color: '#A9A9C0', marginTop: 6, textAlign: 'center' }}>
            Sign in to your account
          </p>
        </div>

        {error && (
          <div className="error-msg" style={{ marginBottom: 12, marginLeft: 0, marginRight: 0 }}>
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: 14 }}>
          <div className="input-group">
            <input
              className="input"
              type="email"
              placeholder="Email address"
              value={form.email}
              onChange={set('email')}
              autoComplete="email"
              required
            />
            <span className="icon-right"><MailIcon size={16} /></span>
          </div>

          <div className="input-group">
            <input
              className="input"
              type={showPw ? 'text' : 'password'}
              placeholder="Password"
              value={form.password}
              onChange={set('password')}
              autoComplete="current-password"
              required
            />
            <button
              type="button"
              className="icon-right clickable"
              onClick={() => setShowPw((v) => !v)}
              tabIndex={-1}
            >
              {showPw ? <EyeOffIcon size={16} /> : <EyeIcon size={16} />}
            </button>
          </div>

          <button
            type="submit"
            className="btn btn-primary"
            disabled={loading}
            style={{ marginTop: 6, width: '100%' }}
          >
            {loading ? 'Logging in…' : 'Log in'}
          </button>
        </form>

        <div className="divider" style={{ margin: '24px 0' }}>or log in with</div>

        <div className="social-row">
          <button className="social-btn" aria-label="Log in with Google"
            style={{ width: 'auto', borderRadius: 13, padding: '0 20px', gap: 10, border: '1.5px solid #EDEDF5' }}>
            <svg width="18" height="18" viewBox="0 0 48 48">
              <path fill="#EA4335" d="M24 9.5c3.54 0 6.71 1.22 9.21 3.6l6.85-6.85C35.9 2.38 30.47 0 24 0 14.62 0 6.51 5.38 2.56 13.22l7.98 6.19C12.43 13.72 17.74 9.5 24 9.5z"/>
              <path fill="#4285F4" d="M46.98 24.55c0-1.57-.15-3.09-.38-4.55H24v9.02h12.94c-.58 2.96-2.26 5.48-4.78 7.18l7.73 6c4.51-4.18 7.09-10.36 7.09-17.65z"/>
              <path fill="#FBBC05" d="M10.53 28.59c-.48-1.45-.76-2.99-.76-4.59s.27-3.14.76-4.59l-7.98-6.19C.92 16.46 0 20.12 0 24c0 3.88.92 7.54 2.56 10.78l7.97-6.19z"/>
              <path fill="#34A853" d="M24 48c6.48 0 11.93-2.13 15.89-5.81l-7.73-6c-2.18 1.48-4.97 2.31-8.16 2.31-6.26 0-11.57-4.22-13.47-9.91l-7.98 6.19C6.51 42.62 14.62 48 24 48z"/>
            </svg>
            <span style={{ fontSize: 14, fontWeight: 600, color: '#1E1B2E' }}>Continue with Google</span>
          </button>
        </div>

        <p style={{ textAlign: 'center', fontSize: 13, color: '#A9A9C0', marginTop: 28 }}>
          Don't have an account?{' '}
          <button
            onClick={() => navigate('signup')}
            style={{
              color: 'var(--primary-btn)', fontWeight: 700,
              background: 'none', border: 'none', cursor: 'pointer', fontSize: 13,
            }}
          >
            Get started
          </button>
        </p>
      </div>
    </div>
  );
}
