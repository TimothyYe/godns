'use client';

import Link from "next/link";
import { useEffect, useContext, useState } from "react";
import { useRouter } from "next/navigation";
import { InfoIcon, DBIcon, TagIcon, ComputerIcon, SettingsIcon, GearIcon } from "@/components/icons";
import { CommonContext } from "./user";
import { Info, get_info, get_hours, get_date } from "@/api/info";
import { DomainCard } from "./domain-card";
import { PageShell, SectionCard } from "./page-shell";

export const Stat = () => {
	const router = useRouter();
	const [info, setInfo] = useState<Info | null>(null);
	const userStore = useContext(CommonContext);
	const { credentials, saveVersion, setCurrentPage } = userStore;

	useEffect(() => {
		if (!credentials) {
			router.push('/login');
			return;
		}

		setCurrentPage('Home');
		get_info(credentials).then((nextInfo) => {
			setInfo(nextInfo);
			saveVersion(nextInfo.version);
		});
	}, [credentials, router, saveVersion, setCurrentPage]);

	if (!info) {
		return (
			<div className="page-shell">
				<div className="panel flex min-h-[24rem] items-center justify-center">
					<span className="loading loading-spinner loading-lg" />
				</div>
			</div>
		);
	}

	const providerLabel = info.provider && (!info.providers || info.providers.length === 0)
		? info.provider
		: info.is_multi_provider
			? `${info.providers?.length || 0} providers`
			: 'Not configured';

	return (
		<PageShell
			eyebrow="Overview"
			title="GoDNS at a glance"
			description="Monitor the current runtime state, verify network configuration, and jump straight into the next configuration task."
			actions={(
				<>
					<Link className="btn btn-primary rounded-full px-5" href="/domains">Manage domains</Link>
					<Link className="btn btn-ghost rounded-full px-5" href="/logs">Review logs</Link>
				</>
			)}
		>
			<section className="metric-grid">
				<div className="metric-card">
					<div className="flex items-start justify-between">
						<p className="metric-label">Uptime</p>
						<InfoIcon />
					</div>
					<p className="metric-value">{get_hours(info.start_time)}</p>
					<p className="metric-meta">Since {get_date(info.start_time)}</p>
				</div>
				<div className="metric-card">
					<div className="flex items-start justify-between">
						<p className="metric-label">Providers</p>
						<GearIcon />
					</div>
					<p className="metric-value">{info.providers?.length || 0}</p>
					<p className="metric-meta">Configured provider profiles</p>
				</div>
				<div className="metric-card">
					<div className="flex items-start justify-between">
						<p className="metric-label">Domains</p>
						<DBIcon />
					</div>
					<p className="metric-value">{info.domain_num || 0}</p>
					<p className="metric-meta">Managed root domains</p>
				</div>
				<div className="metric-card">
					<div className="flex items-start justify-between">
						<p className="metric-label">Subdomains</p>
						<ComputerIcon />
					</div>
					<p className="metric-value">{info.sub_domain_num || 0}</p>
					<p className="metric-meta">Records tracked by GoDNS</p>
				</div>
			</section>

			<SectionCard
				title="Network posture"
				description="These values summarize how GoDNS currently resolves and publishes your public address."
			>
				<div className="grid gap-4 lg:grid-cols-3">
					<div className="panel-muted">
						<div className="flex items-start justify-between">
							<p className="metric-label">Public IP</p>
							<TagIcon />
						</div>
						<p className="mt-4 text-2xl font-semibold tracking-tight">{info.public_ip || 'Unavailable'}</p>
						<p className="metric-meta">Last detected public endpoint.</p>
					</div>
					<div className="panel-muted">
						<div className="flex items-start justify-between">
							<p className="metric-label">IP mode</p>
							<SettingsIcon />
						</div>
						<p className="mt-4 text-2xl font-semibold tracking-tight">{info.ip_mode || 'Unknown'}</p>
						<p className="metric-meta">Controls which upstream URL set GoDNS uses.</p>
					</div>
					<div className="panel-muted">
						<div className="flex items-start justify-between">
							<p className="metric-label">Provider strategy</p>
							<GearIcon />
						</div>
						<p className="mt-4 text-2xl font-semibold tracking-tight">{providerLabel}</p>
						<p className="metric-meta">
							{info.is_multi_provider ? 'Multiple provider profiles are active.' : 'Single provider mode or no provider configured yet.'}
						</p>
					</div>
				</div>
			</SectionCard>

			<SectionCard
				title="Configured domains"
				description="A quick view of the domains currently managed by this GoDNS instance."
				actions={<Link className="btn btn-ghost btn-sm rounded-full" href="/domains">Open domain editor</Link>}
			>
				{info.domains && info.domains.length > 0 ? (
					<div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
						{info.domains.map((domain, index) => (
							<DomainCard key={`${domain.domain_name}-${index}`} domain={domain} index={index} />
						))}
					</div>
				) : (
					<div className="panel-muted flex flex-col gap-3">
						<p className="text-lg font-semibold">No domains configured yet</p>
						<p className="text-sm text-base-content/65">
							Start by adding at least one provider profile, then attach a domain and the subdomains you want GoDNS to keep updated.
						</p>
						<div>
							<Link className="btn btn-primary btn-sm rounded-full px-5" href="/domains">Set up domains</Link>
						</div>
					</div>
				)}
			</SectionCard>
		</PageShell>
	);
};
