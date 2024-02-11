export type SiteConfig = typeof siteConfig;

export const siteConfig = {
	name: "GoDNS",
	description: "Dynamic DNS client with multiple providers support",
	navItems: [
		{
			label: "Status",
			href: "/",
		},
		{
			label: "Config",
			href: "/config",
		}
	],
	links: {
		github: "https://github.com/TimothyYe/godns",
		sponsor: "https://github.com/sponsors/TimothyYe"
	},
};