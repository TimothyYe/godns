import { get_api_server } from '@/api/env';

export async function login(username: string, password: string): Promise<string> {
	// combine username and password into a single string and use Base64 encoding
	const credentials = `${username}:${password}`;
	const encodedCredentials = btoa(credentials);

	// make a GET request to the /api/auth endpoint via basic authentication
	const resp = await fetch(get_api_server() + '/api/v1/auth', {
		method: 'GET',
		headers: {
			'Authorization': `Basic ${encodedCredentials}`
		}
	})

	if (resp.status === 200) {
		// store the credentials in local storage
		return encodedCredentials;
	}

	return '';
}