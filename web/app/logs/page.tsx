'use client';
import React, { useState, useEffect, useContext, useRef, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { CommonContext } from '@/components/user';
import { get_logs, clear_logs, get_log_levels, LogEntry, LogsResponse } from '@/api/logs';
import { LogEntryComponent } from '@/components/log-entry';
import { toast } from 'react-toastify';

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
  const logsScrollRef = useRef<HTMLDivElement>(null);

  const loadLogs = useCallback(async (showLoading = true) => {
    if (!credentials) return;
    
    if (showLoading) setLoading(true);
    
    try {
      const response = await get_logs(credentials, 500); // Get last 500 logs
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
    
    // Load initial data
    loadLogs();
    loadLogLevels();
  }, [credentials, router, setCurrentPage, loadLogs, loadLogLevels]);

  useEffect(() => {
    // Filter logs based on level and search term
    let filtered = logs;
    
    if (selectedLevel) {
      filtered = filtered.filter(log => log.level === selectedLevel);
    }
    
    if (searchTerm) {
      const searchLower = searchTerm.toLowerCase();
      filtered = filtered.filter(log => 
        log.message.toLowerCase().includes(searchLower) ||
        log.level.toLowerCase().includes(searchLower) ||
        (log.fields && JSON.stringify(log.fields).toLowerCase().includes(searchLower))
      );
    }
    
    setFilteredLogs(filtered);
  }, [logs, selectedLevel, searchTerm]);

  useEffect(() => {
    // Auto-refresh functionality
    if (autoRefresh && credentials) {
      intervalRef.current = setInterval(() => {
        loadLogs(false); // Don't show loading for auto-refresh
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

  const scrollToBottom = () => {
    logsScrollRef.current?.scrollTo({ top: logsScrollRef.current.scrollHeight, behavior: 'smooth' });
  };

  const scrollToTop = () => {
    logsScrollRef.current?.scrollTo({ top: 0, behavior: 'smooth' });
  };

  if (loading) {
    return (
      <div className="surface-panel flex min-h-[18rem] items-center justify-center">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    );
  }

  return (
    <main className="page-wrap">
      <div className="page-shell">
		<section className="page-hero page-hero-compact">
		  <div className="section-header">
		    <div className="section-copy">
		      <div className="eyebrow">
		        <span className="inline-block h-2 w-2 rounded-full bg-emerald-400" />
		        Observability
		      </div>
		      <h1 className="page-title">Application logs at a glance.</h1>
		    </div>
		    <div className="flex flex-wrap gap-2">
              <button
                className="theme-primary-sky btn btn-sm rounded-xl border-none px-4"
                onClick={() => loadLogs()}
                disabled={loading}
              >
                {loading ? <span className="loading loading-spinner loading-xs"></span> : 'Refresh'}
              </button>
              <button
                className="theme-danger btn btn-sm rounded-xl border-none px-4"
                onClick={handleClearLogs}
              >
                Clear Logs
              </button>
            </div>
          </div>
        </section>

        <section className="section-shell">
          <div className="section-header">
            <div className="section-copy">
              <div className="section-label ml-0">Filters</div>
              <h2 className="text-2xl font-semibold tracking-tight theme-heading">Refine the stream</h2>
              <p className="mt-2 text-sm leading-7 theme-muted">
                Search by message content, narrow by severity, and optionally keep the stream polling in the background.
              </p>
            </div>
          </div>

          <div className="mt-5 grid gap-4 lg:grid-cols-[minmax(0,1.6fr)_minmax(14rem,0.7fr)_minmax(12rem,0.7fr)]">
            <fieldset className="theme-field">
              <label className="theme-field-label" htmlFor="log-search">Search logs</label>
              <input
                id="log-search"
                type="text"
                className="input theme-input w-full rounded-2xl"
                placeholder="Search messages, levels, or structured fields"
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
              />
            </fieldset>

            <fieldset className="theme-field">
              <label className="theme-field-label" htmlFor="log-level">Filter by level</label>
              <select
                id="log-level"
                className="select theme-input w-full rounded-2xl"
                value={selectedLevel}
                onChange={(e) => setSelectedLevel(e.target.value)}
              >
                <option value="">All levels</option>
                {logLevels.map((level) => (
                  <option key={level} value={level}>{level.toUpperCase()}</option>
                ))}
              </select>
            </fieldset>

            <fieldset className="theme-field">
              <label className="theme-field-label" htmlFor="refresh-interval">Refresh cadence</label>
              <input
                id="refresh-interval"
                type="number"
                className="input theme-input w-full rounded-2xl"
                value={refreshInterval}
                onChange={(e) => setRefreshInterval(parseInt(e.target.value, 10) || 5000)}
                min="1000"
                step="1000"
                disabled={!autoRefresh}
              />
              <label className="mt-2 inline-flex items-center gap-3 text-sm theme-muted" htmlFor="auto-refresh">
                <input
                  id="auto-refresh"
                  type="checkbox"
                  className="toggle toggle-primary"
                  checked={autoRefresh}
                  onChange={(e) => setAutoRefresh(e.target.checked)}
                />
                Auto-refresh logs
              </label>
            </fieldset>
          </div>
        </section>

        <section className="section-shell">
          <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
            <p className="text-sm theme-muted">
              Showing <span className="theme-heading">{filteredLogs.length}</span> of <span className="theme-heading">{logs.length}</span> log entries
            </p>
            <div className="flex gap-2">
              <button className="theme-subtle-btn btn btn-xs rounded-xl" onClick={scrollToTop}>
                ↑ Top
              </button>
              <button className="theme-subtle-btn btn btn-xs rounded-xl" onClick={scrollToBottom}>
                ↓ Bottom
              </button>
            </div>
          </div>

          <div className="log-stream mt-5 rounded-[1.5rem]">
            {filteredLogs.length === 0 ? (
              <div className="px-6 py-12 text-center">
                <p className="theme-muted">
                  {logs.length === 0 ? 'No logs available' : 'No logs match your filters'}
                </p>
              </div>
            ) : (
              <div className="log-stream-scroll" ref={logsScrollRef}>
                {filteredLogs.map((log, index) => (
                  <LogEntryComponent key={index} log={log} />
                ))}
                <div ref={logsEndRef} />
              </div>
            )}
          </div>
        </section>
      </div>
    </main>
  );
}
