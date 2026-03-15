'use client';

import React, { useState, useContext } from 'react';
import { useRouter } from 'next/navigation';
import { CommonContext } from '@/components/user';
import { login } from '@/api/login';
import { toast } from 'react-toastify';

export default function Login() {
  const router = useRouter();
  const [username, setUsername] = useState<string>('');
  const [password, setPassword] = useState<string>('');
  const [submitting, setSubmitting] = useState(false);
  const { loginUser } = useContext(CommonContext);

  const handleLogin = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setSubmitting(true);

    try {
      const credentials = await login(username, password);
      if (!credentials) {
        toast.error('Invalid username or password.');
        return;
      }

      loginUser(credentials);
      router.push('/');
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <main className="mx-auto flex min-h-[calc(100vh-12rem)] w-full max-w-5xl items-center justify-center">
      <div className="grid w-full gap-6 lg:grid-cols-[1.15fr_0.85fr]">
        <section className="panel flex flex-col justify-between gap-8 p-8 sm:p-10">
          <div className="space-y-5">
            <span className="hero-badge">Welcome</span>
            <div className="space-y-3">
              <h1 className="page-title">Manage GoDNS without editing config files by hand.</h1>
              <p className="page-subtitle">
                Sign in to configure provider credentials, attach domains, tune network detection, and monitor update logs from one admin surface.
              </p>
            </div>
          </div>

          <div className="grid gap-4 sm:grid-cols-3">
            <div className="panel-muted">
              <p className="text-sm font-semibold">1. Connect providers</p>
              <p className="mt-2 text-sm text-base-content/65">Store the credentials GoDNS needs for each DNS provider you use.</p>
            </div>
            <div className="panel-muted">
              <p className="text-sm font-semibold">2. Attach domains</p>
              <p className="mt-2 text-sm text-base-content/65">Define domains and subdomains so updates map cleanly to provider profiles.</p>
            </div>
            <div className="panel-muted">
              <p className="text-sm font-semibold">3. Verify runtime</p>
              <p className="mt-2 text-sm text-base-content/65">Inspect IP mode, networking, and logs to confirm updates behave as expected.</p>
            </div>
          </div>
        </section>

        <section className="panel p-8 sm:p-10">
          <form onSubmit={handleLogin} className="space-y-6">
            <div className="space-y-2">
              <h2 className="text-2xl font-semibold tracking-tight">Sign in</h2>
              <p className="text-sm text-base-content/60">
                Use the GoDNS web credentials configured on the server.
              </p>
            </div>

            <label className="field-stack">
              <span className="field-label">Username</span>
              <input
                type="text"
                id="username"
                placeholder="Enter your username"
                className="input input-bordered h-12 w-full rounded-2xl"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
              />
            </label>

            <label className="field-stack">
              <span className="field-label">Password</span>
              <input
                type="password"
                id="password"
                placeholder="Enter your password"
                className="input input-bordered h-12 w-full rounded-2xl"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
              />
            </label>

            <button className="btn btn-primary h-12 w-full rounded-full" disabled={submitting || !username || !password}>
              {submitting ? 'Signing in...' : 'Sign in'}
            </button>
          </form>
        </section>
      </div>
    </main>
  );
}
