import classNames from "classnames";
import { ComputerIcon } from "./icons";

interface ProxyProps {
	EnableProxy: boolean;
	Socks5Proxy: string;
	SkipSSLVerify: boolean;
	onProxyChange?: (data: ProxyProps) => void;
}

export const Proxy = (props: ProxyProps) => {
	return (
		<div className="surface-panel-soft p-5 sm:p-6">
			<div className="flex flex-col gap-5">
				<div className="flex items-start justify-between gap-4">
					<div>
						<div className="metric-kicker">Proxy</div>
						<h3 className="mt-2 text-xl font-semibold tracking-tight theme-heading">Forward requests through SOCKS5</h3>
						<p className="mt-2 text-sm leading-7 theme-muted">Enable a proxy for outbound requests and control SSL verification for constrained environments.</p>
					</div>
					<div className="metric-icon text-sky-300">
						<ComputerIcon />
					</div>
				</div>

				<div className="grid gap-4 md:grid-cols-2">
					<label className="inline-flex items-center gap-3 text-sm theme-muted">
						<input
							type="checkbox"
							className="toggle toggle-primary"
							checked={props.EnableProxy}
							onChange={() => {
								if (props.onProxyChange) {
									props.onProxyChange({
										EnableProxy: !props.EnableProxy,
										Socks5Proxy: props.Socks5Proxy,
										SkipSSLVerify: props.SkipSSLVerify,
									});
								}
							}}
						/>
						<span>Enable proxy</span>
					</label>
					<label className="inline-flex items-center gap-3 text-sm theme-muted">
						<input
							type="checkbox"
							className="toggle toggle-primary"
							checked={props.SkipSSLVerify}
							onChange={() => {
								if (props.onProxyChange) {
									props.onProxyChange({
										EnableProxy: props.EnableProxy,
										Socks5Proxy: props.Socks5Proxy,
										SkipSSLVerify: !props.SkipSSLVerify,
									});
								}
							}}
						/>
						<span>Skip SSL verification</span>
					</label>
				</div>

				<fieldset className="theme-field">
					<label className="theme-field-label" htmlFor="socks5-proxy">SOCKS5 proxy address</label>
					<input
						id="socks5-proxy"
						type="text"
						className={classNames("input theme-input w-full rounded-2xl", {
							'border-rose-400/60': props.EnableProxy && !props.Socks5Proxy
						})}
						disabled={!props.EnableProxy}
						placeholder="Input proxy: e.g. 127.0.0.1:8080"
						value={props.Socks5Proxy}
						onChange={(e) => {
							if (props.onProxyChange) {
								props.onProxyChange({
									EnableProxy: props.EnableProxy,
									Socks5Proxy: e.target.value,
									SkipSSLVerify: props.SkipSSLVerify,
								});
							}
						}}
					/>
					<span className="theme-field-hint">Leave empty to bypass the proxy. When enabled, GoDNS routes HTTP calls through this endpoint.</span>
				</fieldset>
			</div>
		</div>
	);
};
