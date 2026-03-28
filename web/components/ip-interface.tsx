import { InterfaceIcon } from "./icons";

interface IPInterfaceProps {
	IPInterface: string;
	onIPInterfaceChange?: (data: IPInterfaceProps) => void;
}

export const IPInterface = (props: IPInterfaceProps) => {
	return (
		<div className="surface-panel-soft p-5 sm:p-6">
			<div className="flex items-start justify-between gap-4">
				<div className="w-full">
					<div className="metric-kicker">Interface</div>
					<h3 className="mt-2 text-xl font-semibold tracking-tight theme-heading">Fallback network interface</h3>
					<p className="mt-2 text-sm leading-7 theme-muted">Optionally read the public address directly from a local network interface when external lookup URLs are unavailable.</p>
					<fieldset className="theme-field mt-5">
						<label className="theme-field-label" htmlFor="ip-interface">Network interface name</label>
						<input
							id="ip-interface"
							type="text"
							className="input theme-input w-full rounded-2xl"
							placeholder="Input the network interface name: e.g. eth0"
							value={props.IPInterface}
							onChange={(e) => {
								if (props.onIPInterfaceChange) {
									props.onIPInterfaceChange({
										IPInterface: e.target.value
									});
								}
							}}
						/>
						<span className="theme-field-hint">Common values include <code>eth0</code>, <code>en0</code>, or your platform-specific interface alias.</span>
					</fieldset>
				</div>
				<div className="metric-icon text-cyan-300">
					<InterfaceIcon />
				</div>
			</div>
		</div>
	);
};
