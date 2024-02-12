export function get_api_server(): string {
	return process.env.NEXT_PUBLIC_API_SERVER || '';
};