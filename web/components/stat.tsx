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
		get_info(credentials).then((info) => {
			setInfo(info);
			saveVersion(info.version);
		});
	}, [saveVersion, credentials, setCurrentPage, router]);

	return (
		info ? (
			<div className="flex flex-col max-w-screen-lg">
				<span className="text-xl font-semibold text-neutral-500 ml-1 mb-1">Basic Info</span>
				<div className="stats shadow bg-primary-content stats-vertical lg:stats-horizontal">
					<div className="stat">
						<div className="stat-figure text-secondary">
							<InfoIcon />
						</div>
						<div className="stat-title">Uptime</div>
						<div className="stat-value text-primary">{info ? get_hours(info.start_time) : null}</div>
						<div className="stat-desc">Since {info ? get_date(info.start_time) : null}</div>
					</div>

					<div className="stat">
						<div className="stat-figure text-secondary">
							<DBIcon />
						</div>
						<div className="stat-title">Domains</div>
						<div className="stat-value text-info">{info ? info.domain_num : 0}</div>
						<div className="stat-desc">Domains configured</div>
					</div>

					<div className="stat">
						<div className="stat-figure text-secondary">
							<ComputerIcon />
						</div>
						<div className="stat-title">Subdomains</div>
						<div className="stat-value text-error">{info ? info.sub_domain_num : 0}</div>
						<div className="stat-desc">Subdomains configured</div>
					</div>
				</div>
				<span className="text-xl font-semibold text-neutral-500 ml-1 mb-1 mt-5">Network Info</span>
				<div className="stats shadow bg-primary-content stats-vertical lg:stats-horizontal">
					<div className="stat">
						<div className="stat-figure text-secondary">
							<TagIcon />
						</div>
						<div className="stat-title">Public IP</div>
						<div className="stat-value text-primary">{info ? info.public_ip : 'N/A'}</div>
						<div className="stat-desc">The public IP address</div>
					</div>

					<div className="stat">
						<div className="stat-figure text-secondary">
							<SettingsIcon />
						</div>
						<div className="stat-title">IP Mode</div>
						<div className="stat-value text-info">{info ? info.ip_mode : 'N/A'}</div>
						<div className="stat-desc">The IP mode</div>
					</div>

					<div className="stat">
						<div className="stat-figure text-secondary">
							<GearIcon />
						</div>
						<div className="stat-title">Provider</div>
						<div className="stat-value text-error">{info ? info.provider : 'N/A'}</div>
						<div className="stat-desc">Provider configured</div>
					</div>
				</div>
				<span className="text-xl font-semibold text-neutral-500 ml-1 mb-1 mt-5">Domain Info</span>
				<div className="flex flex-wrap gap-2">
					{
						info && info.domains ? info.domains.map((domain, index) => {
							return (
								<DomainCard key={index} domain={domain} index={index} />
							);
						}) : null
					}
				</div>
			</div>
		) : null);
}
