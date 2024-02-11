export function get_api_server(): string {
	console.log('API_SERVER:', process.env.NEXT_PUBLIC_API_SERVER);
	return process.env.NEXT_PUBLIC_API_SERVER || '';
};