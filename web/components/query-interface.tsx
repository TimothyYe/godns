import { InterfaceIcon } from "./icons";

interface QueryInterfaceProps {
	QueryInterface: string;
	onQueryInterfaceChange?: (data: QueryInterfaceProps) => void;
}

export const QueryInterface = (props: QueryInterfaceProps) => {
	return (
		<div className="surface-panel-soft p-5 sm:p-6">
			<div className="flex items-start justify-between gap-4">
				<div className="w-full">
					<div className="metric-kicker">Interface</div>
					<h3 className="mt-2 text-xl font-semibold tracking-tight theme-heading">Query network interface</h3>
					<p className="mt-2 text-sm leading-7 theme-muted">Optionally send the public IP lookup out through a specific network interface instead of the default route. Useful on multi-WAN hosts where the default link sits behind CGNAT.</p>
					<fieldset className="theme-field mt-5">
						<label className="theme-field-label" htmlFor="query-interface">Network interface name</label>
						<input
							id="query-interface"
							type="text"
							className="input theme-input w-full rounded-2xl"
							placeholder="Leave empty to use the default route: e.g. wan0"
							value={props.QueryInterface}
							onChange={(e) => {
								if (props.onQueryInterfaceChange) {
									props.onQueryInterfaceChange({
										QueryInterface: e.target.value
									});
								}
							}}
						/>
						<span className="theme-field-hint">Binds the lookup socket to the device on Linux and macOS. Other platforms only bind the source address.</span>
					</fieldset>
				</div>
				<div className="metric-icon text-cyan-300">
					<InterfaceIcon />
				</div>
			</div>
		</div>
	);
};
