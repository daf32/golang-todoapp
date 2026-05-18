import React from 'react';
import { useNav } from '../App';
import { LogoMark, ArrowRightIcon } from '../components/Icons';

export default function OnboardingPage() {
  const { navigate } = useNav();

  return (
    <div className="page" style={{ background: '#fff', justifyContent: 'space-between', padding: '64px 32px 48px' }}>
      {/* Background blob */}
      <div
        className="blob blob-primary"
        style={{ bottom: -80, left: -80, width: 280, height: 280, opacity: .55, zIndex: 0 }}
      />

      {/* Top decorative dots */}
      <div style={{ position: 'absolute', top: 48, right: 48, display: 'flex', gap: 6 }}>
        <span style={{ width: 8, height: 8, borderRadius: '50%', background: '#7B6CF6', display: 'block' }} />
        <span style={{ width: 8, height: 8, borderRadius: '50%', background: '#F97316', display: 'block' }} />
        <span style={{ width: 8, height: 8, borderRadius: '50%', background: '#EAB308', display: 'block' }} />
      </div>

      {/* Logo + headline */}
      <div style={{ position: 'relative', zIndex: 1 }}>
        <LogoMark />

        <h1 style={{ fontSize: 34, fontWeight: 800, marginTop: 32, lineHeight: 1.2, letterSpacing: '-0.5px' }}>
          Get things<br />done.
        </h1>
        <p style={{ color: 'var(--text-2)', marginTop: 12, fontSize: 15, lineHeight: 1.6 }}>
          Start planning and organizing<br />your tasks
        </p>
      </div>

      {/* Footer: dots + arrow */}
      <div style={{ position: 'relative', zIndex: 1, display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
        <div style={{ display: 'flex', gap: 6 }}>
          <span style={{ width: 22, height: 8, borderRadius: 4, background: 'var(--primary-btn)', display: 'block' }} />
          <span style={{ width: 8, height: 8, borderRadius: '50%', background: 'var(--border)', display: 'block' }} />
          <span style={{ width: 8, height: 8, borderRadius: '50%', background: 'var(--border)', display: 'block' }} />
        </div>

        <button className="fab" onClick={() => navigate('signup')}>
          <ArrowRightIcon size={22} />
        </button>
      </div>
    </div>
  );
}
