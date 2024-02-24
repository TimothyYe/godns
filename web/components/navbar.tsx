'use client'
import { siteConfig } from "@/config/site";
import { MenuIcon, GithubIcon, HeartFilledIcon } from "./icons";
import { LogoutBtn } from "./logout-btn";
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
				<label className="flex cursor-pointer gap-2 items-center">
					<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><circle cx="12" cy="12" r="5" /><path d="M12 1v2M12 21v2M4.2 4.2l1.4 1.4M18.4 18.4l1.4 1.4M1 12h2M21 12h2M4.2 19.8l1.4-1.4M18.4 5.6l1.4-1.4" /></svg>
					<input type="checkbox" value="dark" className="toggle theme-controller" />
					<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"></path></svg>
				</label>
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