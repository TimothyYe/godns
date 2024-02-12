'use client';
import { InfoIcon, DBIcon, TagIcon, ComputerIcon, SettingsIcon, GearIcon } from "@/components/icons";
import { useEffect, useContext, useState } from "react";
import { CommonContext } from "./user";
import { Info, get_info, get_hours, get_date } from "@/api/info";

export const Stat = () => {
	const userStore = useContext(CommonContext);
	const { credentials, setVersion } = userStore;
	const [info, setInfo] = useState<Info>({
		version: '',
		start_time: 0,
		domains: 0,
		sub_domains: 0
	});

	useEffect(() => {
		get_info(credentials).then((info) => {
			setInfo(info);
			setVersion(info.version);
		});
	}, [credentials, setVersion]);

	return (
		<div className="flex flex-col max-w-screen-md">
			<span className="text-xl font-semibold text-neutral-500 ml-1 mb-1">Basic Info</span>
			<div className="stats shadow bg-primary-content stats-vertical lg:stats-horizontal">
				<div className="stat">
					<div className="stat-figure text-secondary">
						<InfoIcon />
					</div>
					<div className="stat-title">Uptime</div>
					<div className="stat-value text-primary">{get_hours(info.start_time)}</div>
					<div className="stat-desc">Since {get_date(info.start_time)}</div>
				</div>

				<div className="stat">
					<div className="stat-figure text-secondary">
						<DBIcon />
					</div>
					<div className="stat-title">Domains</div>
					<div className="stat-value text-info">{info.domains}</div>
					<div className="stat-desc">Domains configured</div>
				</div>

				<div className="stat">
					<div className="stat-figure text-secondary">
						<ComputerIcon />
					</div>
					<div className="stat-title">Subdomains</div>
					<div className="stat-value text-error">{info.sub_domains}</div>
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
					<div className="stat-value text-primary">192.168.1.1</div>
					<div className="stat-desc">The public IP address</div>
				</div>

				<div className="stat">
					<div className="stat-figure text-secondary">
						<SettingsIcon />
					</div>
					<div className="stat-title">IP Mode</div>
					<div className="stat-value text-info">IPv4</div>
					<div className="stat-desc"></div>
				</div>

				<div className="stat">
					<div className="stat-figure text-secondary">
						<GearIcon />
					</div>
					<div className="stat-title">Provider</div>
					<div className="stat-value text-error">Cloudflare</div>
					<div className="stat-desc">Provider configured</div>
				</div>
			</div>
			<span className="text-xl font-semibold text-neutral-500 ml-1 mb-1 mt-5">Domain Info</span>
			<div className="flex flex-wrap gap-2">
				<div className="card w-full md:w-1/3 bg-primary-content shadow-xl">
					<figure></figure>
					<div className="card-body">
						<h2 className="card-title">
							vpsdalao.com
						</h2>
						<div className="flex flex-wrap card-actions justify-start">
							<div className="badge badge-primary">ipv4</div>
							<div className="badge badge-primary">ipv67890</div>
							<div className="badge badge-primary">ipv612345</div>
							<div className="badge badge-primary">ipv63322</div>
						</div>
					</div>
				</div>
				<div className="card w-full md:flex-1 bg-primary-content shadow-xl">
					<figure></figure>
					<div className="card-body">
						<h2 className="card-title">
							zhujijun.com
						</h2>
						<div className="card-actions justify-start">
							<div className="badge badge-primary">ipv4</div>
							<div className="badge badge-primary">ipv6</div>
						</div>
					</div>
				</div>
				<div className="card w-full md:flex-1 bg-primary-content shadow-xl">
					<figure></figure>
					<div className="card-body">
						<h2 className="card-title">
							zhujijun2.com
						</h2>
						<div className="card-actions justify-start">
							<div className="badge badge-primary">ipv4</div>
							<div className="badge badge-primary">ipv6</div>
						</div>
					</div>
				</div>
				<div className="card w-full md:w-1/3 bg-primary-content shadow-xl">
					<figure></figure>
					<div className="card-body">
						<h2 className="card-title">
							zhujijun21.com
						</h2>
						<div className="card-actions justify-start">
							<div className="badge badge-primary">ipv4</div>
							<div className="badge badge-primary">ipv6</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	);
}
