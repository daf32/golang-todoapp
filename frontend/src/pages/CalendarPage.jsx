import React, { useState, useEffect, useCallback } from 'react';
import { useAuth } from '../App';
import { api } from '../api/client';
import AuthLayout from '../components/AuthLayout';
import TaskItem from '../components/TaskItem';
import AddTaskModal from '../components/AddTaskModal';
import { ChevronLeftIcon, ChevronRightIcon } from '../components/Icons';

const DAY_NAMES   = ['Sun','Mon','Tue','Wed','Thu','Fri','Sat'];
const MONTH_NAMES = ['January','February','March','April','May','June',
                     'July','August','September','October','November','December'];

function sameDay(a, b) {
  return a.getFullYear() === b.getFullYear()
      && a.getMonth()    === b.getMonth()
      && a.getDate()     === b.getDate();
}

function buildCalDays(year, month) {
  const first    = new Date(year, month, 1);
  const last     = new Date(year, month + 1, 0);
  const startDow = first.getDay();
  const days     = [];

  for (let i = 0; i < startDow; i++) {
    days.push({ date: new Date(year, month, -startDow + i + 1), thisMonth: false });
  }
  for (let d = 1; d <= last.getDate(); d++) {
    days.push({ date: new Date(year, month, d), thisMonth: true });
  }
  while (days.length % 7 !== 0) {
    days.push({ date: new Date(year, month + 1, days.length - startDow - last.getDate() + 1), thisMonth: false });
  }
  return days;
}

export default function CalendarPage() {
  const { logout }  = useAuth();
  const today       = new Date();
  const [cursor, setCursor]     = useState({ year: today.getFullYear(), month: today.getMonth() });
  const [selected, setSelected] = useState(today);
  const [tasks, setTasks]       = useState([]);
  const [loading, setLoading]   = useState(true);
  const [error, setError]       = useState('');
  const [showAdd, setShowAdd]   = useState(false);

  const fetchTasks = useCallback(async () => {
    setLoading(true); setError('');
    try {
      const data = await api.getTasks();
      setTasks(Array.isArray(data) ? data : []);
    } catch (err) {
      if (err.status === 401) { logout(); return; }
      setError(err.message || 'Failed to load tasks');
    } finally {
      setLoading(false);
    }
  }, [logout]);

  useEffect(() => { fetchTasks(); }, [fetchTasks]);

  function handleUpdate(updated) { setTasks((ts) => ts.map((t) => t.id === updated.id ? updated : t)); }
  function handleDelete(id)      { setTasks((ts) => ts.filter((t) => t.id !== id)); }
  function handleCreated(task)   { setTasks((ts) => [task, ...ts]); }

  function prevMonth() {
    setCursor(({ year, month }) =>
      month === 0 ? { year: year - 1, month: 11 } : { year, month: month - 1 }
    );
  }
  function nextMonth() {
    setCursor(({ year, month }) =>
      month === 11 ? { year: year + 1, month: 0 } : { year, month: month + 1 }
    );
  }

  const calDays       = buildCalDays(cursor.year, cursor.month);
  const tasksOnDay    = (date) => tasks.filter((t) => sameDay(new Date(t.created_at), date));
  const selectedTasks = tasksOnDay(selected);

  const selLabel = sameDay(selected, today)
    ? 'Today'
    : selected.toLocaleDateString('en-US', { weekday: 'long', month: 'short', day: 'numeric' });

  return (
    <AuthLayout title="Calendar" onAdd={() => setShowAdd(true)}>
      {/* Calendar – constrained width so circles stay small on desktop */}
      <div className="cal-wrap">
        <div className="cal-header">
          <button className="cal-nav-btn" onClick={prevMonth}><ChevronLeftIcon size={15} /></button>
          <span className="cal-month-label">{MONTH_NAMES[cursor.month]} {cursor.year}</span>
          <button className="cal-nav-btn" onClick={nextMonth}><ChevronRightIcon size={15} /></button>
        </div>

        {/* Weekday headers */}
        <div className="cal-grid" style={{ marginBottom: 2 }}>
          {DAY_NAMES.map((d) => <div key={d} className="cal-day-label">{d[0]}</div>)}
        </div>

        {/* Days */}
        <div className="cal-grid" style={{ marginBottom: 8 }}>
          {calDays.map(({ date, thisMonth }, i) => {
            const isToday    = sameDay(date, today);
            const isSelected = sameDay(date, selected);
            const hasTasks   = tasksOnDay(date).length > 0;
            return (
              <button
                key={i}
                className={[
                  'cal-day',
                  !thisMonth ? 'other-month' : '',
                  isToday    ? 'today'       : '',
                  isSelected ? 'selected'    : '',
                  hasTasks   ? 'has-tasks'   : '',
                ].join(' ')}
                onClick={() => setSelected(date)}
              >
                {date.getDate()}
              </button>
            );
          })}
        </div>
      </div>

      {/* Separator */}
      <div style={{ height: 1, background: 'var(--border)', margin: '0 20px' }} />

      {/* Day task list */}
      <div className="section-label" style={{ marginTop: 16 }}>{selLabel}</div>

      {error && <div className="error-msg">{error}</div>}

      {loading ? (
        <div className="spinner" />
      ) : selectedTasks.length === 0 ? (
        <div className="empty-state" style={{ paddingTop: 20, paddingBottom: 20 }}>
          <svg width="44" height="44" viewBox="0 0 44 44" fill="none">
            <circle cx="22" cy="22" r="18" fill="var(--primary-light)"/>
            <path d="M14 22h16M22 14v16" stroke="var(--primary-mid)" strokeWidth="2.5" strokeLinecap="round"/>
          </svg>
          <p style={{ fontSize: 13 }}>No tasks on this day. Tap <strong>+</strong> to add.</p>
        </div>
      ) : (
        selectedTasks.map((task, i) => (
          <TaskItem
            key={task.id}
            task={task}
            onUpdate={handleUpdate}
            onDelete={handleDelete}
            style={{ animationDelay: `${i * 35}ms` }}
          />
        ))
      )}

      <div style={{ height: 16 }} />

      {showAdd && (
        <AddTaskModal
          onClose={() => setShowAdd(false)}
          onCreated={handleCreated}
        />
      )}
    </AuthLayout>
  );
}
