'use client';
import React from 'react';
import { LogEntry } from '@/api/logs';

interface LogEntryProps {
  log: LogEntry;
}

const getLevelColor = (level: string): string => {
  switch (level.toLowerCase()) {
    case 'error':
    case 'fatal':
    case 'panic':
      return 'text-error';
    case 'warn':
      return 'text-warning';
    case 'info':
      return 'text-info';
    case 'debug':
      return 'text-primary';
    case 'trace':
      return 'text-neutral-500';
    default:
      return 'text-base-content';
  }
};

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
    <div className="border-b border-base-200 py-2 px-4 hover:bg-base-50 font-mono text-sm">
      <div className="flex items-start gap-3">
        <span className="text-neutral-500 text-xs whitespace-nowrap">
          {formatTimestamp(log.timestamp)}
        </span>
        <span className={`badge badge-sm ${getLevelBadgeClass(log.level)} text-xs uppercase`}>
          {log.level}
        </span>
        <span className={`flex-1 ${getLevelColor(log.level)}`}>
          {log.message}
        </span>
      </div>
      {log.fields && Object.keys(log.fields).length > 0 && (
        <div className="mt-2 ml-16 text-xs text-neutral-600">
          <details className="cursor-pointer">
            <summary className="text-neutral-500">Fields ({Object.keys(log.fields).length})</summary>
            <div className="mt-1 ml-4">
              {Object.entries(log.fields).map(([key, value]) => (
                <div key={key} className="flex gap-2">
                  <span className="text-neutral-400">{key}:</span>
                  <span className="text-neutral-600">
                    {typeof value === 'object' ? JSON.stringify(value) : String(value)}
                  </span>
                </div>
              ))}
            </div>
          </details>
        </div>
      )}
    </div>
  );
};