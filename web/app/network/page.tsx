'use client';

import { IpMode } from "@/components/ip-mode";
import { Proxy } from "@/components/proxy";
import { WebHook } from "@/components/webhook";
import { Resolver } from "@/components/resolver";
import { IPInterface } from "@/components/ip-interface";
import { useRouter } from "next/navigation";
import { CommonContext } from "@/components/user";
import { useEffect, useState, useContext } from "react";
import { get_network_settings, NetworkSettings, update_network_settings } from "@/api/network";
import { get_info } from "@/api/info";
import { ToastContainer, toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';

export default function Network() {
	const router = useRouter();
	const userStore = useContext(CommonContext);
	const { credentials, setCurrentPage, saveVersion } = userStore;
	const [settings, setSettings] = useState<NetworkSettings>({} as NetworkSettings);

	useEffect(() => {
		if (!credentials) {
			router.push('/login');
			return;
		}
		setCurrentPage('Network');
		get_info(credentials).then((info) => {
			saveVersion(info.version);
		});

		get_network_settings(credentials).then((settings) => {
			setSettings(settings);
		});
	}, [credentials, router, setCurrentPage, saveVersion]);


	return (
		<main className="page-wrap">
			<ToastContainer />
			<div className="page-shell">
				<section className="page-hero page-hero-compact">
					<div className="eyebrow">
						<span className="inline-block h-2 w-2 rounded-full bg-amber-400" />
						Network Controls
					</div>
					<h1 className="page-title">Tune how GoDNS discovers and publishes addresses.</h1>
				</section>

				<section className="section-shell">
					<div className="section-header">
						<div className="section-copy">
							<div className="section-label ml-0">Configuration</div>
							<h2 className="text-2xl font-semibold tracking-tight theme-heading">Network and webhook settings</h2>
							<p className="mt-2 text-sm leading-7 theme-muted">
								These controls affect how public addresses are detected and how update events are pushed to downstream systems.
							</p>
						</div>
					</div>

					<div className="mt-5 flex flex-col gap-5">
						<IpMode
							IPMode={settings.ip_mode}
							IPUrls={settings.ip_urls}
							IPV6Urls={settings.ipv6_urls}
							onIpModeChange={(data) => {
								setSettings({
									...settings,
									ip_mode: data.IPMode,
									ip_urls: data.IPUrls,
									ipv6_urls: data.IPV6Urls
								});
							}}
						/>
						<Proxy
							EnableProxy={settings.use_proxy}
							SkipSSLVerify={settings.skip_ssl_verify}
							Socks5Proxy={settings.socks5_proxy}
							onProxyChange={(data) => {
								setSettings({
									...settings,
									use_proxy: data.EnableProxy,
									skip_ssl_verify: data.SkipSSLVerify,
									socks5_proxy: data.Socks5Proxy
								});
							}}
						/>
						{settings.webhook ? (
							<WebHook
								Enabled={settings.webhook.enabled}
								Url={settings.webhook.url}
								RequestBody={settings.webhook.request_body}
								onWebHookChange={(data) => {
									setSettings({
										...settings,
										webhook: {
											enabled: data.Enabled,
											url: data.Url,
											request_body: data.RequestBody
										}
									});
								}}
							/>
						) : null}
						<Resolver
							Resolver={settings.resolver}
							onResolverChange={(data) => {
								setSettings({
									...settings,
									resolver: data.Resolver
								});
							}}
						/>
						<IPInterface
							IPInterface={settings.ip_interface}
							onIPInterfaceChange={(data) => {
								setSettings({
									...settings,
									ip_interface: data.IPInterface
								});
							}}
						/>
						<div className="flex justify-end">
							<button
								className="theme-primary-sky btn rounded-xl border-none px-6"
								onClick={() => {
									if (!credentials) {
										toast.error('Invalid credentials');
										return;
									}

									update_network_settings(credentials, settings).then((success) => {
										if (success) {
											toast.success('Network settings updated successfully');
										} else {
											toast.error('Failed to update network settings');
										}
									});
								}}
							>
								Save
							</button>
						</div>
					</div>
				</section>
			</div>
		</main>
	);
}
