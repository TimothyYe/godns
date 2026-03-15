interface WebHookProps {
	Enabled: boolean;
	Url: string;
	RequestBody: string;
	onWebHookChange?: (data: WebHookProps) => void;
}

export const WebHook = (props: WebHookProps) => {
	return (
		<div className="grid gap-5">
			<div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
				<div className="space-y-1">
					<h3 className="text-lg font-semibold tracking-tight">Webhook notifications</h3>
					<p className="text-sm text-base-content/60">
						Send a webhook when GoDNS completes updates or needs to signal external automation.
					</p>
				</div>
				<label className="flex items-center gap-3 rounded-full border border-base-300 px-4 py-2">
					<span className="text-sm font-medium">Enable webhook</span>
					<input
						type="checkbox"
						className="toggle toggle-primary"
						checked={props.Enabled}
						onChange={(e) => {
							props.onWebHookChange?.({
								Enabled: e.target.checked,
								Url: props.Url,
								RequestBody: props.RequestBody
							});
						}}
					/>
				</label>
			</div>

			<label className="field-stack">
				<span className="field-label">Webhook URL</span>
				<input
					type="text"
					className="input input-bordered h-12 w-full rounded-2xl"
					placeholder="https://example.com/hooks/godns"
					value={props.Url}
					disabled={!props.Enabled}
					onChange={(e) => {
						props.onWebHookChange?.({
							Enabled: props.Enabled,
							Url: e.target.value,
							RequestBody: props.RequestBody
						});
					}}
				/>
				<span className="field-hint">The destination that should receive GoDNS webhook calls.</span>
			</label>

			<label className="field-stack">
				<span className="field-label">Request body</span>
				<textarea
					className="textarea textarea-bordered min-h-32 w-full rounded-[1.25rem]"
					placeholder='{"event":"godns.update"}'
					value={props.RequestBody}
					disabled={!props.Enabled}
					onChange={(e) => {
						props.onWebHookChange?.({
							Enabled: props.Enabled,
							Url: props.Url,
							RequestBody: e.target.value
						});
					}}
				/>
				<span className="field-hint">Optional payload template sent to the webhook endpoint.</span>
			</label>
		</div>
	)
}
