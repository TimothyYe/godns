import { get_api_server } from '@/api/env';

export interface Domain {
	domain_name: string;
	sub_domains: string[];
}

export async function get_domains(credentials: string): Promise<Domain[]> {
	if (credentials) {
		// make a GET request to the /api/auth endpoint via basic authentication
		const resp = await fetch(get_api_server() + '/api/v1/domains', {
			method: 'GET',
			headers: {
				'Authorization': `Basic ${credentials}`
			}
		})

		if (resp.status === 200) {
			return resp.json();
		}
	}

	return {} as Domain[];
}

export async function add_domain(credentials: string, domain: Domain): Promise<boolean> {
	if (credentials) {
		// make a POST request to the /api/auth endpoint via basic authentication
		const resp = await fetch(get_api_server() + '/api/v1/domains', {
			method: 'POST',
			headers: {
				'Authorization': `Basic ${credentials}`,
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(domain)
		})

		if (resp.status === 200) {
			return true;
		}
	}

	return false;
}

export async function remove_domain(credentials: string, domain: string): Promise<boolean> {
	if (credentials) {
		// make a DELETE request to the /api/auth endpoint via basic authentication
		const resp = await fetch(get_api_server() + '/api/v1/domains/' + domain, {
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