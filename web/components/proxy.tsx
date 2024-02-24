import { ComputerIcon } from "./icons"
import classNames from "classnames";

interface ProxyProps {
	EnableProxy: boolean;
	Socks5Proxy: string;
	SkipSSLVerify: boolean;
	onProxyChange?: (data: ProxyProps) => void;
}

export const Proxy = (props: ProxyProps) => {
	let proxyEnabled = props.EnableProxy;
	let skipSSLVerify = props.SkipSSLVerify;

	return (
		<div className="stats shadow bg-primary-content stats-vertical lg:stats-horizontal">
			<div className="stat gap-2">
				<div className="stat-title">Proxy Settings</div>
				<div className="flex flex-col gap-3">
					<div className="flex flex-row items-center justify-start gap-2">
						<span className="label-text text-slate-500 ">Enable Proxy</span>
						<input
							type="checkbox"
							className="toggle toggle-primary"
							checked={proxyEnabled}
							onClick={() => {
								proxyEnabled = !proxyEnabled;
								if (props.onProxyChange) {
									props.onProxyChange({
										EnableProxy: proxyEnabled,
										Socks5Proxy: props.Socks5Proxy,
										SkipSSLVerify: props.SkipSSLVerify,
									});
								}
							}}
						/>
						<span className="label-text text-slate-500 ">Skip SSL Verify</span>
						<input
							type="checkbox"
							className="toggle toggle-primary"
							checked={skipSSLVerify}
							onClick={() => {
								skipSSLVerify = !skipSSLVerify;
								if (props.onProxyChange) {
									props.onProxyChange({
										EnableProxy: proxyEnabled,
										Socks5Proxy: props.Socks5Proxy,
										SkipSSLVerify: skipSSLVerify,
									});
								}
							}}
						/>
						<div className="flex flex-grow justify-end text-secondary">
							<ComputerIcon />
						</div>
					</div>

					<input
						type="text"
						className={classNames("input input-primary w-full", {
							'input-error': proxyEnabled && !props.Socks5Proxy
						})}
						disabled={!props.EnableProxy}
						placeholder="Input Proxy: e.g. 127.0.0.1:8080"
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
				</div>
			</div>
		</div>
	)
}