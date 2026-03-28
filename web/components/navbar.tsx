'use client'
import Link from "next/link";
import classNames from "classnames";
import { siteConfig } from "@/config/site";
import { MenuIcon, GithubIcon, HeartFilledIcon } from "./icons";
import { LogoutBtn } from "./logout-btn";
import { useContext } from "react";
import { CommonContext } from '@/components/user';
import { ThemeSwitch } from "./theme-switch";
import { useIsHydrated } from "./use-is-hydrated";

export const Navbar = () => {
	const userStore = useContext(CommonContext);
	const { credentials, currentPage, version } = userStore;
	const isHydrated = useIsHydrated();

	return (
		<header className="theme-nav sticky top-0 z-40 backdrop-blur-xl">
			<div className="container mx-auto flex max-w-7xl items-center justify-between gap-3 px-4 py-4 sm:px-6">
				<div className="flex items-center gap-3">
					<div className="dropdown">
						<div tabIndex={0} role="button" className="theme-icon-btn btn btn-ghost btn-sm rounded-xl px-3 lg:hidden">
							<MenuIcon />
						</div>
						{isHydrated && credentials ? (
							<ul tabIndex={0} className="menu menu-sm dropdown-content surface-panel z-[1] mt-3 w-56 p-2">
								{siteConfig.navItems.map((item) => (
									<li key={item.label}>
										<Link
											href={item.href}
											className={classNames(
												"rounded-xl px-4 py-2.5 text-sm font-medium transition-colors",
												currentPage === item.label
													? "theme-nav-link-active"
													: "theme-nav-link"
											)}
										>
											{item.label}
										</Link>
									</li>
								))}
							</ul>
						) : null}
					</div>

					<Link href="/" className="flex items-center gap-3">
						<div className="flex h-10 w-10 items-center justify-center rounded-2xl bg-gradient-to-br from-sky-500 to-blue-600 text-base font-bold text-white shadow-lg shadow-sky-900/30">
							G
						</div>
						<div className="flex items-baseline gap-3">
							<span className="text-2xl font-semibold tracking-tight theme-heading">GoDNS</span>
							{isHydrated && version && version !== 'v0.1' ? (
								<span className="theme-chip hidden rounded-full px-2.5 py-1 text-xs font-medium sm:inline-flex">
									{version}
								</span>
							) : null}
						</div>
					</Link>
				</div>

				<div className="navbar-center hidden lg:flex">
					{isHydrated && credentials ? (
						<ul className="menu menu-horizontal rounded-2xl p-1.5 form-shell">
							{siteConfig.navItems.map((item) => (
								<li key={item.label}>
									<Link
										className={classNames(
											"rounded-xl px-4 py-2.5 text-sm font-medium transition-colors",
											currentPage === item.label
												? "theme-nav-link-active"
												: "theme-nav-link"
										)}
										href={item.href}
									>
										{item.label}
									</Link>
								</li>
							))}
						</ul>
					) : null}
				</div>

				<div className="flex items-center gap-2">
					<ThemeSwitch />
					<a className="theme-icon-btn theme-nav-utility-btn hidden h-10 w-10 items-center justify-center rounded-xl transition-colors sm:inline-flex" href={siteConfig.links.github} target="_blank" aria-label="Github">
						<GithubIcon className="h-5 w-5" />
					</a>
					<a className="theme-icon-btn theme-nav-utility-btn hidden h-10 w-10 items-center justify-center rounded-xl transition-colors sm:inline-flex" href={siteConfig.links.sponsor} target="_blank" aria-label="Sponsor">
						<HeartFilledIcon className="h-5 w-5" />
					</a>
					{credentials ? <LogoutBtn /> : null}
				</div>
			</div>
		</header>
	);
}
