'use client';

import React, { useContext, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { CommonContext } from '@/components/user';
import { DomainControl } from '@/components/domain-control';
import { ToastContainer } from 'react-toastify';
import { get_info } from '@/api/info';
import 'react-toastify/dist/ReactToastify.css';
import { ProviderControl } from '@/components/provider';

export default function Domains() {
	const router = useRouter();
	const userStore = useContext(CommonContext);
	const { credentials, setCurrentPage, saveVersion } = userStore;

	useEffect(() => {
		if (!credentials) {
			router.push('/login');
			return;
		}
		setCurrentPage('Domains');
		get_info(credentials).then((info) => {
			saveVersion(info.version);
		});

	}, [setCurrentPage, credentials, saveVersion, router]);

	return (
		<main className="flex min-h-screen flex-col items-center justify-start pt-10 max-w-screen-xl">
			<ToastContainer />
			<div className="flex flex-col items-center w-full bg-base-100 p-10">
				{/* <div className="flex flex-row items-start w-full"> */}
				{/* <span className="text-2xl font-semibold text-neutral-500 ml-1 mb-1">Provider</span> */}
				{/* </div> */}
				<ProviderControl />
				<div className="divider"></div>
				<DomainControl />
			</div>
		</main>
	);
};