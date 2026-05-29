import { useState, useEffect, useRef } from 'react';
import { Check, X } from './Icons.jsx';

export function Segmented({ options, value, onChange }) {
  const ref = useRef(null);
  const [thumb, setThumb] = useState({ left: 0, width: 0 });
  useEffect(() => {
    const el = ref.current;
    if (!el) return;
    const idx = options.findIndex((o) => (o.value ?? o) === value);
    const btn = el.querySelectorAll('button')[idx];
    if (btn) setThumb({ left: btn.offsetLeft, width: btn.offsetWidth });
  }, [value, options]);
  return (
    <div className="segmented" ref={ref}>
      <div className="seg-thumb" style={{ left: thumb.left, width: thumb.width, top: 3, bottom: 3 }} />
      {options.map((o) => {
        const val = o.value ?? o;
        const label = o.label ?? o;
        return (
          <button key={val} className={value === val ? 'active' : ''} onClick={() => onChange(val)} type="button">{label}</button>
        );
      })}
    </div>
  );
}

export function Checkbox({ done, onToggle }) {
  return (
    <button type="button" className={`check ${done ? 'done' : ''}`} onClick={onToggle} aria-label={done ? 'Mark incomplete' : 'Mark complete'}>
      {done && <Check size={13} stroke={3} />}
    </button>
  );
}

function toDateInput(d) { d = new Date(d); return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`; }
function toTimeInput(d) { d = new Date(d); return `${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`; }

export function TaskForm({ initial, defaultDate, onSubmit, onCancel, submitLabel = 'Add task' }) {
  const base = initial?.date
    ? new Date(initial.date)
    : defaultDate
      ? new Date(defaultDate)
      : (() => { const d = new Date(); d.setHours(9, 0, 0, 0); return d; })();
  const [title, setTitle] = useState(initial?.title || '');
  const [desc, setDesc] = useState(initial?.description || '');
  const [dateStr, setDateStr] = useState(toDateInput(base));
  const [timeStr, setTimeStr] = useState(toTimeInput(base));
  const [err, setErr] = useState('');
  const ref = useRef(null);
  useEffect(() => { const id = setTimeout(() => ref.current?.focus({ preventScroll: true }), 60); return () => clearTimeout(id); }, []);

  const today = new Date();
  const tomorrow = new Date(); tomorrow.setDate(tomorrow.getDate() + 1);
  const isToday = dateStr === toDateInput(today);
  const isTomorrow = dateStr === toDateInput(tomorrow);

  function submit(e) {
    e.preventDefault();
    const t = title.trim();
    if (t.length < 1) { setErr('Title is required'); return; }
    if (t.length > 100) { setErr('Title must be under 100 characters'); return; }
    if (desc.length > 1000) { setErr('Description is too long'); return; }
    const [y, m, d] = dateStr.split('-').map(Number);
    const [hh, mm] = (timeStr || '09:00').split(':').map(Number);
    const date = new Date(y, m - 1, d, hh || 0, mm || 0, 0, 0);
    onSubmit({ title: t, description: desc.trim() || null, date: date.toISOString() });
  }

  return (
    <form className="modal-body" onSubmit={submit}>
      <div className="field">
        <input ref={ref} className="input" placeholder="What needs doing?" value={title}
          maxLength={100} onChange={(e) => { setTitle(e.target.value); setErr(''); }} style={{ fontSize: 17, fontWeight: 500 }} />
      </div>
      <div className="field">
        <textarea className="input" rows={3} placeholder="Add a description (optional)" value={desc}
          maxLength={1000} onChange={(e) => setDesc(e.target.value)} />
      </div>
      <div className="field">
        <label className="field-label">When</label>
        <div style={{ display: 'flex', gap: 8, marginBottom: 4 }}>
          <button type="button" className={`pill press ${isToday ? 'ink' : 'outline'}`} onClick={() => setDateStr(toDateInput(today))}>Today</button>
          <button type="button" className={`pill press ${isTomorrow ? 'ink' : 'outline'}`} onClick={() => setDateStr(toDateInput(tomorrow))}>Tomorrow</button>
        </div>
        <div style={{ display: 'flex', gap: 10 }}>
          <input type="date" className="input" value={dateStr} onChange={(e) => setDateStr(e.target.value)} style={{ flex: 2 }} />
          <input type="time" className="input" value={timeStr} onChange={(e) => setTimeStr(e.target.value)} style={{ flex: 1, minWidth: 0 }} />
        </div>
      </div>
      {err && <div className="field-err">{err}</div>}
      <div className="modal-foot">
        <button type="button" className="btn btn-pearl btn-block" onClick={onCancel}>Cancel</button>
        <button type="submit" className="btn btn-primary btn-block">{submitLabel}</button>
      </div>
    </form>
  );
}

export function Sheet({ onClose, children }) {
  const sheetRef = useRef(null);
  const drag = useRef(null);
  const [y, setY] = useState(900);
  const [dragging, setDragging] = useState(false);
  const [closing, setClosing] = useState(false);

  useEffect(() => {
    const id = requestAnimationFrame(() => requestAnimationFrame(() => setY(0)));
    return () => cancelAnimationFrame(id);
  }, []);
  useEffect(() => {
    const h = (e) => { if (e.key === 'Escape') close(); };
    window.addEventListener('keydown', h);
    return () => window.removeEventListener('keydown', h);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  function down(e) {
    const el = sheetRef.current;
    const scale = el ? el.getBoundingClientRect().height / el.offsetHeight || 1 : 1;
    drag.current = { y: e.clientY, t: Date.now(), lastY: e.clientY, lastT: Date.now(), scale };
    setDragging(true);
    try { e.currentTarget.setPointerCapture(e.pointerId); } catch (_) {}
  }
  function move(e) {
    if (!drag.current) return;
    const dy = (e.clientY - drag.current.y) / drag.current.scale;
    drag.current.lastY = e.clientY; drag.current.lastT = Date.now();
    setY(Math.max(0, dy));
  }
  function up() {
    if (!drag.current) return;
    const h = sheetRef.current?.offsetHeight || 600;
    const dt = Math.max(1, drag.current.lastT - drag.current.t);
    const vel = (drag.current.lastY - drag.current.y) / drag.current.scale / dt;
    const dist = y;
    drag.current = null;
    setDragging(false);
    if (dist > h * 0.32 || vel > 0.6) close(h);
    else setY(0);
  }
  function close(h) {
    const fall = (typeof h === 'number' ? h : sheetRef.current?.offsetHeight) || 600;
    setClosing(true);
    setDragging(false);
    setY(fall);
    setTimeout(onClose, 280);
  }

  const sh = sheetRef.current?.offsetHeight || 600;
  const scrimO = closing ? 0 : Math.max(0, Math.min(1, 1 - y / (sh * 1.1)));

  return (
    <div className="sheet-scrim" style={{ background: `rgba(0,0,0,${0.4 * scrimO})`, transition: dragging ? 'none' : 'background 280ms var(--ease)' }} onMouseDown={() => close()}>
      <div ref={sheetRef} className="sheet" onMouseDown={(e) => e.stopPropagation()}
        style={{ transform: `translateY(${y}px)`, transition: dragging ? 'none' : 'transform 280ms var(--ease)' }}>
        <div className="sheet-handle" onPointerDown={down} onPointerMove={move} onPointerUp={up} onPointerCancel={up}>
          <div className="sheet-grip" />
        </div>
        {children}
      </div>
    </div>
  );
}

export function Modal({ title, onClose, children }) {
  useEffect(() => {
    const h = (e) => { if (e.key === 'Escape') onClose(); };
    window.addEventListener('keydown', h);
    return () => window.removeEventListener('keydown', h);
  }, [onClose]);
  return (
    <div className="scrim" style={{ position: 'fixed' }} onMouseDown={onClose}>
      <div className="modal" onMouseDown={(e) => e.stopPropagation()}>
        <div className="modal-head">
          <div className="modal-title">{title}</div>
          <button className="icon-btn" onClick={onClose} aria-label="Close"><X size={18} /></button>
        </div>
        {children}
      </div>
    </div>
  );
}

export function ConfirmModal({ title, body, confirm, danger, onClose, onConfirm }) {
  return (
    <Modal title={title} onClose={onClose}>
      <div className="t-body fg-2" style={{ marginTop: -6 }}>{body}</div>
      <div className="modal-foot">
        <button className="btn btn-pearl btn-block" onClick={onClose}>Cancel</button>
        <button className={`btn btn-block ${danger ? '' : 'btn-primary'}`} onClick={onConfirm}
          style={danger ? { background: 'var(--red)', color: '#fff' } : undefined}>{confirm}</button>
      </div>
    </Modal>
  );
}
