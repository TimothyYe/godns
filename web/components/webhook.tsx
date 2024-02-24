import { GearIcon } from "./icons"
import { useState, useEffect } from "react";

interface WebHookProps {
	Enabled: boolean;
	Url: string;
	RequestBody: string;
	onWebHookChange?: (data: WebHookProps) => void;
}

export const WebHook = ({ Enabled, Url, RequestBody, onWebHookChange }: WebHookProps) => {
	const [webhookEnabled, setWebhookEnabled] = useState(Enabled);
	const [webhookUrl, setWebhookUrl] = useState(Url);
	const [webhookRequestBody, setWebhookRequestBody] = useState(RequestBody);

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
								if (onWebHookChange) {
									onWebHookChange({
										Enabled: !webhookEnabled,
										Url: webhookUrl,
										RequestBody: webhookRequestBody
									});
								}
							}}
						/>
						<div className="flex flex-grow justify-end text-secondary">
							<GearIcon />
						</div>
					</div>
					<input type="text"
						className="input input-primary w-full"
						placeholder="Input the webhhook URL"
						value={Url}
						disabled={!webhookEnabled}
						onChange={(e) => {
							setWebhookUrl(e.target.value);
							if (onWebHookChange) {
								onWebHookChange({
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
						value={RequestBody}
						disabled={!webhookEnabled}
						onChange={(e) => {
							setWebhookRequestBody(e.target.value);
							if (onWebHookChange) {
								onWebHookChange({
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