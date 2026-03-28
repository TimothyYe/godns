import { SettingsIcon } from "@/components/icons";

interface IpModeProps {
	IPMode: string;
	IPUrls?: string[];
	IPV6Urls?: string[];
	onIpModeChange?: (data: IpModeProps) => void;
}

export const IpMode = (props: IpModeProps) => {
	const isIPV6 = props.IPMode === 'IPV6';

	return (
		<div className="surface-panel-soft p-5 sm:p-6">
			<div className="grid gap-5 lg:grid-cols-[minmax(0,0.9fr)_minmax(0,1.6fr)]">
				<div className="flex flex-col gap-4">
					<div className="flex items-start justify-between gap-4">
						<div>
							<div className="metric-kicker">IP Mode</div>
							<h3 className="mt-2 text-xl font-semibold tracking-tight theme-heading">Choose the active address family</h3>
							<p className="mt-2 text-sm leading-7 theme-muted">Switch between IPv4 and IPv6 lookup sources without leaving the dashboard.</p>
						</div>
						<div className="metric-icon text-sky-300">
							<SettingsIcon />
						</div>
					</div>
					<div className="theme-chip-sky w-fit rounded-full px-3 py-1 text-xs font-medium">
						Current mode: {props.IPMode || 'IPV4'}
					</div>
					<label className="inline-flex items-center gap-3 text-sm theme-muted">
						<input
							type="checkbox"
							className="toggle toggle-primary"
							checked={isIPV6}
							onChange={() => {
								if (props.onIpModeChange) {
									props.onIpModeChange({
										IPMode: isIPV6 ? 'IPV4' : 'IPV6',
										IPUrls: props.IPUrls,
										IPV6Urls: props.IPV6Urls,
									});
								}
							}}
						/>
						<span>Enable IPv6 mode</span>
					</label>
				</div>

				<fieldset className="theme-field">
					<label className="theme-field-label" htmlFor="ip-lookup-urls">Lookup URLs</label>
					<textarea
						id="ip-lookup-urls"
						className="textarea theme-input h-32 w-full rounded-2xl"
						placeholder="Input URLs for fetching the online IP"
						value={isIPV6 && props.IPV6Urls ? props.IPV6Urls.join('\n') : !isIPV6 && props.IPUrls ? props.IPUrls.join('\n') : ''}
						onChange={(e) => {
							if (props.onIpModeChange) {
								props.onIpModeChange({
									IPMode: isIPV6 ? 'IPV6' : 'IPV4',
									IPUrls: isIPV6 ? props.IPUrls : e.target.value.split('\n'),
									IPV6Urls: isIPV6 ? e.target.value.split('\n') : props.IPV6Urls
								});
							}
						}}
					/>
					<span className="theme-field-hint">Provide one URL per line. GoDNS rotates through these endpoints until it finds a valid public address.</span>
				</fieldset>
			</div>
		</div>
	);
};
