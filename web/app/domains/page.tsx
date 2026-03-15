'use client';

import Link from 'next/link';
import React, { useContext, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { CommonContext } from '@/components/user';
import { DomainControl } from '@/components/domain-control';
import { get_info } from '@/api/info';
import { MultiProviderControl } from '@/components/multi-provider-control';
import { PageShell } from '@/components/page-shell';

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
	}, [credentials, router, saveVersion, setCurrentPage]);

	return (
		<PageShell
			eyebrow="Configuration"
			title="Providers and domains"
			description="Configure provider credentials first, then attach domains and subdomains to the right provider profile. This keeps multi-provider setups explicit and easier to audit."
			actions={(
				<>
					<Link className="btn btn-primary rounded-full px-5" href="/network">Review network settings</Link>
					<Link className="btn btn-ghost rounded-full px-5" href="/logs">Open logs</Link>
				</>
			)}
		>
			<div className="space-y-6">
				<MultiProviderControl />
				<DomainControl />
			</div>
		</PageShell>
	);
};
