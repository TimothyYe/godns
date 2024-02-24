import { GearIcon } from "./icons"
export const WebHook = () => {
	return (
		<div className="stats shadow bg-primary-content stats-vertical lg:stats-horizontal">
			<div className="stat gap-2">
				<div className="stat-title">Webhook</div>
				<div className="flex flex-col gap-3">
					<div className="flex flex-row items-center justify-start gap-2">
						<span className="label-text text-slate-500 ">Enable Webhook</span>
						<input type="checkbox" className="toggle toggle-primary" checked={false} />
						<div className="flex flex-grow justify-end text-secondary">
							<GearIcon />
						</div>
					</div>
					<textarea className="textarea textarea-primary w-full h-28" placeholder="Input request body"></textarea>
				</div>
			</div>
		</div>
	)
}