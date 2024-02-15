import { get_api_server } from '@/api/env';
import { Domain } from '@/api/domain';

export interface Info {
	version: string;
	start_time: number;
	domain_num: number;
	sub_domain_num: number;
	domains: Domain[];
	public_ip: string;
	ip_mode: string;
	provider: string;
}

export async function get_info(credentials: string): Promise<Info> {
	if (credentials) {
		// make a GET request to the /api/auth endpoint via basic authentication
		const resp = await fetch(get_api_server() + '/api/v1/info', {
			method: 'GET',
			headers: {
				'Authorization': `Basic ${credentials}`
			}
		})

		if (resp.status === 200) {
			return resp.json();
		}
	}

	return {} as Info;
}

export function get_hours(timestamp: number): string {
	if (timestamp === 0) {
		return 'N/A';
	}
	// compute the number of hours between the current time and the timestamp
	const current_time = Date.now() / 1000;
	const diff = (current_time - timestamp);
	const hours = (diff / 3600).toFixed(1);
	return `${hours} Hours`;
}

export function get_date(timestamp: number): string {
	// convert the timestamp to a human-readable date
	const date = new Date(timestamp * 1000);
	// convert date to YYYY-MM-DD HH:mm:ss format with the time in local timezone
	return date.toLocaleString('en-US');
}