import { get_api_server } from "./env";

export interface ProviderSetting {
	name: string;
	username: boolean;
	email: boolean;
	password: boolean;
	login_token: boolean;
	app_key: boolean;
	app_secret: boolean;
	consumer_key: boolean;
}

export interface Provider {
	provider: string;
	username: string;
	email: string;
	password: string;
	login_token: string;
	app_key: string;
	app_secret: string;
	consumer_key: string;
}

export interface MultiProviderConfig {
	[providerName: string]: {
		email?: string;
		password?: string;
		login_token?: string;
		app_key?: string;
		app_secret?: string;
		consumer_key?: string;
	};
}

export async function get_provider_settings(credentials: string): Promise<ProviderSetting[]> {
	if (credentials) {
		const resp = await fetch(get_api_server() + '/api/v1/provider/settings', {
			method: 'GET',
			headers: {
				'Authorization': `Basic ${credentials}`
			}
		})

		if (resp.status === 200) {
			return resp.json();
		}
	}

	return {} as ProviderSetting[];
}

export async function get_provider(credentials: string): Promise<Provider> {
	if (credentials) {
		const resp = await fetch(get_api_server() + '/api/v1/provider', {
			method: 'GET',
			headers: {
				'Authorization': `Basic ${credentials}`
			}
		})

		if (resp.status === 200) {
			return resp.json();
		}
	}

	return {} as Provider;
}

export async function update_provider(credentials: string, provider: Provider): Promise<boolean> {
	if (credentials) {
		const resp = await fetch(get_api_server() + '/api/v1/provider', {
			method: 'PUT',
			headers: {
				'Authorization': `Basic ${credentials}`,
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(provider)
		})

		if (resp.status === 200) {
			return true;
		}
	}

	return false;
}

export async function get_multi_providers(credentials: string): Promise<MultiProviderConfig> {
	if (credentials) {
		const resp = await fetch(get_api_server() + '/api/v1/providers', {
			method: 'GET',
			headers: {
				'Authorization': `Basic ${credentials}`
			}
		})

		if (resp.status === 200) {
			return resp.json();
		}
	}

	return {} as MultiProviderConfig;
}

export async function update_multi_providers(credentials: string, providers: MultiProviderConfig): Promise<boolean> {
	if (credentials) {
		const resp = await fetch(get_api_server() + '/api/v1/providers', {
			method: 'PUT',
			headers: {
				'Authorization': `Basic ${credentials}`,
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(providers)
		})

		if (resp.status === 200) {
			return true;
		}
	}

	return false;
}

export async function add_provider_config(credentials: string, providerName: string, config: any): Promise<boolean> {
	if (credentials) {
		const resp = await fetch(get_api_server() + `/api/v1/providers/${providerName}`, {
			method: 'PUT',
			headers: {
				'Authorization': `Basic ${credentials}`,
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(config)
		})

		if (resp.status === 200) {
			return true;
		}
	}

	return false;
}

export async function delete_provider_config(credentials: string, providerName: string): Promise<boolean> {
	if (credentials) {
		const resp = await fetch(get_api_server() + `/api/v1/providers/${providerName}`, {
			method: 'DELETE',
			headers: {
				'Authorization': `Basic ${credentials}`
			}
		})

		if (resp.status === 200) {
			return true;
		}
	}

	return false;
}