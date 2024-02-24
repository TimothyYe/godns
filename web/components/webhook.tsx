import { GearIcon } from "./icons"
import { useState } from "react";

interface WebHookProps {
	Enabled: boolean;
	Url: string;
	RequestBody: string;
	onWebHookChange?: (data: WebHookProps) => void;
}

export const WebHook = (props: WebHookProps) => {
	const [webhookEnabled, setWebhookEnabled] = useState(props.Enabled);
	const [webhookUrl, setWebhookUrl] = useState(props.Url);
	const [webhookRequestBody, setWebhookRequestBody] = useState(props.RequestBody);

	return (
		<div className="stats shadow bg-primary-content stats-vertical lg:stats-horizontal">
			<div className="stat gap-2">
				<div className="stat-title">Webhook</div>
				<div className="flex flex-col gap-3">
					<div className="flex flex-row items-center justify-start gap-2">
						<span className="label-text text-slate-500 ">Enable Webhook</span>
						<input
							type="checkbox"
							className="toggle toggle-primary"
							checked={webhookEnabled}
							onClick={() => {
								setWebhookEnabled(!webhookEnabled);
								if (props.onWebHookChange) {
									props.onWebHookChange({
										Enabled: !webhookEnabled,
										Url: webhookUrl,
										RequestBody: webhookRequestBody
									});
								}
							}}
							onChange={() => { }}
						/>
						<div className="flex flex-grow justify-end text-secondary">
							<GearIcon />
						</div>
					</div>
					<input type="text"
						className="input input-primary w-full"
						placeholder="Input the webhhook URL"
						value={webhookUrl}
						disabled={!webhookEnabled}
						onChange={(e) => {
							setWebhookUrl(e.target.value);
							if (props.onWebHookChange) {
								props.onWebHookChange({
									Enabled: webhookEnabled,
									Url: e.target.value,
									RequestBody: webhookRequestBody
								});
							}
						}}
					/>
					<textarea
						className="textarea textarea-primary w-full h-28"
						placeholder="Input request body"
						value={props.RequestBody}
						disabled={!webhookEnabled}
						onChange={(e) => {
							setWebhookRequestBody(e.target.value);
							if (props.onWebHookChange) {
								props.onWebHookChange({
									Enabled: webhookEnabled,
									Url: webhookUrl,
									RequestBody: e.target.value
								});
							}
						}}
					/>
				</div>
			</div>
		</div>
	)
}