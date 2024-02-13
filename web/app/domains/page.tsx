'use client';
// components/Login.tsx
import React, { useState, useContext, useEffect } from 'react';
import { CommonContext } from '@/components/user';
import { login } from '@/api/login';
import { ToastContainer, toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';

export default function Domains() {
	const [username, setUsername] = useState<string>('');
	const [password, setPassword] = useState<string>('');
	const userStore = useContext(CommonContext);
	const { credentials, setCurrentPage } = userStore;

	const handleSave = (e: React.FormEvent<HTMLFormElement>) => {
		e.preventDefault();
		console.log('Save button clicked');
	};

	useEffect(() => {
		setCurrentPage('Domains');
	}, [setCurrentPage]);

	return (
		<main className="flex min-h-screen max-w-screen-xl flex-col items-center justify-center">
			<ToastContainer />
			<div className="card w-auto shadow-2xl bg-base-100">
				<form onSubmit={handleSave} className="flex flex-col mb-4">
					<h2 className="card-title text-primary">Login</h2>
					<div className="divider" />
					<div className="mb-4">
						<label className="form-control w-full max-w-xs">
							<div className="label">
								<span className="label-text font-bold">Username</span>
							</div>
							<input type="text" id="username" placeholder="Input the username" className="input input-primary input-bordered w-full max-w-xs"
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
							<input type="password" id="password" placeholder="Input the password" className="input input-primary input-bordered w-full max-w-xs"
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
		</main>
	);
};