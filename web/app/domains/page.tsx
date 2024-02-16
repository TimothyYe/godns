'use client';
// components/Login.tsx
import React, { useContext, useEffect } from 'react';
import { CommonContext } from '@/components/user';
import { DomainControl } from '@/components/domain-control';
import { ToastContainer } from 'react-toastify';
import { get_info } from '@/api/info';
import 'react-toastify/dist/ReactToastify.css';

export default function Domains() {
	const userStore = useContext(CommonContext);
	const { credentials, setCurrentPage, saveVersion } = userStore;

	useEffect(() => {
		if (!credentials) {
			window.location.href = '/login';
			return;
		}
		setCurrentPage('Domains');
		get_info(credentials).then((info) => {
			saveVersion(info.version);
		});

	}, [setCurrentPage, credentials, saveVersion]);

	return (
		<main className="flex min-h-screen max-w-screen-xl flex-col">
			<ToastContainer />
			<div className="card w-auto bg-base-100">
				<DomainControl />
			</div>
		</main>
	);
};