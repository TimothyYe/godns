'use client';
// components/Login.tsx
import React, { useState } from 'react';

export default function Login() {
  const [username, setUsername] = useState<string>('');
  const [password, setPassword] = useState<string>('');

  const handleLogin = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    // Handle login logic here
    console.log('Login attempt with:', username, password);
  };

  return (
    <main className="flex min-h-screen flex-col items-center justify-center">
      <div className="card w-96 shadow-2xl shadow-neutral-950">
        <div className="card-body">
          <form onSubmit={handleLogin} className="flex flex-col mb-4">
            <h2 className="card-title text-primary">Login</h2>
            <div className="divider" />
            <div className="mb-4">
              <label className="form-control w-full max-w-xs">
                <div className="label">
                  <span className="label-text font-bold">Username</span>
                </div>
                <input type="text" id="username" placeholder="Type here" className="input input-primary input-bordered w-full max-w-xs" />
              </label>
            </div>
            <div className="mb-4">
              <label className="form-control w-full max-w-xs">
                <div className="label">
                  <span className="label-text font-bold">Password</span>
                </div>
                <input type="password" id="password" placeholder="Type here" className="input input-primary input-bordered w-full max-w-xs" />
              </label>
            </div>
            <div className="card-actions justify-end">
              <button className="btn btn-primary">Sign In</button>
            </div>
          </form>
        </div>
      </div>
    </main>
  );
};