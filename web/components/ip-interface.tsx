import { InterfaceIcon } from "./icons";

interface IPInterfaceProps {
	IPInterface: string;
	onIPInterfaceChange?: (data: IPInterfaceProps) => void;
}

export const IPInterface = (props: IPInterfaceProps) => {
	return (
		<div className="stats shadow bg-primary-content stats-vertical lg:stats-horizontal">
			<div className="stat gap-2">
				<div className="stat-title">IP Interface</div>
				<div className="flex flex-col gap-3">
					<div className="flex flex-row items-center justify-start gap-2">
						<span className="label-text text-slate-500 ">Set the network interface</span>
						<div className="flex flex-grow justify-end text-secondary">
							<InterfaceIcon />
						</div>
					</div>
					<input
						type="text"
						className="input input-primary w-full"
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
				</div>
			</div>
		</div>
	)
}