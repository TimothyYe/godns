export type SiteConfig = typeof siteConfig;

export const siteConfig = {
	name: "GoDNS",
	description: "Dynamic DNS client with multiple providers support",
	navItems: [
		{
			label: "Home",
			href: "/",
		},
		{
			label: "Domains",
			href: "/domains",
		},
		{
			label: "Network",
			href: "/network",
		},
	],
	links: {
		github: "https://github.com/TimothyYe/godns",
		sponsor: "https://github.com/sponsors/TimothyYe"
	},
};