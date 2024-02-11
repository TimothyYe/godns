import { siteConfig } from "@/config/site";
import { MenuIcon, GithubIcon, HeartFilledIcon } from "./icons";
import { LogoutBtn } from "./logout-btn";
import { ThemeSwitch } from "./theme-switch";

export const Navbar = () => {
	return (
		<div className="navbar bg-base-100">
			<div className="navbar-start">
				<div className="dropdown">
					<div tabIndex={0} role="button" className="btn lg:hidden">
						<MenuIcon />
					</div>
					<ul tabIndex={0} className="menu menu-sm dropdown-content mt-3 z-[1] p-2 shadow rounded-box w-52">
						{
							siteConfig.navItems.map((item) => (
								<li key={item.label}>
									<a className="text-primary" href={item.href}>{item.label}</a>
								</li>
							))
						}
					</ul>
				</div>
				<a className="btn text-xl">GoDNS</a>
			</div>
			<div className="navbar-center hidden lg:flex">
				<ul className="menu menu-horizontal px-1">
					{
						siteConfig.navItems.map((item) => (
							<li key={item.label}>
								<a className="text-primary" href={item.href}>{item.label}</a>
							</li>
						))
					}
				</ul>
			</div>
			<div className="navbar-end gap-2">
				<ThemeSwitch />
				<a className="hidden sm:flex link" href={siteConfig.links.github} aria-label="Github">
					<GithubIcon className="text-default-500" />
				</a>
				<a className="hidden sm:flex link" href={siteConfig.links.sponsor} aria-label="Sponsor">
					<HeartFilledIcon className="text-red-500" />
				</a>
				<LogoutBtn />
			</div>
		</div>
	);
}