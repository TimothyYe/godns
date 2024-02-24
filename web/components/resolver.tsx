import { DBIcon } from "./icons";

interface ResolverProps {
	Resolver: string;
	onResolverChange?: (data: ResolverProps) => void;
}

export const Resolver = (props: ResolverProps) => {
	return (
		<div className="stats shadow bg-primary-content stats-vertical lg:stats-horizontal">
			<div className="stat gap-2">
				<div className="stat-title">Resolver</div>
				<div className="flex flex-col gap-3">
					<div className="flex flex-row items-center justify-start gap-2">
						<span className="label-text text-slate-500 ">Set the DNS resolver</span>
						<div className="flex flex-grow justify-end text-secondary">
							<DBIcon />
						</div>
					</div>
					<input
						type="text"
						className="input input-primary w-full"
						placeholder="Input DNS resolver: e.g. 8.8.8.8"
						value={props.Resolver}
						onChange={(e) => {
							if (props.onResolverChange) {
								props.onResolverChange({
									Resolver: e.target.value
								});
							}
						}}
					/>
				</div>
			</div>
		</div>
	)
}