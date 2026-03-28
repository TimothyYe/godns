'use client';
import React, { useState, useContext } from 'react';
import { useRouter } from 'next/navigation';
import { CommonContext } from '@/components/user';
import { login } from '@/api/login';
import { ToastContainer, toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';

export default function Login() {
  const router = useRouter();
  const [username, setUsername] = useState<string>('');
  const [password, setPassword] = useState<string>('');
  const { loginUser } = useContext(CommonContext);

  const handleLogin = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    // Handle login logic here
    login(username, password).then((credentials) => {
      if (!credentials) {
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
      } else {
        loginUser(credentials);
        // Redirect to the home page
        router.push('/');
      }
    });
  };

  return (
    <main className="flex min-h-[calc(100dvh-10rem)] items-center justify-center py-3 sm:min-h-[calc(100dvh-10.5rem)] sm:py-4">
      <ToastContainer />
		<div className="page-shell max-w-md gap-5">
		  <section className="page-hero page-hero-compact">
			<div className="eyebrow">
			  <span className="inline-block h-2 w-2 rounded-full bg-sky-400" />
			  Secure Access
			</div>
			<h1 className="page-title text-3xl sm:text-4xl">Sign in to GoDNS.</h1>
		  </section>

        <section className="section-shell">
          <form onSubmit={handleLogin} className="flex flex-col gap-5">
            <div>
              <h2 className="text-2xl font-semibold tracking-tight theme-heading">Login</h2>
              <p className="mt-2 text-sm leading-7 theme-muted">Use your web panel credentials to continue.</p>
            </div>

            <fieldset className="theme-field">
              <label className="theme-field-label" htmlFor="username">Username</label>
              <input
                type="text"
                id="username"
                placeholder="Input the username"
                className="input theme-input w-full rounded-2xl"
                onChange={(e) => setUsername(e.target.value)}
              />
            </fieldset>

            <fieldset className="theme-field">
              <label className="theme-field-label" htmlFor="password">Password</label>
              <input
                type="password"
                id="password"
                placeholder="Input the password"
                className="input theme-input w-full rounded-2xl"
                onChange={(e) => setPassword(e.target.value)}
              />
            </fieldset>

            <div className="flex justify-end">
              <button className="theme-primary-sky btn rounded-xl border-none px-6">Sign In</button>
            </div>
          </form>
        </section>
      </div>
    </main>
  );
};
