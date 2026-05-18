import React from 'react';
import Sidebar from './Sidebar';
import BottomNav from './BottomNav';
import { useNav } from '../App';

const DAY_NAMES   = ['Sun','Mon','Tue','Wed','Thu','Fri','Sat'];
const MONTH_NAMES = ['January','February','March','April','May','June',
                     'July','August','September','October','November','December'];

function todayLabel() {
  const d = new Date();
  return `Today · ${DAY_NAMES[d.getDay()]} ${d.getDate()} ${MONTH_NAMES[d.getMonth()]}`;
}

export default function AuthLayout({ title, date, actions, onAdd, children }) {
  return (
    <div className="auth-layout">
      <Sidebar />

      <div className="auth-main">
        {/* Page header */}
        <div className="page-header">
          <div>
            <div className="page-date">{date ?? todayLabel()}</div>
            <div className="page-title">{title}</div>
          </div>
          {actions && (
            <div className="page-header-right">
              {actions}
            </div>
          )}
        </div>

        {/* Scrollable body */}
        <div className="scroll-area" style={{ flex: 1 }}>
          {children}
        </div>

        {/* Mobile bottom nav */}
        <BottomNav onAdd={onAdd} />
      </div>
    </div>
  );
}
