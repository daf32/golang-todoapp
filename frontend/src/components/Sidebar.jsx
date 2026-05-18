import React from 'react';
import { useNav, useAuth } from '../App';
import { ListIcon, CalendarIcon, UserIcon, SettingsIcon, CheckIcon } from './Icons';
import { parseJwt } from '../utils/jwt';

const NAV = [
  { key: 'tasks',    label: 'My Tasks',  Icon: ListIcon },
  { key: 'calendar', label: 'Calendar',  Icon: CalendarIcon },
  { key: 'profile',  label: 'Profile',   Icon: UserIcon },
  { key: 'settings', label: 'Settings',  Icon: SettingsIcon },
];

export default function Sidebar() {
  const { screen, navigate } = useNav();
  const { token }            = useAuth();

  const claims   = parseJwt(token);
  const initials = (claims?.full_name ?? claims?.name ?? 'U')
    .split(' ').map((w) => w[0]).join('').slice(0, 2).toUpperCase();

  return (
    <aside className="sidebar">
      {/* Logo — click goes to main page */}
      <div className="sidebar-logo">
        <div className="logo-mark">
          <CheckIcon size={22} />
        </div>
        <button
          className="sidebar-logo-name"
          onClick={() => navigate('tasks')}
        >
          Taskly
        </button>
      </div>

      {/* Nav items */}
      <nav className="sidebar-nav">
        {NAV.map(({ key, label, Icon }) => (
          <button
            key={key}
            className={`sidebar-item ${screen === key ? 'active' : ''}`}
            onClick={() => navigate(key)}
          >
            <Icon size={19} />
            {label}
          </button>
        ))}
      </nav>

      {/* User footer */}
      <div className="sidebar-footer">
        <button
          className="sidebar-item"
          style={{ width: '100%' }}
          onClick={() => navigate('profile')}
        >
          <div style={{
            width: 32, height: 32, borderRadius: '50%',
            background: 'var(--primary-btn)',
            color: '#fff', display: 'flex', alignItems: 'center', justifyContent: 'center',
            fontSize: 12, fontWeight: 700, flexShrink: 0,
          }}>
            {initials}
          </div>
          <span style={{ flex: 1, textAlign: 'left', fontSize: 13 }}>
            {claims?.full_name ?? 'Account'}
          </span>
        </button>
      </div>
    </aside>
  );
}
