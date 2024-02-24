import { ComputerIcon } from "./icons"
export const Proxy = () => {
	return (
		<div className="stats shadow bg-primary-content stats-vertical lg:stats-horizontal">
			<div className="stat gap-2">
				<div className="stat-title">Proxy Settings</div>
				<div className="flex flex-col gap-3">
					<div className="flex flex-row items-center justify-start gap-2">
						<span className="label-text text-slate-500 ">Enable IPv6</span>
						<input type="checkbox" className="toggle toggle-primary" checked={false} />
						<span className="label-text text-slate-500 ">Skip SSL Verify</span>
						<input type="checkbox" className="toggle toggle-primary" checked={false} />
						<div className="flex w-1/2 justify-end text-secondary">
							<ComputerIcon />
						</div>
					</div>

					<input type="text" className="input input-primary w-full input-disabled" placeholder="Input Proxy: e.g. 127.0.0.1:8080"></input>
					<div className="flex justify-end">
						<button className="flex justify-end btn btn-sm btn-primary">Save</button>
					</div>
				</div>
			</div>
		</div>
	)
}