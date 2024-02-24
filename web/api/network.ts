import { get_api_server } from "./env";

export interface WebHook {
	enabled: boolean;
	url: string;
	request_body: string;
}

export interface NetworkSettings {
	ip_mode: string;
	ip_urls?: string[];
	ipv6_urls?: string[];
	use_proxy: boolean;
	skip_ssl_verify: boolean;
	socks5_proxy: string;
	webhook: WebHook;
	resolver: string;
	ip_interface: string;
}

export async function get_network_settings(credentials: string): Promise<NetworkSettings> {
	if (credentials) {
		const resp = await fetch(get_api_server() + '/api/v1/network', {
			method: 'GET',
			headers: {
				'Authorization': `Basic ${credentials}`
			}
		})

		if (resp.status === 200) {
			return resp.json();
		}
	}

	return {} as NetworkSettings;
}

export async function update_network_settings(credentials: string, settings: NetworkSettings): Promise<boolean> {
	if (credentials) {
		const resp = await fetch(get_api_server() + '/api/v1/network', {
			method: 'PUT',
			headers: {
				'Authorization': `Basic ${credentials}`,
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(settings)
		})

		if (resp.status === 200) {
			return true;
		}
	}

	return false;
}