import React, { useState, useEffect, useCallback } from 'react';
import { useAuth } from '../App';
import { api } from '../api/client';
import AuthLayout from '../components/AuthLayout';
import TaskItem from '../components/TaskItem';
import AddTaskModal from '../components/AddTaskModal';
import { SearchIcon } from '../components/Icons';

function groupTasks(tasks) {
  const now     = new Date();
  const today   = new Date(now.getFullYear(), now.getMonth(), now.getDate());
  const weekAgo = new Date(today); weekAgo.setDate(today.getDate() - 6);

  const groups = { Today: [], 'This Week': [], Earlier: [] };
  for (const t of tasks) {
    const d   = new Date(t.created_at);
    const day = new Date(d.getFullYear(), d.getMonth(), d.getDate());
    if (day >= today)   groups['Today'].push(t);
    else if (day >= weekAgo) groups['This Week'].push(t);
    else                groups['Earlier'].push(t);
  }
  return groups;
}

const TABS = ['All', 'Active', 'Done'];

export default function TasksPage() {
  const { logout }                    = useAuth();
  const [tasks,     setTasks]         = useState([]);
  const [loading,   setLoading]       = useState(true);
  const [error,     setError]         = useState('');
  const [showAdd,   setShowAdd]       = useState(false);
  const [search,    setSearch]        = useState('');
  const [activeTab, setActiveTab]     = useState('All');

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

  function handleUpdate(updated) {
    setTasks((ts) => ts.map((t) => (t.id === updated.id ? updated : t)));
  }
  function handleDelete(id) {
    setTasks((ts) => ts.filter((t) => t.id !== id));
  }
  function handleCreated(task) {
    setTasks((ts) => [task, ...ts]);
  }

  const filtered = tasks
    .filter((t) => {
      if (activeTab === 'Active') return !t.completed;
      if (activeTab === 'Done')   return  t.completed;
      return true;
    })
    .filter((t) =>
      !search || t.title.toLowerCase().includes(search.toLowerCase())
    );

  const groups = groupTasks(filtered);

  return (
    <AuthLayout title="My Tasks" onAdd={() => setShowAdd(true)}>
      {/* Search + tabs */}
      <div style={{ padding: '12px 20px 0', display: 'flex', gap: 10 }}>
        <label className="search-bar">
          <SearchIcon size={15} />
          <input
            placeholder="Search tasks…"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
        </label>
        <button
          className="btn btn-primary"
          style={{ padding: '10px 16px', fontSize: 13, flexShrink: 0 }}
          onClick={() => setShowAdd(true)}
        >
          + Add
        </button>
      </div>

      <div className="tab-bar">
        {TABS.map((tab) => (
          <button
            key={tab}
            className={`tab-btn ${activeTab === tab ? 'active' : 'inactive'}`}
            onClick={() => setActiveTab(tab)}
          >
            {tab}
          </button>
        ))}
      </div>

      {/* List */}
      {error && <div className="error-msg">{error}</div>}

      {loading ? (
        <div className="spinner" />
      ) : filtered.length === 0 ? (
        <div className="empty-state">
          <svg width="64" height="64" viewBox="0 0 64 64" fill="none">
            <rect x="8" y="12" width="48" height="40" rx="6" fill="var(--primary-light)" />
            <path d="M20 28h24M20 36h16" stroke="var(--primary-mid)" strokeWidth="2.5" strokeLinecap="round" />
          </svg>
          <p>No tasks yet. Tap <strong>+ Add</strong> to create one.</p>
        </div>
      ) : (
        Object.entries(groups).map(([label, list]) =>
          list.length === 0 ? null : (
            <div key={label}>
              <div className="section-label">{label}</div>
              {list.map((task, i) => (
                <TaskItem
                  key={task.id}
                  task={task}
                  onUpdate={handleUpdate}
                  onDelete={handleDelete}
                  style={{ animationDelay: `${i * 35}ms` }}
                />
              ))}
            </div>
          )
        )
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
