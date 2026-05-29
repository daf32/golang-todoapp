import { useState, useEffect, useMemo, useCallback, useRef } from 'react';
import { useTheme } from './hooks/useTheme.js';
import { api } from './api/client.js';
import {
  CheckSquare, Inbox, Calendar, Chart, Settings, Search, Plus, X, User,
} from './components/Icons.jsx';
import { Segmented, TaskForm, Sheet, Modal, ConfirmModal } from './components/UI.jsx';
import { TaskListView, TaskRow, CalendarGrid, StatsView, CalendarEmpty } from './screens/Tasks.jsx';
import { ProfileView, SettingsView, ChangePasswordForm } from './screens/Account.jsx';
import { Login, SignUp, EmailConfirm, EmailConfirmed } from './screens/Auth.jsx';
import { initials, sameDay, fmtDayMonth, DOW } from './lib/format.js';

const NAV = [
  { id: 'inbox', label: 'Inbox', icon: Inbox },
  { id: 'calendar', label: 'Calendar', icon: Calendar },
  { id: 'stats', label: 'Statistics', icon: Chart },
];

function BootSplash({ visible }) {
  return (
    <div className={`boot ${visible ? '' : 'hidden'}`} aria-hidden={!visible}>
      <div className="boot-mark">
        <div className="boot-logo"><CheckSquare size={32} /></div>
        <div className="boot-word">todo</div>
      </div>
      <div className="boot-bar" />
    </div>
  );
}

function decodeJwtUserId(token) {
  try {
    const payload = JSON.parse(atob(token.split('.')[1].replace(/-/g, '+').replace(/_/g, '/')));
    return payload.user_id ?? payload.sub ?? payload.uid ?? null;
  } catch (_) { return null; }
}

async function fetchUserFromToken(token) {
  const id = decodeJwtUserId(token);
  if (id == null) return null;
  try { return await api.getUser(id); } catch (_) { return null; }
}

function useAuth() {
  const [user, setUser] = useState(() => {
    try { return JSON.parse(localStorage.getItem('user') || 'null'); } catch { return null; }
  });
  const [authed, setAuthed] = useState(() => !!localStorage.getItem('access_token'));

  const login = useCallback(async (accessToken, refreshToken, u) => {
    localStorage.setItem('access_token', accessToken);
    if (refreshToken) localStorage.setItem('refresh_token', refreshToken);
    let finalUser = u || null;
    if (!finalUser) finalUser = await fetchUserFromToken(accessToken);
    if (finalUser) {
      localStorage.setItem('user', JSON.stringify(finalUser));
      setUser(finalUser);
    }
    setAuthed(true);
  }, []);

  // If we already have a token but no user (e.g. older session), hydrate the user.
  useEffect(() => {
    const t = localStorage.getItem('access_token');
    if (t && !user) {
      fetchUserFromToken(t).then((u) => {
        if (u) { localStorage.setItem('user', JSON.stringify(u)); setUser(u); }
      });
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const logout = useCallback(async () => {
    try { await api.logout(); } catch (_) {}
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    localStorage.removeItem('user');
    setUser(null);
    setAuthed(false);
  }, []);

  const updateUser = useCallback((patch) => {
    setUser((u) => {
      const next = { ...u, ...patch };
      localStorage.setItem('user', JSON.stringify(next));
      return next;
    });
  }, []);

  return { user, authed, login, logout, updateUser };
}

function useTasks(authed) {
  const [tasks, setTasks] = useState([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (!authed) { setTasks([]); return; }
    setLoading(true);
    api.getTasks({ limit: 200 }).then((res) => {
      const list = Array.isArray(res) ? res : (res?.tasks || res?.items || []);
      setTasks(list);
    }).catch(() => setTasks([])).finally(() => setLoading(false));
  }, [authed]);

  const toggleTask = useCallback(async (task) => {
    const completed = !task.completed;
    setTasks((ts) => ts.map((t) => t.id === task.id ? { ...t, completed, completed_at: completed ? new Date().toISOString() : null } : t));
    try {
      const upd = await api.patchTask(task.id, { completed });
      if (upd) setTasks((ts) => ts.map((t) => t.id === task.id ? upd : t));
    } catch (_) { /* revert on error */ }
  }, []);

  const deleteTask = useCallback(async (id) => {
    setTasks((ts) => ts.filter((t) => t.id !== id));
    try { await api.deleteTask(id); } catch (_) {}
  }, []);

  const addTask = useCallback(async ({ title, description, date }) => {
    const optimistic = {
      id: `tmp-${Date.now()}`, version: 1, title, description: description || null,
      completed: false, created_at: new Date().toISOString(), completed_at: null,
      date: date || new Date().toISOString(),
    };
    setTasks((ts) => [optimistic, ...ts]);
    try {
      const real = await api.createTask({ title, description: description || undefined, date });
      if (real) setTasks((ts) => [real, ...ts.filter((t) => t.id !== optimistic.id)]);
    } catch (_) {}
  }, []);

  const editTask = useCallback(async (id, patch) => {
    setTasks((ts) => ts.map((t) => t.id === id ? { ...t, ...patch } : t));
    try {
      const upd = await api.patchTask(id, patch);
      if (upd) setTasks((ts) => ts.map((t) => t.id === id ? upd : t));
    } catch (_) {}
  }, []);

  return { tasks, loading, toggleTask, deleteTask, addTask, editTask };
}

function useViewport() {
  const [w, setW] = useState(() => window.innerWidth);
  useEffect(() => {
    const onR = () => setW(window.innerWidth);
    window.addEventListener('resize', onR);
    return () => window.removeEventListener('resize', onR);
  }, []);
  return { isMobile: w <= 768 };
}

/* ════════════ Desktop Shell ════════════ */
function DesktopApp({ user, tasks, theme, setTheme, onLogout, onUpdateUser, onChangePassword, onDeleteAccount, taskOps }) {
  const { toggleTask, deleteTask, addTask, editTask } = taskOps;
  const [route, setRoute] = useState('inbox');
  const [query, setQuery] = useState('');
  const [filter, setFilter] = useState('all');
  const [modal, setModal] = useState(null);
  const [cursor, setCursor] = useState(() => { const d = new Date(); return { year: d.getFullYear(), month: d.getMonth() }; });
  const [selected, setSelected] = useState(new Date());
  const [range, setRange] = useState('7d');

  const visibleTasks = useMemo(() => {
    let list = tasks;
    if (filter === 'active') list = list.filter((t) => !t.completed);
    if (filter === 'completed') list = list.filter((t) => t.completed);
    if (query.trim()) {
      const q = query.toLowerCase();
      list = list.filter((t) => t.title.toLowerCase().includes(q) || (t.description || '').toLowerCase().includes(q));
    }
    return list;
  }, [tasks, filter, query]);

  const dayTasks = useMemo(() => tasks.filter((t) => sameDay(t.date || t.created_at, selected)), [tasks, selected]);
  const activeCount = tasks.filter((t) => !t.completed).length;

  function closeModal() { setModal(null); }
  function submitTask(data) {
    if (modal?.type === 'edit') editTask(modal.task.id, data);
    else addTask(data);
    closeModal();
  }

  return (
    <div className="shell">
      <nav className="rail">
        <div className="rail-logo"><CheckSquare size={24} /></div>
        {NAV.map((n) => {
          const Icon = n.icon;
          return (
            <button key={n.id} className={`rail-item ${route === n.id ? 'active' : ''}`} onClick={() => setRoute(n.id)}>
              <Icon size={22} />
              <span className="rail-tip">{n.label}</span>
            </button>
          );
        })}
        <div className="rail-spacer" />
        <button className={`rail-item ${route === 'settings' ? 'active' : ''}`} onClick={() => setRoute('settings')}>
          <Settings size={22} /><span className="rail-tip">Settings</span>
        </button>
        <button className={`rail-avatar ${route === 'profile' ? 'active' : ''}`} onClick={() => setRoute('profile')}>
          {initials(user.full_name)}
        </button>
      </nav>

      <main className="content">
        {route === 'inbox' && (
          <>
            <header className="page-head">
              <div className="page-head-inner">
                <div className="page-title-row">
                  <div style={{ flex: 1 }}>
                    <div className="page-title">Inbox</div>
                    <div className="page-sub">{activeCount} active · {tasks.length - activeCount} done</div>
                  </div>
                  <button className="btn btn-primary" onClick={() => setModal({ type: 'add' })}><Plus size={18} />New task</button>
                </div>
                <div className="toolbar">
                  <div className="search">
                    <Search size={17} />
                    <input type="text" placeholder="Search tasks" value={query} onChange={(e) => setQuery(e.target.value)} autoComplete="off" spellCheck={false} />
                    {query && <button type="button" className="icon-btn" onClick={() => setQuery('')} aria-label="Clear search"><X size={15} /></button>}
                  </div>
                  <Segmented value={filter} onChange={setFilter} options={[{ value: 'all', label: 'All' }, { value: 'active', label: 'Active' }, { value: 'completed', label: 'Done' }]} />
                </div>
              </div>
            </header>
            <div className="content-scroll">
              <div className="content-inner">
                <TaskListView tasks={visibleTasks} onToggle={toggleTask} onDelete={deleteTask} onEdit={(t) => setModal({ type: 'edit', task: t })} searchQuery={query} />
              </div>
            </div>
          </>
        )}

        {route === 'calendar' && (
          <>
            <header className="page-head"><div className="page-head-inner"><div className="page-title">Calendar</div><div className="page-sub">Plan tasks across your week</div></div></header>
            <div className="content-scroll">
              <div className="content-inner">
                <div className="cal-layout">
                  <CalendarGrid cursor={cursor} setCursor={setCursor} selected={selected} setSelected={setSelected} tasks={tasks} />
                  <div className="cal-side">
                    <div className="cal-side-head">
                      <div className="cal-side-date">{DOW[selected.getDay()]}, {fmtDayMonth(selected)}</div>
                      <button className="icon-btn" onClick={() => setModal({ type: 'add-day' })} aria-label="Add task this day"><Plus size={18} /></button>
                    </div>
                    {dayTasks.length ? (
                      <div style={{ marginTop: 12 }}>
                        {dayTasks.map((t) => <TaskRow key={t.id} task={t} onToggle={toggleTask} onDelete={deleteTask} onEdit={(tk) => setModal({ type: 'edit', task: tk })} />)}
                      </div>
                    ) : <CalendarEmpty />}
                  </div>
                </div>
              </div>
            </div>
          </>
        )}

        {route === 'stats' && (
          <>
            <header className="page-head">
              <div className="page-head-inner">
                <div className="page-title-row">
                  <div style={{ flex: 1 }}><div className="page-title">Statistics</div><div className="page-sub">Your productivity at a glance</div></div>
                  <Segmented value={range} onChange={setRange} options={[{ value: '7d', label: '7 days' }, { value: '30d', label: '30 days' }, { value: 'all', label: 'All' }]} />
                </div>
              </div>
            </header>
            <div className="content-scroll"><div className="content-inner"><StatsView tasks={tasks} range={range} /></div></div>
          </>
        )}

        {route === 'profile' && (
          <>
            <header className="page-head"><div className="page-head-inner"><div className="page-title">Profile</div></div></header>
            <div className="content-scroll"><div className="content-inner"><ProfileView user={user} tasks={tasks} onSave={onUpdateUser} /></div></div>
          </>
        )}

        {route === 'settings' && (
          <>
            <header className="page-head"><div className="page-head-inner"><div className="page-title">Settings</div></div></header>
            <div className="content-scroll"><div className="content-inner">
              <SettingsView user={user} theme={theme} setTheme={setTheme}
                onChangePassword={() => setModal({ type: 'password' })}
                onLogout={() => setModal({ type: 'logout' })}
                onDelete={() => setModal({ type: 'delete' })} />
            </div></div>
          </>
        )}
      </main>

      {modal?.type === 'add' && <Modal title="New task" onClose={closeModal}><TaskForm onSubmit={submitTask} onCancel={closeModal} /></Modal>}
      {modal?.type === 'add-day' && <Modal title={`New task · ${fmtDayMonth(selected)}`} onClose={closeModal}><TaskForm defaultDate={selected} onSubmit={submitTask} onCancel={closeModal} /></Modal>}
      {modal?.type === 'edit' && <Modal title="Edit task" onClose={closeModal}><TaskForm initial={modal.task} submitLabel="Save changes" onSubmit={submitTask} onCancel={closeModal} /></Modal>}
      {modal?.type === 'password' && <Modal title="Change password" onClose={closeModal}><ChangePasswordForm onClose={closeModal} onSubmit={onChangePassword} /></Modal>}
      {modal?.type === 'logout' && <ConfirmModal title="Log out?" body="You'll need to sign in again to access your tasks." confirm="Log out" onClose={closeModal} onConfirm={() => { closeModal(); onLogout(); }} />}
      {modal?.type === 'delete' && <ConfirmModal danger title="Delete account?" body="This permanently removes your account and all tasks. This can't be undone." confirm="Delete" onClose={closeModal} onConfirm={() => { closeModal(); onDeleteAccount(); }} />}
    </div>
  );
}

/* ════════════ Mobile Shell ════════════ */
function MobileApp({ user, tasks, theme, setTheme, onLogout, onUpdateUser, onChangePassword, onDeleteAccount, taskOps }) {
  const { toggleTask, deleteTask, addTask, editTask } = taskOps;
  const [route, setRoute] = useState('inbox');
  const [query, setQuery] = useState('');
  const [filter, setFilter] = useState('all');
  const [sheet, setSheet] = useState(null);
  const [searchOpen, setSearchOpen] = useState(false);
  const searchInputRef = useRef(null);
  useEffect(() => {
    if (searchOpen) searchInputRef.current?.focus({ preventScroll: true });
  }, [searchOpen]);
  const [cursor, setCursor] = useState(() => { const d = new Date(); return { year: d.getFullYear(), month: d.getMonth() }; });
  const [selected, setSelected] = useState(new Date());
  const [range, setRange] = useState('7d');

  const visibleTasks = useMemo(() => {
    let list = tasks;
    if (filter === 'active') list = list.filter((t) => !t.completed);
    if (filter === 'completed') list = list.filter((t) => t.completed);
    if (query.trim()) { const q = query.toLowerCase(); list = list.filter((t) => t.title.toLowerCase().includes(q) || (t.description || '').toLowerCase().includes(q)); }
    return list;
  }, [tasks, filter, query]);
  const dayTasks = useMemo(() => tasks.filter((t) => sameDay(t.date || t.created_at, selected)), [tasks, selected]);

  function closeSheet() { setSheet(null); }
  function submitTask(data) {
    if (sheet?.type === 'edit') editTask(sheet.task.id, data);
    else addTask(data);
    closeSheet();
  }

  const TABS = [...NAV, { id: 'profile', label: 'Profile', icon: User }];
  const titles = { inbox: 'Inbox', calendar: 'Calendar', stats: 'Stats', profile: 'Profile', settings: 'Settings' };

  return (
    <div className="m-shell">
      <div className="m-status-pad" />
      <div className="m-head">
        {route === 'inbox' && (
          <>
            <div className="m-head-row">
              {!searchOpen ? <div className="m-title">Inbox</div> : (
                <div className="search" style={{ flex: 1 }}>
                  <Search size={17} />
                  <input ref={searchInputRef} type="text" placeholder="Search" value={query} onChange={(e) => setQuery(e.target.value)} autoComplete="off" spellCheck={false} />
                  {query && <button type="button" className="icon-btn" onClick={() => setQuery('')} aria-label="Clear search"><X size={15} /></button>}
                </div>
              )}
              <button type="button" className="icon-btn" style={{ width: 38, height: 38 }} onClick={() => { setSearchOpen((s) => !s); if (searchOpen) setQuery(''); }}>
                {searchOpen ? <X size={20} /> : <Search size={20} />}
              </button>
            </div>
            <div style={{ marginTop: 10 }}>
              <Segmented value={filter} onChange={setFilter} options={[{ value: 'all', label: 'All' }, { value: 'active', label: 'Active' }, { value: 'completed', label: 'Done' }]} />
            </div>
          </>
        )}
        {route !== 'inbox' && (
          <div className="m-head-row">
            <div className="m-title">{titles[route]}</div>
            {route === 'stats' && <Segmented value={range} onChange={setRange} options={[{ value: '7d', label: '7d' }, { value: '30d', label: '30d' }, { value: 'all', label: 'All' }]} />}
          </div>
        )}
      </div>

      <div className="m-scroll">
        {route === 'inbox' && <TaskListView tasks={visibleTasks} onToggle={toggleTask} onDelete={deleteTask} onEdit={(t) => setSheet({ type: 'edit', task: t })} mobile searchQuery={query} />}
        {route === 'calendar' && (
          <div>
            <CalendarGrid cursor={cursor} setCursor={setCursor} selected={selected} setSelected={setSelected} tasks={tasks} />
            <div className="cal-side-head" style={{ marginTop: 22 }}>
              <div className="cal-side-date">{DOW[selected.getDay()]}, {fmtDayMonth(selected)}</div>
            </div>
            <div style={{ marginTop: 12 }}>
              {dayTasks.length ? dayTasks.map((t) => <TaskRow key={t.id} task={t} onToggle={toggleTask} onDelete={deleteTask} onEdit={(tk) => setSheet({ type: 'edit', task: tk })} mobile />)
                : <CalendarEmpty />}
            </div>
          </div>
        )}
        {route === 'stats' && <StatsView tasks={tasks} range={range} />}
        {route === 'profile' && <ProfileView user={user} tasks={tasks} onSave={onUpdateUser} />}
        {route === 'settings' && (
          <SettingsView user={user} theme={theme} setTheme={setTheme}
            onChangePassword={() => setSheet({ type: 'password' })}
            onLogout={onLogout}
            onDelete={onDeleteAccount} />
        )}
        <div className="m-pad-bottom" />
      </div>

      {(route === 'inbox' || route === 'calendar') && (
        <button className="fab" onClick={() => setSheet({ type: route === 'calendar' ? 'add-day' : 'add' })}><Plus size={26} /></button>
      )}

      <nav className="m-tabbar">
        {TABS.map((t) => {
          const Icon = t.icon;
          return (
            <button key={t.id} className={`m-tab ${route === t.id ? 'active' : ''}`} onClick={() => setRoute(t.id)}>
              <Icon size={23} /><span className="m-tab-label">{t.label}</span>
            </button>
          );
        })}
        <button className={`m-tab ${route === 'settings' ? 'active' : ''}`} onClick={() => setRoute('settings')}>
          <Settings size={23} /><span className="m-tab-label">Settings</span>
        </button>
      </nav>

      {sheet && (
        <Sheet onClose={closeSheet}>
          {(sheet.type === 'add' || sheet.type === 'add-day') && <><div className="modal-title" style={{ marginBottom: 16 }}>New task</div><TaskForm defaultDate={sheet.type === 'add-day' ? selected : undefined} onSubmit={submitTask} onCancel={closeSheet} /></>}
          {sheet.type === 'edit' && <><div className="modal-title" style={{ marginBottom: 16 }}>Edit task</div><TaskForm initial={sheet.task} submitLabel="Save changes" onSubmit={submitTask} onCancel={closeSheet} /></>}
          {sheet.type === 'password' && <><div className="modal-title" style={{ marginBottom: 16 }}>Change password</div><ChangePasswordForm onClose={closeSheet} onSubmit={onChangePassword} /></>}
        </Sheet>
      )}
    </div>
  );
}

/* ════════════ Root ════════════ */
export default function App() {
  const [theme, setTheme] = useTheme();
  const auth = useAuth();
  const taskOps = useTasks(auth.authed);
  const { isMobile } = useViewport();
  const [authScreen, setAuthScreen] = useState('login');
  const [pendingEmail, setPendingEmail] = useState('');
  const [booting, setBooting] = useState(true);
  const [confirmed, setConfirmed] = useState(false);

  useEffect(() => {
    const id = setTimeout(() => setBooting(false), 1150);
    return () => clearTimeout(id);
  }, []);

  useEffect(() => {
    const params = new URLSearchParams(window.location.search);
    if (params.get('status') === 'success' || window.location.pathname.includes('email-confirmed')) {
      setConfirmed(true);
    }
  }, []);

  // Capture OAuth tokens redirected back from the backend in the URL fragment.
  useEffect(() => {
    if (!window.location.hash || window.location.hash.length < 2) return;
    const frag = new URLSearchParams(window.location.hash.slice(1));
    const at = frag.get('access_token');
    const rt = frag.get('refresh_token');
    if (!at) return;
    // Strip the fragment so tokens don't linger in the address bar.
    window.history.replaceState({}, '', window.location.pathname);
    auth.login(at, rt);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  async function handleChangePassword(current, next) {
    if (!auth.user) throw new Error('Not signed in');
    await api.changePassword(auth.user.id, current, next);
  }
  async function handleDeleteAccount() {
    if (!auth.user) return;
    try { await api.deleteUser(auth.user.id); } catch (_) {}
    await auth.logout();
  }
  async function handleUpdateUser(patch) {
    if (!auth.user?.id) throw new Error('Not signed in');
    const upd = await api.patchUser(auth.user.id, patch);
    auth.updateUser(upd || patch);
  }

  if (!auth.authed) {
    let inner;
    if (confirmed) inner = <EmailConfirmed onAuth={() => { setConfirmed(false); setAuthScreen('login'); }} />;
    else if (authScreen === 'signup') inner = <SignUp go={setAuthScreen} onPending={(em) => { setPendingEmail(em); setAuthScreen('confirm'); }} />;
    else if (authScreen === 'confirm') inner = <EmailConfirm email={pendingEmail} go={setAuthScreen} />;
    else inner = <Login go={setAuthScreen} onAuth={auth.login} />;
    return (
      <div className="app-root">
        <div className="auth-stage">{inner}</div>
        <BootSplash visible={booting} />
      </div>
    );
  }

  // Wait for the user record to load before mounting the shell — otherwise
  // the profile/header would briefly render with placeholder values.
  if (!auth.user) {
    return <div className="app-root"><BootSplash visible={true} /></div>;
  }

  const Shell = isMobile ? MobileApp : DesktopApp;
  return (
    <div className="app-root">
      <Shell
        user={auth.user}
        tasks={taskOps.tasks}
        theme={theme}
        setTheme={setTheme}
        onLogout={auth.logout}
        onUpdateUser={handleUpdateUser}
        onChangePassword={handleChangePassword}
        onDeleteAccount={handleDeleteAccount}
        taskOps={taskOps}
      />
      <BootSplash visible={booting} />
    </div>
  );
}
