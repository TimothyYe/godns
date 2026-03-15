import classNames from "classnames";
import { SettingsIcon } from "@/components/icons";

interface IpModeProps {
	IPMode: string;
	IPUrls?: string[];
	IPV6Urls?: string[];
	onIpModeChange?: (data: IpModeProps) => void;
}

export const IpMode = (props: IpModeProps) => {
	const isIPV6 = props.IPMode === 'IPV6';
	const activeUrls = isIPV6 ? props.IPV6Urls || [] : props.IPUrls || [];

	return (
		<div className="grid gap-5">
			<div className="flex items-start justify-between gap-4">
				<div className="space-y-1">
					<div className="flex items-center gap-2">
						<SettingsIcon />
						<h3 className="text-lg font-semibold tracking-tight">IP detection mode</h3>
					</div>
					<p className="text-sm text-base-content/60">
						Choose whether GoDNS should detect and update IPv4 or IPv6 addresses.
					</p>
				</div>
				<div className="badge badge-primary badge-lg">{props.IPMode || 'Unset'}</div>
			</div>

			<div className="grid gap-3 sm:grid-cols-2">
				<button
					type="button"
					className={classNames(
						"rounded-[1.25rem] border p-4 text-left transition-colors",
						!isIPV6 ? "border-primary bg-primary/10" : "border-base-300 bg-base-100 hover:bg-base-200/60"
					)}
					onClick={() => props.onIpModeChange?.({
						IPMode: 'IPV4',
						IPUrls: props.IPUrls,
						IPV6Urls: props.IPV6Urls,
					})}
				>
					<p className="font-semibold">IPv4</p>
					<p className="mt-2 text-sm text-base-content/60">Use IPv4 detection URLs and update A records.</p>
				</button>
				<button
					type="button"
					className={classNames(
						"rounded-[1.25rem] border p-4 text-left transition-colors",
						isIPV6 ? "border-primary bg-primary/10" : "border-base-300 bg-base-100 hover:bg-base-200/60"
					)}
					onClick={() => props.onIpModeChange?.({
						IPMode: 'IPV6',
						IPUrls: props.IPUrls,
						IPV6Urls: props.IPV6Urls,
					})}
				>
					<p className="font-semibold">IPv6</p>
					<p className="mt-2 text-sm text-base-content/60">Use IPv6 detection URLs and update AAAA records.</p>
				</button>
			</div>

			<label className="field-stack">
				<span className="field-label">{isIPV6 ? 'IPv6 detection URLs' : 'IPv4 detection URLs'}</span>
				<textarea
					className="textarea textarea-bordered min-h-40 w-full rounded-[1.25rem]"
					placeholder={isIPV6 ? "https://api6.ipify.org" : "https://api.ipify.org"}
					value={activeUrls.join('\n')}
					onChange={(e) => {
						const nextUrls = e.target.value.split('\n').map((value) => value.trim()).filter(Boolean);
						props.onIpModeChange?.({
							IPMode: isIPV6 ? 'IPV6' : 'IPV4',
							IPUrls: isIPV6 ? props.IPUrls : nextUrls,
							IPV6Urls: isIPV6 ? nextUrls : props.IPV6Urls
						});
					}}
				/>
				<span className="field-hint">One URL per line. GoDNS will query these endpoints to determine the address it should publish.</span>
			</label>
		</div>
	);
}
