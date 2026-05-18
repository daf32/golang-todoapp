import React from 'react';
import { LogoMark } from '../components/Icons';

export default function LoadingPage() {
  return (
    <div className="loading-page">
      <div className="loading-logo-wrap">
        <div className="logo-mark logo-mark-lg">
          <svg width="38" height="38" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.8" strokeLinecap="round" strokeLinejoin="round">
            <polyline points="20 6 9 17 4 12" />
          </svg>
        </div>
      </div>

      <p className="loading-app-name">Taskly</p>

      <div className="loading-dots">
        <span /><span /><span />
      </div>
    </div>
  );
}
