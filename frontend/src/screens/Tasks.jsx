import { useState, useEffect } from 'react';
import { Checkbox } from '../components/UI.jsx';
import { CheckSquare, Clock, Check, Pencil, Trash, Calendar as CalendarIcon, ChevronLeft, ChevronRight, Plus, Target, Search } from '../components/Icons.jsx';
import { fmtTaskMeta, fmtDayMonth, sameDay, MONTHS, DOW, groupTasks, computeStats, weeklyBuckets, fmtDuration } from '../lib/format.js';

export function TaskRow({ task, onToggle, onDelete, onEdit, style, mobile }) {
  const [removing, setRemoving] = useState(false);
  function del(e) { e && e.stopPropagation(); setRemoving(true); setTimeout(() => onDelete(task.id), 180); }
  const cls = mobile ? 'm-task' : 'task';
  return (
    <div className={`${cls} ${task.completed ? 'done' : ''} ${removing ? 'completing' : ''}`}
      style={{ ...style, ...(removing ? { opacity: 0, transform: 'scale(0.97)' } : {}) }}>
      <Checkbox done={task.completed} onToggle={() => onToggle(task)} />
      <div className="task-body" onClick={() => onEdit && onEdit(task)} style={onEdit ? { cursor: 'pointer' } : undefined}>
        <div className="task-title">{task.title}</div>
        {task.description && <div className="task-desc">{task.description}</div>}
        <div className="task-meta">
          <span className="task-time"><Clock size={13} />{fmtTaskMeta(task)}</span>
          {task.completed && task.completed_at && (
            <span className="task-time" style={{ color: 'var(--green)' }}><Check size={13} stroke={2.4} />done</span>
          )}
        </div>
      </div>
      {!mobile && (
        <div className="task-actions">
          {onEdit && <button className="icon-btn" onClick={(e) => { e.stopPropagation(); onEdit(task); }} aria-label="Edit"><Pencil size={16} /></button>}
          <button className="icon-btn danger" onClick={del} aria-label="Delete"><Trash size={16} /></button>
        </div>
      )}
    </div>
  );
}

export function TaskListView({ tasks, onToggle, onDelete, onEdit, mobile, searchQuery }) {
  if (!tasks.length) {
    const searching = searchQuery && searchQuery.trim().length > 0;
    return (
      <div className="empty">
        <div className="empty-glyph">{searching ? <Search size={26} /> : <CheckSquare size={28} />}</div>
        <div className="empty-title">{searching ? 'No matches' : 'All clear'}</div>
        <div className="empty-sub">
          {searching ? <>Nothing matches &ldquo;{searchQuery}&rdquo;.</> : 'Nothing here yet. Add a task to get going.'}
        </div>
      </div>
    );
  }
  const groups = groupTasks(tasks);
  let n = 0;
  return (
    <div>
      {groups.map(([label, list]) => (
        <div className="task-group" key={label}>
          <div className="task-group-head">
            <span className="section-label">{label}</span>
            <span className="fg-4 t-fine">{list.length}</span>
          </div>
          <div className="task-list">
            {list.map((t) => {
              const delay = `${Math.min(n++, 12) * 32}ms`;
              return <TaskRow key={t.id} task={t} onToggle={onToggle} onDelete={onDelete} onEdit={onEdit} mobile={mobile} style={{ animationDelay: delay }} />;
            })}
          </div>
        </div>
      ))}
    </div>
  );
}

function buildCalDays(year, month) {
  const first = new Date(year, month, 1);
  const startDow = first.getDay();
  const last = new Date(year, month + 1, 0).getDate();
  const days = [];
  for (let i = 0; i < startDow; i++) days.push({ date: new Date(year, month, -startDow + i + 1), other: true });
  for (let d = 1; d <= last; d++) days.push({ date: new Date(year, month, d), other: false });
  while (days.length % 7 !== 0) days.push({ date: new Date(year, month + 1, days.length - startDow - last + 1), other: true });
  return days;
}

export function CalendarGrid({ cursor, setCursor, selected, setSelected, tasks }) {
  const today = new Date();
  const days = buildCalDays(cursor.year, cursor.month);
  const tasksOn = (d) => tasks.filter((t) => sameDay(t.date || t.created_at, d));
  function shift(n) {
    setCursor((c) => {
      let m = c.month + n, y = c.year;
      if (m < 0) { m = 11; y--; }
      if (m > 11) { m = 0; y++; }
      return { year: y, month: m };
    });
  }
  return (
    <div className="cal-card card">
      <div className="cal-head">
        <div className="cal-month">{MONTHS[cursor.month]} {cursor.year}</div>
        <div className="cal-nav">
          <button className="icon-btn" onClick={() => shift(-1)} aria-label="Previous month"><ChevronLeft size={18} /></button>
          <button className="icon-btn" onClick={() => shift(1)} aria-label="Next month"><ChevronRight size={18} /></button>
        </div>
      </div>
      <div className="cal-grid">{DOW.map((d) => <div className="cal-dow" key={d}>{d[0]}</div>)}</div>
      <div className="cal-grid" style={{ marginTop: 4 }}>
        {days.map(({ date, other }, i) => {
          const isToday = sameDay(date, today);
          const isSel = sameDay(date, selected);
          const has = tasksOn(date).length > 0;
          return (
            <button key={i} className={`cal-day ${other ? 'other' : ''} ${isToday ? 'today' : ''} ${isSel ? 'selected' : ''}`}
              style={{ animationDelay: `${i * 9}ms` }}
              onClick={() => setSelected(new Date(date))}>
              {date.getDate()}
              {has && <span className="cal-dot" />}
            </button>
          );
        })}
      </div>
    </div>
  );
}

function StatCard({ icon, value, label, sub, delay }) {
  return (
    <div className="stat-card card" style={{ animationDelay: delay }}>
      <div className="stat-icon">{icon}</div>
      <div className="stat-value">{value}</div>
      <div className="stat-label">{label}</div>
      {sub && <div className="stat-delta">{sub}</div>}
    </div>
  );
}

export { StatCard };

function CompletionRing({ rate }) {
  const r = 50, c = 2 * Math.PI * r;
  const pct = rate == null ? 0 : Math.round(rate * 100);
  const [off, setOff] = useState(c);
  useEffect(() => { const id = requestAnimationFrame(() => setOff(c - (pct / 100) * c)); return () => cancelAnimationFrame(id); }, [pct, c]);
  return (
    <div className="ring">
      <svg width="116" height="116" viewBox="0 0 116 116">
        <circle className="ring-track" cx="58" cy="58" r={r} strokeWidth="11" />
        <circle className="ring-fill" cx="58" cy="58" r={r} strokeWidth="11" strokeDasharray={c} strokeDashoffset={off} />
      </svg>
      <div className="ring-center"><div className="ring-pct">{pct}%</div></div>
    </div>
  );
}

function BarChart({ buckets }) {
  // Each day's bar represents total activity = max(created, completed).
  // The ink portion is `completed`; the muted portion on top is open work
  // (created today but not yet completed). Using max() means a day with a
  // completion but no creation still renders a full ink bar.
  const max = Math.max(1, ...buckets.flatMap((b) => [b.created, b.completed]));
  return (
    <div>
      <div className="bars">
        {buckets.map((b, i) => {
          const total = Math.max(b.created, b.completed);
          const open = Math.max(0, b.created - b.completed);
          const empty = total === 0;
          return (
            <div className="bar-col" key={i}>
              <div className="bar-stack" title={`${b.created} created · ${b.completed} completed`}>
                {empty ? (
                  <div className="bar-empty" />
                ) : (
                  <>
                    {open > 0 && (
                      <div className="bar created"
                        style={{ height: `${(open / max) * 100}%`, animationDelay: `${i * 50}ms` }} />
                    )}
                    {b.completed > 0 && (
                      <div className="bar completed"
                        style={{ height: `${(b.completed / max) * 100}%`, animationDelay: `${i * 50 + 40}ms` }} />
                    )}
                  </>
                )}
              </div>
              <div className="bar-label">{b.label}</div>
            </div>
          );
        })}
      </div>
      <div className="legend">
        <span className="legend-item"><span className="legend-swatch" style={{ background: 'var(--accent)' }} />Completed</span>
        <span className="legend-item"><span className="legend-swatch" style={{ background: 'var(--surface-muted)' }} />Open</span>
      </div>
    </div>
  );
}

export function StatsView({ tasks, range }) {
  const fromDays = range === '7d' ? -7 : range === '30d' ? -30 : -365;
  const stats = computeStats(tasks, fromDays, 0);
  const buckets = weeklyBuckets(tasks);
  return (
    <div>
      <div className="stat-grid">
        <StatCard icon={<Plus size={18} />} value={stats.tasks_created} label="Tasks created" delay="0ms" />
        <StatCard icon={<Check size={18} stroke={2.4} />} value={stats.tasks_completed} label="Tasks completed" delay="60ms" />
        <StatCard icon={<Target size={18} />} value={stats.tasks_completed_rate == null ? '—' : `${Math.round(stats.tasks_completed_rate * 100)}%`} label="Completion rate" delay="120ms" />
        <StatCard icon={<Clock size={18} />} value={fmtDuration(stats.tasks_average_completion_time_ms)} label="Avg. completion" delay="180ms" />
      </div>

      <div className="chart-card card">
        <div className="ring-wrap">
          <CompletionRing rate={stats.tasks_completed_rate} />
          <div>
            <div className="t-h2">Completion rate</div>
            <div className="fg-3 t-caption" style={{ marginTop: 4, maxWidth: 220 }}>
              {stats.tasks_completed} of {stats.tasks_created} tasks completed in the selected range.
            </div>
          </div>
        </div>
      </div>

      <div className="chart-card card">
        <div className="chart-head">
          <div className="t-h2">This week</div>
          <span className="pill">Created vs completed</span>
        </div>
        <BarChart buckets={buckets} />
      </div>
    </div>
  );
}

export function CalendarEmpty() {
  return (
    <div className="empty" style={{ padding: '48px 16px' }}>
      <div className="empty-glyph"><CalendarIcon size={26} /></div>
      <div className="empty-sub">No tasks this day.</div>
    </div>
  );
}
