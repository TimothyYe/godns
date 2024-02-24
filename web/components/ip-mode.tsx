import { SettingsIcon } from "@/components/icons";
import { useState } from "react";

interface IpModeProps {
	IPMode: string;
	IPUrls?: string[];
	IPV6Urls?: string[];
}

export const IpMode = (props: IpModeProps) => {
	const [isIPV6, setIsIPV6] = useState(props.IPMode === 'IPV6' ? true : false);

	return (
		<div className="stats shadow bg-primary-content stats-vertical lg:stats-horizontal">
			<div className="stat">
				<div className="stat-figure text-secondary">
					<SettingsIcon />
				</div>
				<div className="stat-title">IP Mode</div>
				<div className="stat-value text-primary">{props.IPMode}</div>
				<div className="stat-desc">The current IP mode</div>
				<div className="flex flex-row items-center gap-2 stat-actions">
					<span className="label-text text-slate-500 ">Enable IPv6</span>
					<input type="checkbox" className="toggle toggle-primary" checked={isIPV6} />
				</div>
			</div>

			<div className="stat gap-2">
				<div className="stat-title">URLs</div>
				<div className="flex flex-col gap-2">
					<textarea
						className="textarea textarea-primary w-96 h-28 "
						placeholder="Input urls for fetching the online IP"
						value={isIPV6 ? props.IPV6Urls?.join('\n') : props.IPUrls?.join('\n')}
					/>
				</div>
			</div>
		</div>
	)
}