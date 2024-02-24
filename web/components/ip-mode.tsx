import { SettingsIcon } from "@/components/icons";

export const IpMode = () => {
	return (
		<div className="stats shadow bg-primary-content stats-vertical lg:stats-horizontal">
			<div className="stat">
				<div className="stat-figure text-secondary">
					<SettingsIcon />
				</div>
				<div className="stat-title">IP Mode</div>
				<div className="stat-value text-primary">IPV4</div>
				<div className="stat-desc">The current IP mode</div>
				<div className="flex flex-row items-center gap-2 stat-actions">
					<span className="label-text text-slate-500 ">Enable IPv6</span>
					<input type="checkbox" className="toggle toggle-primary" checked={false} />
				</div>
			</div>

			<div className="stat gap-2">
				<div className="stat-title">URLs</div>
				<div className="flex flex-col gap-2">
					<textarea className="textarea textarea-primary w-96 h-28 " placeholder="Input IP urls"></textarea>
					<div className="flex justify-end">
						<button className="flex justify-end btn btn-sm btn-primary">Save</button>
					</div>
				</div>
			</div>
		</div>
	)
}