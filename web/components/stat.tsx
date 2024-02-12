'use client';
import { InfoIcon, DBIcon, TagIcon, ComputerIcon } from "@/components/icons";
import { useEffect, useContext, useState } from "react";
import { UserContext } from "./user";
import { Info, get_info, get_hours, get_date } from "@/api/info";

export const Stat = () => {
	const userStore = useContext(UserContext);
	const { credentials } = userStore;
	const [info, setInfo] = useState<Info>({
		version: '',
		start_time: 0,
		domains: 0,
		sub_domains: 0
	});

	useEffect(() => {
		get_info(credentials).then((info) => {
			setInfo(info);
		});
	}, [credentials]);

	return (
		<div className="stats shadow bg-primary-content">
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

			<div className="stat">
				<div className="stat-figure text-secondary">
					<TagIcon />
				</div>
				<div className="stat-title">Version</div>
				<div className="stat-value text-success justify-center items-center">{info.version}</div>
				<div className="stat-desc">Current version</div>
			</div>
		</div>
	);
}
