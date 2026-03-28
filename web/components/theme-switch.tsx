'use client';
import { useEffect, useState } from "react";
import { SunFilledIcon, MoonFilledIcon } from "@/components/icons";

const getPreferredTheme = () => {
	if (typeof window === 'undefined') {
		return 'light';
	}

	const localTheme = localStorage.getItem("theme");
	if (localTheme === 'light' || localTheme === 'dark') {
		return localTheme;
	}

	return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
};

export const ThemeSwitch = () => {
	const [theme, setTheme] = useState<string | null>(null);
	const [mounted, setMounted] = useState(false);

	useEffect(() => {
		setMounted(true);
		setTheme(getPreferredTheme());
	}, []);

	useEffect(() => {
		if (theme) {
			document.documentElement.setAttribute("data-theme", theme);
			localStorage.setItem("theme", theme);
		}
	}, [theme]);

	// Prevent hydration mismatch by not rendering until mounted
	if (!mounted || !theme) {
		return (
			<div className="theme-icon-btn theme-nav-utility-btn flex h-10 w-10 items-center justify-center rounded-xl">
				<div className="w-[22px] h-[22px]" /> {/* Placeholder to prevent layout shift */}
			</div>
		);
	}

	return (
		<button
			type="button"
			className="theme-icon-btn theme-nav-utility-btn flex h-10 w-10 items-center justify-center rounded-xl transition-colors"
			aria-label={`Switch to ${theme === "light" ? "dark" : "light"} mode`}
			title={theme === "light" ? "Switch to dark mode" : "Switch to light mode"}
			onClick={() => {
				const newTheme = theme === "light" ? "dark" : "light";
				setTheme(newTheme);
			}}
		>
			{theme === "dark" ? <SunFilledIcon size={20} /> : <MoonFilledIcon size={20} />}
		</button>
	);
};
