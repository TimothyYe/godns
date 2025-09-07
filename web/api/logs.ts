import { get_api_server } from '@/api/env';

export interface LogEntry {
	timestamp: string;
	level: string;
	message: string;
	fields?: { [key: string]: any };
}

export interface LogsResponse {
	logs: LogEntry[];
	total: number;
}

export async function get_logs(credentials: string, limit?: number, level?: string): Promise<LogsResponse | null> {
	if (!credentials) {
		return null;
	}

	const params = new URLSearchParams();
	if (limit) params.append('limit', limit.toString());
	if (level) params.append('level', level);

	try {
		const resp = await fetch(get_api_server() + '/api/v1/logs?' + params.toString(), {
			method: 'GET',
			headers: {
				'Authorization': `Basic ${credentials}`
			}
		});

		if (resp.status === 200) {
			return resp.json();
		}
	} catch (error) {
		console.error('Error fetching logs:', error);
	}

	return null;
}

export async function clear_logs(credentials: string): Promise<boolean> {
	if (!credentials) {
		return false;
	}

	try {
		const resp = await fetch(get_api_server() + '/api/v1/logs', {
			method: 'DELETE',
			headers: {
				'Authorization': `Basic ${credentials}`
			}
		});

		return resp.status === 200;
	} catch (error) {
		console.error('Error clearing logs:', error);
	}

	return false;
}

export async function get_log_levels(credentials: string): Promise<string[]> {
	if (!credentials) {
		return [];
	}

	try {
		const resp = await fetch(get_api_server() + '/api/v1/logs/levels', {
			method: 'GET',
			headers: {
				'Authorization': `Basic ${credentials}`
			}
		});

		if (resp.status === 200) {
			const data = await resp.json();
			return data.levels || [];
		}
	} catch (error) {
		console.error('Error fetching log levels:', error);
	}

	return [];
}