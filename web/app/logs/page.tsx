'use client';

import Link from 'next/link';
import React, { useState, useEffect, useContext, useRef, useCallback, useMemo } from 'react';
import { useRouter } from 'next/navigation';
import { CommonContext } from '@/components/user';
import { get_logs, clear_logs, get_log_levels, LogEntry } from '@/api/logs';
import { LogEntryComponent } from '@/components/log-entry';
import { toast } from 'react-toastify';
import { PageShell, SectionCard } from '@/components/page-shell';

const refreshOptions = [
  { label: '2s', value: 2000 },
  { label: '5s', value: 5000 },
  { label: '15s', value: 15000 },
  { label: '30s', value: 30000 },
];

export default function LogsPage() {
  const router = useRouter();
  const userStore = useContext(CommonContext);
  const { credentials, setCurrentPage } = userStore;

  const [logs, setLogs] = useState<LogEntry[]>([]);
  const [filteredLogs, setFilteredLogs] = useState<LogEntry[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedLevel, setSelectedLevel] = useState<string>('');
  const [searchTerm, setSearchTerm] = useState<string>('');
  const [logLevels, setLogLevels] = useState<string[]>([]);
  const [autoRefresh, setAutoRefresh] = useState(false);
  const [refreshInterval, setRefreshInterval] = useState(5000);
  const intervalRef = useRef<NodeJS.Timeout | null>(null);
  const logsEndRef = useRef<HTMLDivElement>(null);

  const loadLogs = useCallback(async (showLoading = true) => {
    if (!credentials) return;

    if (showLoading) setLoading(true);

    try {
      const response = await get_logs(credentials, 500);
      if (response) {
        setLogs(response.logs);
      }
    } catch (error) {
      toast.error('Failed to load logs');
    } finally {
      if (showLoading) setLoading(false);
    }
  }, [credentials]);

  const loadLogLevels = useCallback(async () => {
    if (!credentials) return;

    try {
      const levels = await get_log_levels(credentials);
      setLogLevels(levels);
    } catch (error) {
      console.error('Failed to load log levels:', error);
    }
  }, [credentials]);

  useEffect(() => {
    if (!credentials) {
      router.push('/login');
      return;
    }
    setCurrentPage('Logs');
    loadLogs();
    loadLogLevels();
  }, [credentials, router, setCurrentPage, loadLogs, loadLogLevels]);

  useEffect(() => {
    let filtered = logs;

    if (selectedLevel) {
      filtered = filtered.filter((log) => log.level === selectedLevel);
    }

    if (searchTerm) {
      const searchLower = searchTerm.toLowerCase();
      filtered = filtered.filter((log) =>
        log.message.toLowerCase().includes(searchLower)
        || log.level.toLowerCase().includes(searchLower)
        || (log.fields && JSON.stringify(log.fields).toLowerCase().includes(searchLower))
      );
    }

    setFilteredLogs(filtered);
  }, [logs, selectedLevel, searchTerm]);

  useEffect(() => {
    if (autoRefresh && credentials) {
      intervalRef.current = setInterval(() => {
        loadLogs(false);
      }, refreshInterval);
    } else if (intervalRef.current) {
      clearInterval(intervalRef.current);
      intervalRef.current = null;
    }

    return () => {
      if (intervalRef.current) {
        clearInterval(intervalRef.current);
      }
    };
  }, [autoRefresh, refreshInterval, credentials, loadLogs]);

  const handleClearLogs = async () => {
    if (!credentials) return;

    try {
      const success = await clear_logs(credentials);
      if (success) {
        setLogs([]);
        toast.success('Logs cleared successfully');
      } else {
        toast.error('Failed to clear logs');
      }
    } catch (error) {
      toast.error('Failed to clear logs');
    }
  };

  const levelCounts = useMemo(() => {
    return logs.reduce<Record<string, number>>((acc, log) => {
      const key = log.level.toLowerCase();
      acc[key] = (acc[key] || 0) + 1;
      return acc;
    }, {});
  }, [logs]);

  return (
    <PageShell
      eyebrow="Observability"
      title="Runtime logs"
      description="Inspect recent GoDNS activity, filter by severity, and keep a live tail open while you validate configuration changes."
      actions={(
        <>
          <Link className="btn btn-ghost rounded-full px-5" href="/network">Back to network</Link>
          <button className="btn btn-primary rounded-full px-5" onClick={() => loadLogs()} disabled={loading}>
            {loading ? 'Refreshing...' : 'Refresh'}
          </button>
        </>
      )}
    >
      <SectionCard
        title="Filters and controls"
        description="Narrow the current view or switch the log stream into live-refresh mode."
        actions={(
          <button className="btn btn-error btn-sm rounded-full px-4" onClick={handleClearLogs}>
            Clear logs
          </button>
        )}
      >
        <div className="grid gap-5 lg:grid-cols-[1.3fr_0.7fr]">
          <div className="grid gap-4 sm:grid-cols-2">
            <label className="field-stack">
              <span className="field-label">Search logs</span>
              <input
                type="text"
                className="input input-bordered h-12 w-full rounded-2xl"
                placeholder="Message text, severity, or structured field"
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
              />
            </label>

            <label className="field-stack">
              <span className="field-label">Severity</span>
              <select
                className="select select-bordered h-12 w-full rounded-2xl"
                value={selectedLevel}
                onChange={(e) => setSelectedLevel(e.target.value)}
              >
                <option value="">All levels</option>
                {logLevels.map((level) => (
                  <option key={level} value={level}>{level.toUpperCase()}</option>
                ))}
              </select>
            </label>
          </div>

          <div className="grid gap-4 rounded-[1.5rem] border border-base-300/70 p-4">
            <label className="flex items-center justify-between rounded-[1.25rem] bg-base-200/50 px-4 py-3">
              <div className="space-y-1">
                <p className="text-sm font-semibold">Live refresh</p>
                <p className="text-xs text-base-content/60">{autoRefresh ? 'Streaming the latest logs.' : 'Manual refresh only.'}</p>
              </div>
              <input
                type="checkbox"
                className="toggle toggle-primary"
                checked={autoRefresh}
                onChange={(e) => setAutoRefresh(e.target.checked)}
              />
            </label>

            <label className="field-stack">
              <span className="field-label">Refresh cadence</span>
              <select
                className="select select-bordered h-12 w-full rounded-2xl"
                value={refreshInterval}
                disabled={!autoRefresh}
                onChange={(e) => setRefreshInterval(Number(e.target.value))}
              >
                {refreshOptions.map((option) => (
                  <option key={option.value} value={option.value}>{option.label}</option>
                ))}
              </select>
            </label>

            <div className="flex gap-2">
              <button className="btn btn-ghost btn-sm rounded-full" onClick={() => {
                setSelectedLevel('');
                setSearchTerm('');
              }}>
                Clear filters
              </button>
              <button className="btn btn-ghost btn-sm rounded-full" onClick={() => logsEndRef.current?.scrollIntoView({ behavior: 'smooth' })}>
                Jump to latest
              </button>
            </div>
          </div>
        </div>
      </SectionCard>

      <section className="grid gap-4 sm:grid-cols-2 xl:grid-cols-5">
        <div className="metric-card xl:col-span-2">
          <p className="metric-label">Visible entries</p>
          <p className="metric-value">{filteredLogs.length}</p>
          <p className="metric-meta">Showing {filteredLogs.length} of {logs.length} buffered log entries.</p>
        </div>
        {['error', 'warn', 'info'].map((level) => (
          <div key={level} className="metric-card">
            <p className="metric-label">{level.toUpperCase()}</p>
            <p className="metric-value">{levelCounts[level] || 0}</p>
            <p className="metric-meta">Entries at {level} level.</p>
          </div>
        ))}
      </section>

      <SectionCard
        title="Log stream"
        description="Recent GoDNS events are buffered here. Expand structured fields to inspect provider responses or update payloads."
      >
        {loading ? (
          <div className="flex min-h-80 items-center justify-center">
            <span className="loading loading-spinner loading-lg" />
          </div>
        ) : filteredLogs.length === 0 ? (
          <div className="panel-muted">
            <p className="font-semibold">{logs.length === 0 ? 'No logs available' : 'No logs match your filters'}</p>
            <p className="mt-2 text-sm text-base-content/60">
              {logs.length === 0 ? 'Generate activity or confirm logging is enabled in GoDNS.' : 'Reset the search term or severity filter to broaden the results.'}
            </p>
          </div>
        ) : (
          <div className="overflow-hidden rounded-[1.5rem] border border-base-300/70 bg-base-100">
            {filteredLogs.map((log, index) => (
              <LogEntryComponent key={`${log.timestamp}-${index}`} log={log} />
            ))}
            <div ref={logsEndRef} />
          </div>
        )}
      </SectionCard>
    </PageShell>
  );
}
