'use client';
import { useEffect, useState } from "react";
import { SunFilledIcon, MoonFilledIcon } from "@/components/icons";

export const ThemeSwitch = () => {
	const [theme, setTheme] = useState<string>("light");

	useEffect(() => {
		if (typeof window !== 'undefined') {
			const localTheme = localStorage.getItem("theme");
			setTheme(localTheme ? localTheme : "light");
		}
	}, []);

	useEffect(() => {
		document.documentElement.setAttribute(
			"data-theme",
			theme
		);
	}, [theme]);

	return (
		<div className="w-auto h-auto bg-transparent rounded-lg flex items-center justify-center group-data-[selected=true]:bg-transparent !text-default-500 pt-px px-0 mx-0">
			<div onClick={
				() => {
					const newTheme = theme === "light" ? "night" : "light";
					localStorage.setItem("theme", newTheme);
					setTheme(newTheme);
				}
			}>
				{theme === "light" ? <SunFilledIcon className="hover:text-gray-700" size={22} /> : <MoonFilledIcon size={22} className="hover:text-gray-700" />}
			</div>
		</div>
	);
};
