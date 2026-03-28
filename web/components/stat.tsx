'use client';
import { InfoIcon, DBIcon, TagIcon, ComputerIcon, SettingsIcon, GearIcon } from "@/components/icons";
import { useEffect, useContext, useState } from "react";
import { useRouter } from "next/navigation";
import { CommonContext } from "./user";
import { Info, get_info, get_hours, get_date } from "@/api/info";
import { DomainCard } from "./domain-card";

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
	}, [saveVersion, credentials, setCurrentPage, router]);

	return info ? (
		<div className="flex flex-col gap-8">
			<section className="page-hero page-hero-compact">
				<div className="eyebrow">
					<span className="inline-block h-2 w-2 rounded-full bg-sky-400" />
					System Overview
				</div>
				<h1 className="page-title">Manage your DNS updates with less friction.</h1>
				<p className="page-description">
					Check runtime health, network mode, and the domains currently under management.
				</p>
			</section>

			<section className="section-shell">
				<div className="section-header">
					<div className="section-copy">
						<div className="section-label ml-0">Basic Info</div>
						<h2 className="text-2xl font-semibold tracking-tight theme-heading">Runtime snapshot</h2>
						<p className="mt-2 text-sm leading-7 theme-muted">
							A compact view of how long this node has been up and how much DNS configuration it is carrying.
						</p>
					</div>
				</div>
				<div className="metric-grid">
					<div className="metric-card">
						<div className="metric-head">
							<div>
								<div className="metric-kicker">Uptime</div>
								<div className="metric-value">{get_hours(info.start_time)}</div>
							</div>
							<div className="metric-icon text-sky-300">
								<InfoIcon />
							</div>
						</div>
						<div className="metric-note">Since {get_date(info.start_time)}</div>
					</div>
					<div className="metric-card">
						<div className="metric-head">
							<div>
								<div className="metric-kicker">Providers</div>
								<div className="metric-value">{info.providers ? info.providers.length : 0}</div>
							</div>
							<div className="metric-icon text-violet-300">
								<GearIcon />
							</div>
						</div>
						<div className="metric-note">Providers configured</div>
					</div>
					<div className="metric-card">
						<div className="metric-head">
							<div>
								<div className="metric-kicker">Domains</div>
								<div className="metric-value">{info.domain_num}</div>
							</div>
							<div className="metric-icon text-cyan-300">
								<DBIcon />
							</div>
						</div>
						<div className="metric-note">Root domains configured</div>
					</div>
					<div className="metric-card">
						<div className="metric-head">
							<div>
								<div className="metric-kicker">Subdomains</div>
								<div className="metric-value">{info.sub_domain_num}</div>
							</div>
							<div className="metric-icon text-amber-300">
								<ComputerIcon />
							</div>
						</div>
						<div className="metric-note">Subdomains configured</div>
					</div>
				</div>
			</section>

			<section className="section-shell">
				<div className="section-header">
					<div className="section-copy">
						<div className="section-label ml-0">Network Info</div>
						<h2 className="text-2xl font-semibold tracking-tight theme-heading">Current network state</h2>
						<p className="mt-2 text-sm leading-7 theme-muted">
							The current public address, update mode, and provider strategy the daemon is operating with.
						</p>
					</div>
				</div>
				<div className="metric-grid-compact">
					<div className="metric-card">
						<div className="metric-head">
							<div>
								<div className="metric-kicker">Public IP</div>
								<div className="metric-value text-2xl">{info.public_ip || 'N/A'}</div>
							</div>
							<div className="metric-icon text-sky-300">
								<TagIcon />
							</div>
						</div>
						<div className="metric-note">Detected public address</div>
					</div>
					<div className="metric-card">
						<div className="metric-head">
							<div>
								<div className="metric-kicker">IP Mode</div>
								<div className="metric-value">{info.ip_mode || 'N/A'}</div>
							</div>
							<div className="metric-icon text-emerald-300">
								<SettingsIcon />
							</div>
						</div>
						<div className="metric-note">Current update mode</div>
					</div>
					<div className="metric-card">
						<div className="metric-head">
							<div>
								<div className="metric-kicker">Provider Strategy</div>
								<div className="metric-value text-2xl">
									{info.provider && (!info.providers || info.providers.length === 0) ? info.provider : 'Multiple'}
								</div>
							</div>
							<div className="metric-icon text-rose-300">
								<GearIcon />
							</div>
						</div>
						<div className="metric-note">
							{info.provider && (!info.providers || info.providers.length === 0) ? 'Single provider configured' : 'Multiple providers configured'}
						</div>
					</div>
				</div>
			</section>

			<section className="section-shell">
				<div className="section-header">
					<div className="section-copy">
						<div className="section-label ml-0">Domain Info</div>
						<h2 className="text-2xl font-semibold tracking-tight theme-heading">Tracked domains</h2>
						<p className="mt-2 text-sm leading-7 theme-muted">
							These are the root domains and subdomains currently being managed by the running GoDNS instance.
						</p>
					</div>
					{info.domains?.length ? (
						<div className="theme-chip rounded-full px-3 py-1 text-xs font-medium">
							{info.domains.length} tracked
						</div>
					) : null}
				</div>
				<div className="mt-5 grid gap-4 lg:grid-cols-3">
					{info.domains && info.domains.length > 0 ? info.domains.map((domain, index) => (
						<DomainCard key={index} domain={domain} index={index} />
					)) : (
						<div className="surface-panel-soft col-span-full px-6 py-10 text-center">
							<p className="text-lg font-medium theme-heading">No domains configured yet</p>
							<p className="mt-2 text-sm theme-muted">Add a provider and domain from the Domains page to start managing DNS updates.</p>
						</div>
					)}
				</div>
			</section>
		</div>
	) : (
		<div className="surface-panel flex min-h-[18rem] items-center justify-center">
			<span className="loading loading-spinner loading-lg text-sky-400" />
		</div>
	);
}
