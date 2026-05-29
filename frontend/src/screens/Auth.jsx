import { useState } from 'react';
import { CheckSquare, Lock, Mail, User, Eye, EyeOff, X, Check, BadgeCheck, GoogleMark } from '../components/Icons.jsx';
import { api } from '../api/client.js';

export function PasswordField({ value, onChange, placeholder = 'Password', autoComplete }) {
  const [show, setShow] = useState(false);
  return (
    <div className="input-wrap">
      <span className="icon-lead"><Lock size={17} /></span>
      <input className="input" type={show ? 'text' : 'password'} value={value} placeholder={placeholder}
        autoComplete={autoComplete} onChange={(e) => onChange(e.target.value)} />
      <button type="button" className="icon-btn" style={{ position: 'absolute', right: 6 }}
        onClick={() => setShow((s) => !s)} aria-label={show ? 'Hide password' : 'Show password'}>
        {show ? <EyeOff size={17} /> : <Eye size={17} />}
      </button>
    </div>
  );
}

export function Login({ go, onAuth }) {
  const [email, setEmail] = useState('');
  const [pw, setPw] = useState('');
  const [busy, setBusy] = useState(false);
  const [err, setErr] = useState('');
  async function submit(e) {
    e.preventDefault();
    setErr('');
    if (!email || !pw) { setErr('Enter your email and password'); return; }
    setBusy(true);
    try {
      const res = await api.login(email, pw);
      await onAuth(res.access_token, res.refresh_token);
    } catch (ex) {
      setErr(ex.message || 'Sign-in failed');
    } finally {
      setBusy(false);
    }
  }
  return (
    <div className="auth-card">
      <div className="auth-brand">
        <div className="auth-logo"><CheckSquare size={28} /></div>
        <div>
          <div className="auth-title">Welcome back</div>
          <div className="auth-sub">Sign in to pick up where you left off.</div>
        </div>
      </div>
      <form className="auth-form" onSubmit={submit}>
        {err && <div className="banner err"><X size={15} />{err}</div>}
        <div className="field">
          <label className="field-label">Email</label>
          <div className="input-wrap">
            <span className="icon-lead"><Mail size={17} /></span>
            <input className="input" type="email" value={email} placeholder="you@example.com"
              autoComplete="email" onChange={(e) => setEmail(e.target.value)} />
          </div>
        </div>
        <div className="field">
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'baseline' }}>
            <label className="field-label">Password</label>
            <span className="link" style={{ fontSize: 13 }}>Forgot?</span>
          </div>
          <PasswordField value={pw} onChange={setPw} autoComplete="current-password" />
        </div>
        <button className="btn btn-primary btn-block" type="submit" disabled={busy} style={{ marginTop: 4 }}>
          {busy ? <span className="spinner" style={{ width: 18, height: 18, borderTopColor: 'var(--on-accent)', borderColor: 'rgba(255,255,255,0.3)' }} /> : 'Sign in'}
        </button>
        <div className="auth-divider">or</div>
        <button type="button" className="btn btn-google btn-block" onClick={() => { window.location.href = '/api/v1/auth/oauth/google'; }}>
          <GoogleMark /> Continue with Google
        </button>
      </form>
      <div className="auth-foot">New here? <a onClick={() => go('signup')}>Create an account</a></div>
    </div>
  );
}

export function SignUp({ go, onPending }) {
  const [form, setForm] = useState({ full_name: '', email: '', password: '' });
  const [busy, setBusy] = useState(false);
  const [errs, setErrs] = useState({});
  const [topErr, setTopErr] = useState('');
  function set(k, v) { setForm((f) => ({ ...f, [k]: v })); setErrs((e) => ({ ...e, [k]: '' })); setTopErr(''); }
  async function submit(e) {
    e.preventDefault();
    const next = {};
    if (form.full_name.trim().length < 3) next.full_name = 'Use at least 3 characters';
    if (!/^[^@\s]+@[^@\s]+\.[^@\s]{2,}$/.test(form.email)) next.email = 'Enter a valid email';
    if (form.password.length < 8) next.password = 'Use at least 8 characters';
    setErrs(next);
    if (Object.keys(next).length) return;
    setBusy(true);
    try {
      await api.register(form.full_name.trim(), form.email, form.password);
      onPending(form.email);
    } catch (ex) {
      setTopErr(ex.message || 'Sign-up failed');
    } finally {
      setBusy(false);
    }
  }
  return (
    <div className="auth-card">
      <div className="auth-brand">
        <div className="auth-logo"><CheckSquare size={28} /></div>
        <div>
          <div className="auth-title">Create your account</div>
          <div className="auth-sub">Start organizing your day in a minute.</div>
        </div>
      </div>
      <form className="auth-form" onSubmit={submit}>
        {topErr && <div className="banner err"><X size={15} />{topErr}</div>}
        <div className="field">
          <label className="field-label">Full name</label>
          <div className="input-wrap">
            <span className="icon-lead"><User size={17} /></span>
            <input className="input" value={form.full_name} placeholder="Jane Doe"
              maxLength={100} onChange={(e) => set('full_name', e.target.value)} />
          </div>
          {errs.full_name && <span className="field-err">{errs.full_name}</span>}
        </div>
        <div className="field">
          <label className="field-label">Email</label>
          <div className="input-wrap">
            <span className="icon-lead"><Mail size={17} /></span>
            <input className="input" type="email" value={form.email} placeholder="you@example.com"
              onChange={(e) => set('email', e.target.value)} />
          </div>
          {errs.email && <span className="field-err">{errs.email}</span>}
        </div>
        <div className="field">
          <label className="field-label">Password</label>
          <PasswordField value={form.password} onChange={(v) => set('password', v)} placeholder="At least 8 characters" autoComplete="new-password" />
          {errs.password && <span className="field-err">{errs.password}</span>}
        </div>
        <button className="btn btn-primary btn-block" type="submit" disabled={busy} style={{ marginTop: 4 }}>
          {busy ? <span className="spinner" style={{ width: 18, height: 18, borderTopColor: 'var(--on-accent)', borderColor: 'rgba(255,255,255,0.3)' }} /> : 'Create account'}
        </button>
        <div className="auth-divider">or</div>
        <button type="button" className="btn btn-google btn-block" onClick={() => { window.location.href = '/api/v1/auth/oauth/google'; }}>
          <GoogleMark /> Continue with Google
        </button>
      </form>
      <div className="auth-foot">Already have an account? <a onClick={() => go('login')}>Sign in</a></div>
    </div>
  );
}

export function EmailConfirm({ email, go }) {
  const [resent, setResent] = useState(false);
  const [resending, setResending] = useState(false);
  async function resend() {
    setResending(true);
    try { await api.resendConfirmation(email); } catch (_) {}
    setResending(false);
    setResent(true);
    setTimeout(() => setResent(false), 2400);
  }
  return (
    <div className="auth-card" style={{ textAlign: 'center' }}>
      <div className="confirm-icon"><Mail size={34} /></div>
      <div className="auth-title">Check your inbox</div>
      <div className="auth-sub" style={{ margin: '6px auto 0' }}>
        We sent a confirmation link to <strong style={{ color: 'var(--text)' }}>{email}</strong>. Click it to verify your email, then sign in.
      </div>
      {resent && <div className="banner ok" style={{ marginTop: 20, justifyContent: 'center' }}><Check size={15} />Confirmation resent</div>}
      <div style={{ display: 'flex', flexDirection: 'column', gap: 10, marginTop: 26 }}>
        <button className="btn btn-primary btn-block" onClick={() => go('login')}>Back to sign in</button>
        <button className="btn btn-pearl btn-block" onClick={resend} disabled={resending}>{resending ? 'Sending…' : 'Resend email'}</button>
      </div>
      <div className="auth-foot">Wrong address? <a onClick={() => go('signup')}>Go back</a></div>
    </div>
  );
}

export function EmailConfirmed({ onAuth }) {
  return (
    <div className="auth-card" style={{ textAlign: 'center' }}>
      <div className="confirm-icon ok"><BadgeCheck size={36} /></div>
      <div className="auth-title">Email confirmed</div>
      <div className="auth-sub" style={{ margin: '6px auto 0' }}>You&#39;re all set. Let&#39;s get to work.</div>
      <button className="btn btn-primary btn-block" style={{ marginTop: 26 }} onClick={onAuth}>Open my tasks</button>
    </div>
  );
}
