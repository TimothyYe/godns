'use client';
// components/Login.tsx
import React, { useState, useContext, useEffect } from 'react';
import { CommonContext } from '@/components/user';
import { DomainControl } from '@/components/domain-control';
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
		<main className="flex min-h-screen max-w-screen-xl flex-col">
			<ToastContainer />
			<div className="card w-auto bg-base-100">
				<form onSubmit={handleSave}>
					<DomainControl />
				</form>
			</div>
		</main>
	);
};