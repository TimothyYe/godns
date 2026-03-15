'use client';
import { useEffect, useState } from "react";
import { SunFilledIcon, MoonFilledIcon } from "@/components/icons";

export const ThemeSwitch = () => {
	const [theme, setTheme] = useState<string | null>(null);
	const [mounted, setMounted] = useState(false);

	useEffect(() => {
		setMounted(true);
		const localTheme = localStorage.getItem("theme");
		setTheme(localTheme ? localTheme : "dark");
	}, []);

	useEffect(() => {
		if (theme) {
			// set theme attribute to the <html> tag
			document.documentElement.setAttribute(
				"data-theme",
				theme
			);
		}
	}, [theme]);

	// Prevent hydration mismatch by not rendering until mounted
	if (!mounted || !theme) {
		return (
			<div className="inline-flex h-10 w-10 items-center justify-center rounded-full border border-base-300">
				<div className="w-[22px] h-[22px]" /> {/* Placeholder to prevent layout shift */}
			</div>
		);
	}

	return (
		<button
			type="button"
			className="inline-flex h-10 w-10 items-center justify-center rounded-full border border-base-300 text-base-content/70 transition-colors hover:border-base-content/20 hover:text-base-content"
			aria-label={`Switch to ${theme === "light" ? "dark" : "light"} mode`}
			onClick={() => {
				const newTheme = theme === "light" ? "dark" : "light";
				localStorage.setItem("theme", newTheme);
				setTheme(newTheme);
			}}
		>
			{theme === "dark" ? <SunFilledIcon className="h-[22px] w-[22px]" size={22} /> : <MoonFilledIcon size={22} className="h-[22px] w-[22px]" />}
		</button>
	);
};
