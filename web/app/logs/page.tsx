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
    logsEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  const scrollToTop = () => {
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center min-h-screen">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-6 max-w-7xl">
      <div className="flex flex-col gap-4">
        {/* Header */}
        <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
          <div>
            <h1 className="text-2xl font-bold text-base-content">Application Logs</h1>
            <p className="text-base-content/60">View and monitor application logs in real-time</p>
          </div>
          
          <div className="flex flex-wrap gap-2">
            <button 
              className="btn btn-primary btn-sm"
              onClick={() => loadLogs()}
              disabled={loading}
            >
              {loading ? <span className="loading loading-spinner loading-xs"></span> : 'Refresh'}
            </button>
            <button 
              className="btn btn-error btn-sm"
              onClick={handleClearLogs}
            >
              Clear Logs
            </button>
          </div>
        </div>

        {/* Controls */}
        <div className="card bg-base-100 shadow-sm">
          <div className="card-body p-4">
            <div className="flex flex-col sm:flex-row gap-4">
              {/* Search */}
              <div className="flex-1">
                <label className="form-control w-full">
                  <div className="label">
                    <span className="label-text text-sm">Search logs</span>
                  </div>
                  <input
                    type="text"
                    className="input input-bordered input-sm w-full"
                    placeholder="Search in messages, levels, or fields..."
                    value={searchTerm}
                    onChange={(e) => setSearchTerm(e.target.value)}
                  />
                </label>
              </div>
              
              {/* Level Filter */}
              <div className="sm:w-48">
                <label className="form-control w-full">
                  <div className="label">
                    <span className="label-text text-sm">Filter by level</span>
                  </div>
                  <select 
                    className="select select-bordered select-sm w-full"
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
            </div>

            {/* Auto-refresh controls */}
            <div className="divider my-2"></div>
            <div className="flex flex-col sm:flex-row gap-4 items-start sm:items-end">
              <div className="form-control">
                <label className="label cursor-pointer gap-2">
                  <input 
                    type="checkbox" 
                    className="checkbox checkbox-sm"
                    checked={autoRefresh}
                    onChange={(e) => setAutoRefresh(e.target.checked)}
                  />
                  <span className="label-text text-sm">Auto-refresh</span>
                </label>
              </div>
              
              {autoRefresh && (
                <div className="sm:w-32">
                  <label className="form-control w-full">
                    <div className="label">
                      <span className="label-text text-xs">Interval (ms)</span>
                    </div>
                    <input
                      type="number"
                      className="input input-bordered input-sm w-full"
                      value={refreshInterval}
                      onChange={(e) => setRefreshInterval(parseInt(e.target.value) || 5000)}
                      min="1000"
                      step="1000"
                    />
                  </label>
                </div>
              )}
            </div>
          </div>
        </div>

        {/* Log Count and Actions */}
        <div className="flex justify-between items-center">
          <span className="text-sm text-base-content/60">
            Showing {filteredLogs.length} of {logs.length} log entries
          </span>
          <div className="flex gap-2">
            <button className="btn btn-ghost btn-xs" onClick={scrollToTop}>
              ↑ Top
            </button>
            <button className="btn btn-ghost btn-xs" onClick={scrollToBottom}>
              ↓ Bottom
            </button>
          </div>
        </div>

        {/* Logs Display */}
        <div className="card bg-base-100 shadow-sm">
          <div className="card-body p-0">
            {filteredLogs.length === 0 ? (
              <div className="text-center py-12">
                <p className="text-base-content/60">
                  {logs.length === 0 ? 'No logs available' : 'No logs match your filters'}
                </p>
              </div>
            ) : (
              <div className="max-h-[600px] overflow-y-auto">
                {filteredLogs.map((log, index) => (
                  <LogEntryComponent key={index} log={log} />
                ))}
                <div ref={logsEndRef} />
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}