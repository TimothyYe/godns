import { InfoIcon, DBIcon, TagIcon, ComputerIcon } from "@/components/icons";

export const Stat = () => {
	return (
		<div className="stats shadow bg-primary-content">
			<div className="stat">
				<div className="stat-figure text-secondary">
					<InfoIcon />
				</div>
				<div className="stat-title">Uptime</div>
				<div className="stat-value text-primary">31K</div>
				<div className="stat-desc">Jan 1st - Feb 1st</div>
			</div>

			<div className="stat">
				<div className="stat-figure text-secondary">
					<DBIcon />
				</div>
				<div className="stat-title">Domains</div>
				<div className="stat-value text-info">4,200</div>
				<div className="stat-desc">Subdomains configured</div>
			</div>

			<div className="stat">
				<div className="stat-figure text-secondary">
					<ComputerIcon />
				</div>
				<div className="stat-title">Subdomains</div>
				<div className="stat-value text-error">4,200</div>
				<div className="stat-desc">Subdomains configured</div>
			</div>

			<div className="stat">
				<div className="stat-figure text-secondary">
					<TagIcon />
				</div>
				<div className="stat-title">Status</div>
				<div className="stat-value text-success justify-center items-center">Running</div>
				<div className="stat-desc">Server status</div>
			</div>
		</div>
	);
}
