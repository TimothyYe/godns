'use client'
import { siteConfig } from "@/config/site";
import { MenuIcon, GithubIcon, HeartFilledIcon } from "./icons";
import { LogoutBtn } from "./logout-btn";
import { ThemeSwitch } from "./theme-switch";
import { useContext } from "react";
import { CommonContext } from '@/components/user';
import { useState, useEffect } from "react";

export const Navbar = () => {
	const userStore = useContext(CommonContext);
	const { credentials, currentPage, version } = userStore;
	const [isClient, setIsClient] = useState(false);

	useEffect(() => {
		// Set isClient to true once the component has mounted
		setIsClient(true);
	}, []);

	return (
		<div className="navbar bg-base-100">
			<div className="navbar-start gap-2">
				<div className="dropdown">
					<div tabIndex={0} role="button" className="btn lg:hidden">
						<MenuIcon />
					</div>
					{isClient && credentials ?
						<ul tabIndex={0} className="menu menu-sm dropdown-content mt-3 z-[1] p-2 shadow rounded-box w-52 bg-white">
							{
								siteConfig.navItems.map((item) => (
									<li key={item.label}>
										<a href={item.href}>{item.label}</a>
									</li>
								))
							}
						</ul> : null}
				</div>
				<span className="text-2xl font-bold">GoDNS</span>
				<span className="text-sm mt-2">{isClient && version && version !== 'v0.1' ? `${version}` : ''}</span>
			</div>
			<div className="navbar-center hidden lg:flex">
				{
					isClient && credentials ?
						<ul className="menu menu-horizontal px-1">
							{
								siteConfig.navItems.map((item) => (
									<li key={item.label}>
										<a className={currentPage === item.label ? "font-semibold bg-slate-100" : "font-semibold"} href={item.href}>{item.label}</a>
									</li>
								))}
						</ul> : null
				}
			</div>
			<div className="hidden sm:flex navbar-end gap-2">
				<ThemeSwitch />
				<a className="hidden sm:flex link" href={siteConfig.links.github} target="_blank" aria-label="Github">
					<GithubIcon className="text-default-500" />
				</a>
				<a className="hidden sm:flex link" href={siteConfig.links.sponsor} target="_blank" aria-label="Sponsor">
					<HeartFilledIcon className="text-red-500" />
				</a>
				{
					credentials ? <LogoutBtn /> : null
				}
			</div>
		</div>
	);
}