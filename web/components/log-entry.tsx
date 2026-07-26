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
      return 'theme-log-level-danger';
    case 'warn':
      return 'theme-log-level-warn';
    case 'info':
      return 'theme-log-level-info';
    case 'debug':
      return 'theme-log-level-debug';
    case 'trace':
      return 'theme-log-level-trace';
    default:
      return 'theme-log-level-default';
  }
};

const getLevelBadgeClass = (level: string): string => {
  switch (level.toLowerCase()) {
    case 'error':
    case 'fatal':
    case 'panic':
      return 'theme-badge-danger';
    case 'warn':
      return 'theme-badge-amber';
    case 'info':
      return 'theme-badge-sky';
    case 'debug':
      return 'theme-badge-violet';
    case 'trace':
      return 'theme-badge-neutral';
    default:
      return 'theme-badge-neutral';
  }
};

export const LogEntryComponent: React.FC<LogEntryProps> = ({ log }) => {
  const formatTimestamp = (timestamp: string) => {
    return new Date(timestamp).toLocaleString();
  };

  return (
    <div className="log-row px-4 py-2 text-sm">
      <div className="flex items-start gap-3">
        <span className="log-time whitespace-nowrap text-xs">
          {formatTimestamp(log.timestamp)}
        </span>
        <span className={`badge badge-sm border ${getLevelBadgeClass(log.level)} text-xs uppercase`}>
          {log.level}
        </span>
        <span className={`log-message flex-1 ${getLevelColor(log.level)}`}>
          {log.message}
        </span>
      </div>
      {log.fields && Object.keys(log.fields).length > 0 && (
        <div className="log-fields mt-2 ml-16 text-xs">
          <details className="cursor-pointer">
            <summary className="log-fields-summary">Fields ({Object.keys(log.fields).length})</summary>
            <div className="mt-1 ml-4">
              {Object.entries(log.fields).map(([key, value]) => (
                <div key={key} className="flex gap-2">
                  <span className="log-field-key">{key}:</span>
                  <span className="log-field-value">
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
