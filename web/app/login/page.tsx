'use client';
// components/Login.tsx
import React, { useState } from 'react';
import { login } from '@/api/login';
import { ToastContainer, toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';

export default function Login() {
  const [username, setUsername] = useState<string>('');
  const [password, setPassword] = useState<string>('');

  const handleLogin = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    // Handle login logic here
    login(username, password).then((success) => {
      if (!success) {
        toast.error('Invalid username or password!', {
          position: "top-right",
          autoClose: 3000,
          hideProgressBar: false,
          closeOnClick: true,
          pauseOnHover: true,
          draggable: true,
          progress: undefined,
          theme: "light",
        });
      }
    });
  };

  return (
    <main className="flex min-h-screen flex-col items-center justify-center">
      <ToastContainer />
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
                <input type="text" id="username" placeholder="Type here" className="input input-primary input-bordered w-full max-w-xs"
                  onChange={
                    (e) => setUsername(e.target.value)
                  } />
              </label>
            </div>
            <div className="mb-4">
              <label className="form-control w-full max-w-xs">
                <div className="label">
                  <span className="label-text font-bold">Password</span>
                </div>
                <input type="password" id="password" placeholder="Type here" className="input input-primary input-bordered w-full max-w-xs"
                  onChange={
                    (e) => setPassword(e.target.value)
                  } />
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