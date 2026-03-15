'use client';

import Link from "next/link";
import classNames from "classnames";
import { useContext, useEffect, useState } from "react";
import { siteConfig } from "@/config/site";
import { MenuIcon, GithubIcon, HeartFilledIcon } from "./icons";
import { LogoutBtn } from "./logout-btn";
import { ThemeSwitch } from "./theme-switch";
import { CommonContext } from '@/components/user';

export const Navbar = () => {
	const userStore = useContext(CommonContext);
	const { credentials, currentPage, version } = userStore;
	const [isClient, setIsClient] = useState(false);

	useEffect(() => {
		setIsClient(true);
	}, []);

	return (
		<header className="sticky top-0 z-30 border-b border-base-300/60 bg-base-100/88 backdrop-blur-xl">
			<div className="mx-auto flex w-full max-w-7xl items-center justify-between gap-3 px-4 py-4 sm:px-6 lg:px-8">
				<div className="flex items-center gap-3">
					<div className="dropdown">
						<div tabIndex={0} role="button" className="btn btn-ghost btn-circle border border-base-300 lg:hidden">
							<MenuIcon />
						</div>
						{isClient && credentials ? (
							<ul tabIndex={0} className="menu menu-sm dropdown-content mt-3 z-[1] w-64 rounded-2xl border border-base-300 bg-base-100 p-2 shadow-xl">
								{siteConfig.navItems.map((item) => (
									<li key={item.label}>
										<Link href={item.href}>{item.label}</Link>
									</li>
								))}
							</ul>
						) : null}
					</div>

					<Link href="/" className="flex items-center gap-3">
						<div className="flex h-10 w-10 items-center justify-center rounded-2xl bg-primary text-primary-content shadow-lg shadow-primary/20">
							<span className="text-lg font-black">G</span>
						</div>
						<div className="space-y-0.5">
							<p className="text-lg font-semibold tracking-tight text-base-content">GoDNS</p>
							<p className="text-xs uppercase tracking-[0.24em] text-base-content/50">Control Panel</p>
						</div>
					</Link>

					{isClient && version && version !== 'v0.1' ? (
						<span className="hidden rounded-full border border-base-300 bg-base-200 px-2.5 py-1 text-xs font-medium text-base-content/70 sm:inline-flex">
							{version}
						</span>
					) : null}
				</div>

				{isClient && credentials ? (
					<nav className="hidden lg:flex">
						<ul className="menu menu-horizontal gap-1 rounded-full border border-base-300/80 bg-base-100/80 p-1">
							{siteConfig.navItems.map((item) => (
								<li key={item.label}>
									<Link
										className={classNames(
											"rounded-full px-4 py-2 font-medium transition-colors",
											currentPage === item.label
												? "bg-primary text-primary-content"
												: "text-base-content/70 hover:bg-base-200 hover:text-base-content"
										)}
										href={item.href}
									>
										{item.label}
									</Link>
								</li>
							))}
						</ul>
					</nav>
				) : (
					<div />
				)}

				<div className="flex items-center gap-2">
					<ThemeSwitch />
					<a className="hidden rounded-full border border-base-300 p-2 text-base-content/60 transition-colors hover:border-base-content/20 hover:text-base-content sm:inline-flex" href={siteConfig.links.github} target="_blank" aria-label="Github">
						<GithubIcon className="h-5 w-5" />
					</a>
					<a className="hidden rounded-full border border-base-300 p-2 text-red-500 transition-colors hover:border-red-200 hover:text-red-600 sm:inline-flex" href={siteConfig.links.sponsor} target="_blank" aria-label="Sponsor">
						<HeartFilledIcon className="h-5 w-5" />
					</a>
					{credentials ? <LogoutBtn /> : null}
				</div>
			</div>
		</header>
	);
};
