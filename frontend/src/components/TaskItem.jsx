import React, { useState } from 'react';
import { CheckIcon, TrashIcon } from './Icons';
import { api } from '../api/client';

function getPriorityTag(task) {
  if (task.completed) return null;
  const id = task.id % 3;
  if (id === 0) return { cls: 'tag-high', label: 'High' };
  if (id === 1) return { cls: 'tag-med',  label: 'Med'  };
  return null;
}

export default function TaskItem({ task, onUpdate, onDelete, style }) {
  const [loading, setLoading] = useState(false);

  async function toggleComplete() {
    setLoading(true);
    try {
      const updated = await api.patchTask(task.id, { completed: !task.completed });
      onUpdate?.(updated);
    } finally {
      setLoading(false);
    }
  }

  async function handleDelete() {
    setLoading(true);
    try {
      await api.deleteTask(task.id);
      onDelete?.(task.id);
    } finally {
      setLoading(false);
    }
  }

  const tag = getPriorityTag(task);

  return (
    <div
      className="task-item"
      style={{ opacity: loading ? .5 : 1, ...style }}
    >
      <button
        className={`task-check ${task.completed ? 'done' : ''}`}
        onClick={toggleComplete}
        disabled={loading}
        aria-label={task.completed ? 'Mark incomplete' : 'Mark complete'}
      >
        {task.completed && <CheckIcon size={11} />}
      </button>

      <span className={`task-title ${task.completed ? 'done' : ''}`}>
        {task.title}
      </span>

      <div className="task-tags">
        {task.completed && <span className="tag tag-done">Done</span>}
        {!task.completed && tag && <span className={`tag ${tag.cls}`}>{tag.label}</span>}
      </div>

      <button
        onClick={handleDelete}
        disabled={loading}
        style={{ color: 'var(--text-3)', padding: '4px 6px', marginLeft: 2, borderRadius: 8,
                 transition: 'background .15s, color .15s' }}
        onMouseEnter={(e) => { e.currentTarget.style.background = 'var(--tag-high-bg)'; e.currentTarget.style.color = 'var(--tag-high)'; }}
        onMouseLeave={(e) => { e.currentTarget.style.background = 'transparent'; e.currentTarget.style.color = 'var(--text-3)'; }}
        aria-label="Delete task"
      >
        <TrashIcon size={15} />
      </button>
    </div>
  );
}
