import React, { useState, useEffect, useRef } from 'react';
import { api } from '../api/client';
import { PlusIcon } from './Icons';

const CLOSE_THRESHOLD = 100; // px dragged down to trigger dismiss
const EXIT_MS         = 280; // must match exit transition duration

export default function AddTaskModal({ onClose, onCreated }) {
  const [title,   setTitle]   = useState('');
  const [desc,    setDesc]    = useState('');
  const [loading, setLoading] = useState(false);
  const [error,   setError]   = useState('');

  // 'entering' → 'idle' → 'dragging' → 'exiting'
  const [phase, setPhase] = useState('entering');
  const [dragY, setDragY] = useState(0);
  const dragStartY        = useRef(0);
  const inputRef          = useRef(null);

  // Two-frame trick: render at bottom, then snap to idle so the CSS
  // transition creates the slide-up entrance.
  useEffect(() => {
    const id = requestAnimationFrame(() =>
      requestAnimationFrame(() => setPhase('idle'))
    );
    return () => cancelAnimationFrame(id);
  }, []);

  // Focus the title input once the sheet has arrived
  useEffect(() => {
    if (phase !== 'idle') return;
    const t = setTimeout(() => inputRef.current?.focus(), 60);
    return () => clearTimeout(t);
  }, [phase]);

  // Trigger slide-down exit, then call onClose after the transition finishes
  function close() {
    if (phase === 'exiting') return;
    setDragY(0);
    setPhase('exiting');
    setTimeout(onClose, EXIT_MS);
  }

  /* ── Drag-to-dismiss ─────────────────────────────── */
  function handlePointerDown(e) {
    if (phase === 'exiting') return;
    dragStartY.current = e.clientY;
    setPhase('dragging');
    e.currentTarget.setPointerCapture(e.pointerId);
  }

  function handlePointerMove(e) {
    if (phase !== 'dragging') return;
    setDragY(Math.max(0, e.clientY - dragStartY.current));
  }

  function handlePointerUp(e) {
    if (phase !== 'dragging') return;
    const dy = e.clientY - dragStartY.current;
    if (dy >= CLOSE_THRESHOLD) {
      close();
    } else {
      setPhase('idle');
      setDragY(0);
    }
  }

  /* ── Derive sheet transform / transition from phase ── */
  const ENTER_SPRING  = 'transform .34s cubic-bezier(.22,1,.36,1), opacity .25s ease';
  const EXIT_EASE     = `transform ${EXIT_MS}ms cubic-bezier(.55,.05,.68,.19), opacity ${EXIT_MS}ms ease`;
  const SNAP_SPRING   = 'transform .3s cubic-bezier(.22,1,.36,1)';

  let sheetStyle;
  switch (phase) {
    case 'entering':
      sheetStyle = { transform: 'translateY(100%)', opacity: 0, transition: 'none' };
      break;
    case 'idle':
      sheetStyle = { transform: 'translateY(0)', opacity: 1, transition: ENTER_SPRING };
      break;
    case 'dragging':
      sheetStyle = {
        transform:  `translateY(${dragY}px)`,
        opacity:    Math.max(0.4, 1 - dragY / 280),
        transition: 'none',
        cursor:     'grabbing',
      };
      break;
    case 'exiting':
      sheetStyle = { transform: 'translateY(110%)', opacity: 0, transition: EXIT_EASE };
      break;
    default:
      sheetStyle = {};
  }

  const overlayOpacity =
    phase === 'exiting' ? 0 :
    phase === 'dragging' ? Math.max(0, 1 - dragY / 280) :
    1;

  /* ── Submit ── */
  async function handleSubmit(e) {
    e.preventDefault();
    const t = title.trim();
    if (!t) return;
    setLoading(true);
    setError('');
    try {
      const task = await api.createTask(t, desc.trim() || undefined);
      onCreated?.(task);
      close();
    } catch (err) {
      setError(err.message || 'Failed to create task');
      setLoading(false);
    }
  }

  return (
    <div
      className="modal-overlay"
      onClick={(e) => e.target === e.currentTarget && close()}
      style={{
        opacity:    overlayOpacity,
        transition: phase === 'exiting' ? `opacity ${EXIT_MS}ms ease` :
                    phase === 'dragging' ? 'none' : undefined,
        /* override the CSS fadeIn — we drive opacity ourselves */
        animation: 'none',
      }}
    >
      <div
        className="modal-sheet"
        /* override the CSS slideUpModal — we drive transform ourselves */
        style={{ ...sheetStyle, animation: 'none' }}
      >
        {/* ── Drag handle — wide tap target ── */}
        <div
          style={{
            padding: '4px 0 14px',
            display: 'flex',
            justifyContent: 'center',
            cursor: phase === 'dragging' ? 'grabbing' : 'grab',
            touchAction: 'none',
            userSelect: 'none',
          }}
          onPointerDown={handlePointerDown}
          onPointerMove={handlePointerMove}
          onPointerUp={handlePointerUp}
        >
          <div
            style={{
              width: 40, height: 4,
              borderRadius: 2,
              background: 'var(--border)',
              /* Widen the visual handle a bit when dragging */
              transition: 'width .15s ease',
              ...(phase === 'dragging' ? { width: 56 } : {}),
            }}
          />
        </div>

        <p className="modal-title">What would you like to do?</p>

        {error && (
          <div className="error-msg" style={{ marginLeft: 0, marginRight: 0, marginBottom: 12 }}>
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
          <input
            ref={inputRef}
            className="input input-no-icon"
            placeholder="Task title…"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            maxLength={100}
          />

          <input
            className="input input-no-icon"
            placeholder="Description (optional)…"
            value={desc}
            onChange={(e) => setDesc(e.target.value)}
            maxLength={1000}
          />

          <button
            type="submit"
            className="btn btn-primary"
            disabled={loading || !title.trim()}
            style={{ marginTop: 4 }}
          >
            <PlusIcon size={18} />
            {loading ? 'Adding…' : 'Add Task'}
          </button>
        </form>
      </div>
    </div>
  );
}
