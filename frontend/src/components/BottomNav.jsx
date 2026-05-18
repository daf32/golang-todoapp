import React from 'react';
import { useNav } from '../App';
import { ListIcon, CalendarIcon, PlusIcon, UserIcon, SettingsIcon } from './Icons';

export default function BottomNav({ onAdd }) {
  const { screen, navigate } = useNav();

  const item = (key, Icon, label) => (
    <button
      className={`bottom-nav-item ${screen === key ? 'active' : ''}`}
      onClick={() => navigate(key)}
    >
      <span style={{ display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
        <Icon size={21} />
      </span>
      <span>{label}</span>
    </button>
  );

  return (
    <nav className="bottom-nav">
      {item('tasks',    ListIcon,     'Tasks')}
      {item('calendar', CalendarIcon, 'Calendar')}

      <button className="fab" onClick={onAdd} style={{ width: 48, height: 48 }}>
        <PlusIcon size={20} />
      </button>

      {item('profile',  UserIcon,     'Profile')}
      {item('settings', SettingsIcon, 'Settings')}
    </nav>
  );
}
