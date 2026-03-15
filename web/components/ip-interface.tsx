interface IPInterfaceProps {
	IPInterface: string;
	onIPInterfaceChange?: (data: IPInterfaceProps) => void;
}

export const IPInterface = (props: IPInterfaceProps) => {
	return (
		<label className="field-stack">
			<span className="field-label">Network interface</span>
			<input
				type="text"
				className="input input-bordered h-12 w-full rounded-2xl"
				placeholder="eth0"
				value={props.IPInterface}
				onChange={(e) => {
					props.onIPInterfaceChange?.({
						IPInterface: e.target.value
					});
				}}
			/>
			<span className="field-hint">Set this only when GoDNS should bind to a specific interface to discover the local address.</span>
		</label>
	)
}
