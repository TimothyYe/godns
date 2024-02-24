'use client';

import { IpMode } from "@/components/ip-mode";
import { Proxy } from "@/components/proxy";
import { WebHook } from "@/components/webhook";
import { Resolver } from "@/components/resolver";
import { IPInterface } from "@/components/ip-interface";
import { useRouter } from "next/navigation";
import { CommonContext } from "@/components/user";
import { useEffect, useState, useContext } from "react";
import { get_network_settings, NetworkSettings, WebHook as WebHookSetting, update_network_settings } from "@/api/network";
import { get_info } from "@/api/info";

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
			console.log('settings', settings);
			setSettings(settings);
		});
	}, [credentials, router, setCurrentPage, saveVersion]);


	return (
		<main className="flex min-h-screen flex-col items-center justify-between pt-10 max-w-screen-xl">
			<div className="p-5">
				<div className="flex flex-col max-w-screen-lg gap-5">
					<IpMode IPMode={settings.ip_mode} IPUrls={settings.ip_urls} IPV6Urls={settings.ipv6_urls} />
					<Proxy />
					<WebHook />
					<Resolver />
					<IPInterface />
					<div className="flex justify-center">
						<button className="flex btn btn-primary">Save</button>
					</div>
				</div>
			</div>
		</main>
	);
}