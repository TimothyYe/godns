import { DBIcon } from "./icons";

interface ResolverProps {
	Resolver: string;
	onResolverChange?: (data: ResolverProps) => void;
}

export const Resolver = (props: ResolverProps) => {
	return (
		<div className="surface-panel-soft p-5 sm:p-6">
			<div className="flex items-start justify-between gap-4">
				<div className="w-full">
					<div className="metric-kicker">Resolver</div>
					<h3 className="mt-2 text-xl font-semibold tracking-tight theme-heading">DNS lookup resolver</h3>
					<p className="mt-2 text-sm leading-7 theme-muted">Choose the resolver GoDNS queries when it checks existing records before pushing updates.</p>
					<fieldset className="theme-field mt-5">
						<label className="theme-field-label" htmlFor="resolver">Resolver address</label>
						<input
							id="resolver"
							type="text"
							className="input theme-input w-full rounded-2xl"
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
						<span className="theme-field-hint">Use a public resolver like <code>1.1.1.1</code> or <code>8.8.8.8</code>, or point at your own internal DNS endpoint.</span>
					</fieldset>
				</div>
				<div className="metric-icon text-violet-300">
					<DBIcon />
				</div>
			</div>
		</div>
	);
};
