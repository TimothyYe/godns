interface ResolverProps {
	Resolver: string;
	onResolverChange?: (data: ResolverProps) => void;
}

export const Resolver = (props: ResolverProps) => {
	return (
		<label className="field-stack">
			<span className="field-label">DNS resolver</span>
			<input
				type="text"
				className="input input-bordered h-12 w-full rounded-2xl"
				placeholder="8.8.8.8"
				value={props.Resolver}
				onChange={(e) => {
					props.onResolverChange?.({
						Resolver: e.target.value
					});
				}}
			/>
			<span className="field-hint">Optional override for DNS resolution during provider lookups and verification.</span>
		</label>
	)
}
