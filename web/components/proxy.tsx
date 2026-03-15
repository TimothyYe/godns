interface ProxyProps {
	EnableProxy: boolean;
	Socks5Proxy: string;
	SkipSSLVerify: boolean;
	onProxyChange?: (data: ProxyProps) => void;
}

export const Proxy = (props: ProxyProps) => {
	return (
		<div className="grid gap-5">
			<div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
				<div className="space-y-1">
					<h3 className="text-lg font-semibold tracking-tight">Proxy access</h3>
					<p className="text-sm text-base-content/60">
						Route outbound requests through a SOCKS5 proxy if the GoDNS runtime cannot access provider APIs directly.
					</p>
				</div>
				<label className="flex items-center gap-3 rounded-full border border-base-300 px-4 py-2">
					<span className="text-sm font-medium">Enable proxy</span>
					<input
						type="checkbox"
						className="toggle toggle-primary"
						checked={props.EnableProxy}
						onChange={(e) => {
							props.onProxyChange?.({
								EnableProxy: e.target.checked,
								Socks5Proxy: props.Socks5Proxy,
								SkipSSLVerify: props.SkipSSLVerify,
							});
						}}
					/>
				</label>
			</div>

			<label className="field-stack">
				<span className="field-label">SOCKS5 proxy endpoint</span>
				<input
					type="text"
					className="input input-bordered h-12 w-full rounded-2xl"
					disabled={!props.EnableProxy}
					placeholder="127.0.0.1:8080"
					value={props.Socks5Proxy}
					onChange={(e) => {
						props.onProxyChange?.({
							EnableProxy: props.EnableProxy,
							Socks5Proxy: e.target.value,
							SkipSSLVerify: props.SkipSSLVerify,
						});
					}}
				/>
				<span className="field-hint">Only required when proxying outbound traffic.</span>
			</label>

			<label className="flex items-center justify-between rounded-[1.25rem] border border-base-300 bg-base-100 px-4 py-3">
				<div className="space-y-1">
					<p className="text-sm font-semibold">Skip SSL verification</p>
					<p className="text-xs text-base-content/60">Use only when your environment requires insecure TLS connections.</p>
				</div>
				<input
					type="checkbox"
					className="toggle toggle-warning"
					checked={props.SkipSSLVerify}
					onChange={(e) => {
						props.onProxyChange?.({
							EnableProxy: props.EnableProxy,
							Socks5Proxy: props.Socks5Proxy,
							SkipSSLVerify: e.target.checked,
						});
					}}
				/>
			</label>
		</div>
	)
}
