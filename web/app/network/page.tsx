'use client';

import Link from "next/link";
import { useRouter } from "next/navigation";
import { CommonContext } from "@/components/user";
import { useEffect, useMemo, useState, useContext } from "react";
import { get_network_settings, NetworkSettings, update_network_settings } from "@/api/network";
import { get_info } from "@/api/info";
import { toast } from 'react-toastify';
import { IpMode } from "@/components/ip-mode";
import { Proxy } from "@/components/proxy";
import { WebHook } from "@/components/webhook";
import { Resolver } from "@/components/resolver";
import { IPInterface } from "@/components/ip-interface";
import { PageShell, SectionCard } from "@/components/page-shell";

const emptySettings: NetworkSettings = {
	ip_mode: 'IPV4',
	ip_urls: [],
	ipv6_urls: [],
	use_proxy: false,
	skip_ssl_verify: false,
	socks5_proxy: '',
	webhook: {
		enabled: false,
		url: '',
		request_body: ''
	},
	resolver: '',
	ip_interface: ''
};

export default function Network() {
	const router = useRouter();
	const userStore = useContext(CommonContext);
	const { credentials, setCurrentPage, saveVersion } = userStore;
	const [settings, setSettings] = useState<NetworkSettings>(emptySettings);
	const [initialSettings, setInitialSettings] = useState<NetworkSettings>(emptySettings);
	const [loading, setLoading] = useState(true);
	const [saving, setSaving] = useState(false);

	useEffect(() => {
		if (!credentials) {
			router.push('/login');
			return;
		}

		setCurrentPage('Network');
		setLoading(true);
		Promise.all([
			get_info(credentials),
			get_network_settings(credentials),
		]).then(([info, nextSettings]) => {
			saveVersion(info.version);
			const normalizedSettings = {
				...emptySettings,
				...nextSettings,
				webhook: {
					...emptySettings.webhook,
					...nextSettings.webhook,
				}
			};
			setSettings(normalizedSettings);
			setInitialSettings(normalizedSettings);
		}).finally(() => {
			setLoading(false);
		});
	}, [credentials, router, saveVersion, setCurrentPage]);

	const hasChanges = useMemo(
		() => JSON.stringify(settings) !== JSON.stringify(initialSettings),
		[initialSettings, settings]
	);

	const saveSettings = async () => {
		if (!credentials) {
			toast.error('Invalid credentials');
			return;
		}

		setSaving(true);
		const success = await update_network_settings(credentials, settings);
		setSaving(false);

		if (!success) {
			toast.error('Failed to update network settings');
			return;
		}

		setInitialSettings(settings);
		toast.success('Network settings updated successfully');
	};

	return (
		<PageShell
			eyebrow="Networking"
			title="How GoDNS reaches the internet"
			description="Tune public IP discovery, network access, and outbound integration behavior. These settings affect how GoDNS resolves your address and talks to provider APIs."
			actions={(
				<>
					<Link className="btn btn-ghost rounded-full px-5" href="/logs">View logs</Link>
					<button className="btn btn-primary rounded-full px-5" onClick={saveSettings} disabled={!hasChanges || saving || loading}>
						{saving ? 'Saving...' : hasChanges ? 'Save changes' : 'Saved'}
					</button>
				</>
			)}
		>
			{loading ? (
				<div className="panel flex min-h-[24rem] items-center justify-center">
					<span className="loading loading-spinner loading-lg" />
				</div>
			) : (
				<>
					<SectionCard
						title="IP detection"
						description="Choose which IP family GoDNS should detect and which endpoints it should query to resolve your public address."
					>
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
					</SectionCard>

					<SectionCard
						title="Outbound access"
						description="Adjust how GoDNS reaches provider APIs and performs DNS/network resolution."
					>
						<div className="grid gap-6 lg:grid-cols-[1.25fr_0.75fr]">
							<div className="rounded-[1.5rem] border border-base-300/70 p-5">
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
							</div>
							<div className="grid gap-5 rounded-[1.5rem] border border-base-300/70 p-5">
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
							</div>
						</div>
					</SectionCard>

					<SectionCard
						title="Notifications"
						description="Optional outbound hooks for integrating GoDNS updates with other systems."
					>
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
					</SectionCard>

					<div className="panel flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
						<div className="space-y-1">
							<p className="text-lg font-semibold">Save before leaving</p>
							<p className="text-sm text-base-content/60">
								GoDNS uses these values directly during runtime. Unsaved edits are not applied.
							</p>
						</div>
						<div className="flex flex-wrap gap-2">
							<button
								className="btn rounded-full px-5"
								onClick={() => setSettings(initialSettings)}
								disabled={!hasChanges || saving}
							>
								Reset edits
							</button>
							<button
								className="btn btn-primary rounded-full px-5"
								onClick={saveSettings}
								disabled={!hasChanges || saving}
							>
								{saving ? 'Saving...' : 'Save settings'}
							</button>
						</div>
					</div>
				</>
			)}
		</PageShell>
	);
}
