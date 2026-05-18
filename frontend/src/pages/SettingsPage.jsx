import React from 'react';
import { useAuth, useNav } from '../App';
import { useTheme } from '../context/ThemeContext';
import AuthLayout from '../components/AuthLayout';
import {
  SettingsIcon, UserIcon, ChevronRightIcon,
} from '../components/Icons';

const SunIcon = () => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <circle cx="12" cy="12" r="5"/>
    <line x1="12" y1="1" x2="12" y2="3"/>
    <line x1="12" y1="21" x2="12" y2="23"/>
    <line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/>
    <line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/>
    <line x1="1" y1="12" x2="3" y2="12"/>
    <line x1="21" y1="12" x2="23" y2="12"/>
    <line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/>
    <line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/>
  </svg>
);
const MoonIcon = () => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <path d="M21 12.79A9 9 0 1111.21 3 7 7 0 0021 12.79z"/>
  </svg>
);
const MonitorIcon = () => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <rect x="2" y="3" width="20" height="14" rx="2"/>
    <line x1="8" y1="21" x2="16" y2="21"/>
    <line x1="12" y1="17" x2="12" y2="21"/>
  </svg>
);
const LogOutIcon = () => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <path d="M9 21H5a2 2 0 01-2-2V5a2 2 0 012-2h4"/>
    <polyline points="16 17 21 12 16 7"/>
    <line x1="21" y1="12" x2="9" y2="12"/>
  </svg>
);

const THEMES = [
  { key: 'light',  label: 'Light',  icon: <SunIcon /> },
  { key: 'dark',   label: 'Dark',   icon: <MoonIcon /> },
  { key: 'system', label: 'System', icon: <MonitorIcon /> },
];

export default function SettingsPage() {
  const { logout }        = useAuth();
  const { navigate }      = useNav();
  const { theme, setTheme } = useTheme();

  return (
    <AuthLayout title="Settings" date="Preferences">
      {/* Appearance */}
      <div className="settings-section" style={{ animationDelay: '0ms' }}>
        <div className="settings-section-label">Appearance</div>
        <div className="theme-options">
          {THEMES.map(({ key, label, icon }) => (
            <button
              key={key}
              className={`theme-option ${theme === key ? 'active' : ''}`}
              onClick={() => setTheme(key)}
            >
              {icon}
              {label}
            </button>
          ))}
        </div>
      </div>

      {/* Account */}
      <div className="settings-section" style={{ animationDelay: '60ms' }}>
        <div className="settings-section-label">Account</div>

        <button
          className="settings-row"
          style={{ width: '100%', textAlign: 'left' }}
          onClick={() => navigate('profile')}
        >
          <span className="settings-row-icon"><UserIcon size={18} /></span>
          <span className="settings-row-text">
            <span className="settings-row-title">My Profile</span>
            <div className="settings-row-sub">View and edit your account info</div>
          </span>
          <span className="settings-row-right"><ChevronRightIcon size={16} /></span>
        </button>
      </div>

      {/* Danger zone */}
      <div className="settings-section" style={{ animationDelay: '120ms' }}>
        <div className="settings-section-label">Session</div>
        <button
          className="settings-row"
          style={{ width: '100%', textAlign: 'left' }}
          onClick={logout}
        >
          <span className="settings-row-icon" style={{ background: 'var(--tag-high-bg)', color: 'var(--tag-high)' }}>
            <LogOutIcon />
          </span>
          <span className="settings-row-text">
            <span className="settings-row-title" style={{ color: 'var(--tag-high)' }}>Log out</span>
            <div className="settings-row-sub">Sign out of your account</div>
          </span>
        </button>
      </div>

      <div style={{ height: 32 }} />
    </AuthLayout>
  );
}
