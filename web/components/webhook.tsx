import { GearIcon } from "./icons";

interface WebHookProps {
	Enabled: boolean;
	Url: string;
	RequestBody: string;
	onWebHookChange?: (data: WebHookProps) => void;
}

export const WebHook = (props: WebHookProps) => {
	return (
		<div className="surface-panel-soft p-5 sm:p-6">
			<div className="flex flex-col gap-5">
				<div className="flex items-start justify-between gap-4">
					<div>
						<div className="metric-kicker">Webhook</div>
						<h3 className="mt-2 text-xl font-semibold tracking-tight theme-heading">Push change notifications downstream</h3>
						<p className="mt-2 text-sm leading-7 theme-muted">Trigger a GET or POST request whenever GoDNS publishes a fresh address update.</p>
					</div>
					<div className="metric-icon text-violet-300">
						<GearIcon />
					</div>
				</div>

				<label className="inline-flex items-center gap-3 text-sm theme-muted">
					<input
						type="checkbox"
						className="toggle toggle-primary"
						checked={props.Enabled}
						onChange={() => {
							if (props.onWebHookChange) {
								props.onWebHookChange({
									Enabled: !props.Enabled,
									Url: props.Url,
									RequestBody: props.RequestBody
								});
							}
						}}
					/>
					<span>Enable webhook delivery</span>
				</label>

				<fieldset className="theme-field">
					<label className="theme-field-label" htmlFor="webhook-url">Target URL</label>
					<input
						id="webhook-url"
						type="text"
						className="input theme-input w-full rounded-2xl"
						placeholder="Input the webhook URL"
						value={props.Url}
						disabled={!props.Enabled}
						onChange={(e) => {
							if (props.onWebHookChange) {
								props.onWebHookChange({
									Enabled: props.Enabled,
									Url: e.target.value,
									RequestBody: props.RequestBody
								});
							}
						}}
					/>
				</fieldset>

				<fieldset className="theme-field">
					<label className="theme-field-label" htmlFor="webhook-request-body">Request body</label>
					<textarea
						id="webhook-request-body"
						className="textarea theme-input h-28 w-full rounded-2xl"
						placeholder="Input request body"
						value={props.RequestBody}
						disabled={!props.Enabled}
						onChange={(e) => {
							if (props.onWebHookChange) {
								props.onWebHookChange({
									Enabled: props.Enabled,
									Url: props.Url,
									RequestBody: e.target.value
								});
							}
						}}
					/>
					<span className="theme-field-hint">Leave the body empty to send a GET request instead of POST.</span>
				</fieldset>
			</div>
		</div>
	);
};
