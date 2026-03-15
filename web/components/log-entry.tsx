'use client';

import React from 'react';
import { LogEntry } from '@/api/logs';

interface LogEntryProps {
  log: LogEntry;
}

const getLevelBadgeClass = (level: string): string => {
  switch (level.toLowerCase()) {
    case 'error':
    case 'fatal':
    case 'panic':
      return 'badge-error';
    case 'warn':
      return 'badge-warning';
    case 'info':
      return 'badge-info';
    case 'debug':
      return 'badge-primary';
    case 'trace':
      return 'badge-neutral';
    default:
      return 'badge-ghost';
  }
};

export const LogEntryComponent: React.FC<LogEntryProps> = ({ log }) => {
  const formatTimestamp = (timestamp: string) => {
    return new Date(timestamp).toLocaleString();
  };

  return (
    <article className="border-b border-base-300/70 px-4 py-4 last:border-b-0">
      <div className="flex flex-col gap-3 md:flex-row md:items-start">
        <div className="min-w-52 text-xs uppercase tracking-[0.2em] text-base-content/45">
          {formatTimestamp(log.timestamp)}
        </div>
        <div className="flex-1 space-y-3">
          <div className="flex flex-wrap items-start gap-3">
            <span className={`badge badge-sm ${getLevelBadgeClass(log.level)} uppercase`}>
              {log.level}
            </span>
            <p className="font-mono text-sm leading-6 text-base-content">
              {log.message}
            </p>
          </div>
          {log.fields && Object.keys(log.fields).length > 0 ? (
            <details className="rounded-2xl border border-base-300 bg-base-200/40 p-3">
              <summary className="cursor-pointer text-sm font-medium text-base-content/70">
                Structured fields ({Object.keys(log.fields).length})
              </summary>
              <div className="mt-3 space-y-2">
                {Object.entries(log.fields).map(([key, value]) => (
                  <div key={key} className="code-surface grid gap-1">
                    <span className="text-xs uppercase tracking-[0.2em] text-base-content/45">{key}</span>
                    <span className="break-all text-base-content/80">
                      {typeof value === 'object' ? JSON.stringify(value) : String(value)}
                    </span>
                  </div>
                ))}
              </div>
            </details>
          ) : null}
        </div>
      </div>
    </article>
  );
};
