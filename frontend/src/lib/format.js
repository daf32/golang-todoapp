export const MONTHS = ['January','February','March','April','May','June','July','August','September','October','November','December'];
export const MON_SHORT = ['Jan','Feb','Mar','Apr','May','Jun','Jul','Aug','Sep','Oct','Nov','Dec'];
export const DOW = ['Sun','Mon','Tue','Wed','Thu','Fri','Sat'];

const pad = (n) => String(n).padStart(2, '0');

export function fmtTimeH(d) { d = new Date(d); return `${pad(d.getHours())}:${pad(d.getMinutes())}h`; }
export function fmtDayMonth(d) { d = new Date(d); return `${d.getDate()} ${MON_SHORT[d.getMonth()]}`; }
export function sameDay(a, b) { a = new Date(a); b = new Date(b); return a.getFullYear() === b.getFullYear() && a.getMonth() === b.getMonth() && a.getDate() === b.getDate(); }

export function fmtTaskMeta(t) {
  const d = new Date(t.date || t.created_at);
  const now = new Date();
  const diffDays = Math.round((new Date(d.getFullYear(), d.getMonth(), d.getDate()) - new Date(now.getFullYear(), now.getMonth(), now.getDate())) / 86400000);
  const time = fmtTimeH(d);
  if (diffDays === 0) return `today at ${time}`;
  if (diffDays === 1) return `tomorrow at ${time}`;
  if (diffDays === -1) return `yesterday at ${time}`;
  return `${fmtDayMonth(d)} at ${time}`;
}

export function fmtDuration(ms) {
  if (ms == null) return '—';
  const mins = Math.round(ms / 60000);
  if (mins < 60) return `${mins}m`;
  const hrs = mins / 60;
  if (hrs < 24) return `${hrs % 1 === 0 ? hrs : hrs.toFixed(1)}h`;
  const days = hrs / 24;
  return `${days % 1 === 0 ? days : days.toFixed(1)}d`;
}

export function initials(name) {
  if (!name) return '?';
  return name.trim().split(/\s+/).map((w) => w[0]).slice(0, 2).join('').toUpperCase();
}

export function atDay(offsetDays, h = 9, m = 0) {
  const d = new Date();
  d.setDate(d.getDate() + offsetDays);
  d.setHours(h, m, 0, 0);
  return d;
}

export function computeStats(tasks, fromDays, toDays) {
  let list = tasks;
  if (fromDays != null || toDays != null) {
    list = tasks.filter((t) => {
      const d = new Date(t.created_at);
      const lo = fromDays != null ? atDay(fromDays, 0, 0) : null;
      const hi = toDays != null ? atDay(toDays, 23, 59) : null;
      return (!lo || d >= lo) && (!hi || d <= hi);
    });
  }
  const created = list.length;
  const completedList = list.filter((t) => t.completed && t.completed_at);
  const completed = completedList.length;
  const rate = created > 0 ? completed / created : null;
  let avg = null;
  if (completedList.length) {
    const total = completedList.reduce((s, t) => s + (new Date(t.completed_at) - new Date(t.created_at)), 0);
    avg = total / completedList.length;
  }
  return { tasks_created: created, tasks_completed: completed, tasks_completed_rate: rate, tasks_average_completion_time_ms: avg };
}

export function weeklyBuckets(tasks) {
  const out = [];
  for (let i = 6; i >= 0; i--) {
    const day = atDay(-i, 0, 0);
    const created = tasks.filter((t) => sameDay(t.created_at, day)).length;
    const completed = tasks.filter((t) => t.completed && t.completed_at && sameDay(t.completed_at, day)).length;
    out.push({ label: DOW[day.getDay()][0], date: day, created, completed });
  }
  return out;
}

export function groupTasks(tasks) {
  const map = {};
  const order = [];
  for (const t of tasks) {
    const d = new Date(t.date || t.created_at);
    const now = new Date();
    const diff = Math.round((new Date(d.getFullYear(), d.getMonth(), d.getDate()) - new Date(now.getFullYear(), now.getMonth(), now.getDate())) / 86400000);
    let key;
    if (diff < 0) key = 'Earlier';
    else if (diff === 0) key = 'Today';
    else if (diff === 1) key = 'Tomorrow';
    else if (diff <= 7) key = 'This week';
    else key = 'Later';
    if (!map[key]) { map[key] = []; order.push(key); }
    map[key].push(t);
  }
  const RANK = { Today: 0, Tomorrow: 1, 'This week': 2, Later: 3, Earlier: 4 };
  order.sort((a, b) => RANK[a] - RANK[b]);
  return order.map((k) => [k, map[k]]);
}
