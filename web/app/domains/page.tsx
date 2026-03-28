'use client';

import React, { useContext, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { CommonContext } from '@/components/user';
import { DomainControl } from '@/components/domain-control';
import { get_info } from '@/api/info';
import { MultiProviderControl } from '@/components/multi-provider-control';

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
		<main className="page-wrap">
			<div className="page-shell">
				<section className="page-hero page-hero-compact">
					<div className="eyebrow">
						<span className="inline-block h-2 w-2 rounded-full bg-violet-400" />
						Configuration
					</div>
					<h1 className="page-title">Organize provider credentials and managed domains.</h1>
				</section>

				<section className="section-shell">
					<MultiProviderControl />
				</section>

				<section className="section-shell">
					<DomainControl />
				</section>
			</div>
		</main>
	);
};
