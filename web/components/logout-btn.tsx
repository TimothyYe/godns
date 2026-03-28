'use client'
import { useContext } from "react";
import { useRouter } from "next/navigation";
import { CommonContext } from '@/components/user';
import { useIsHydrated } from "./use-is-hydrated";

export const LogoutBtn = () => {
	const router = useRouter();
	const { logoutUser } = useContext(CommonContext);
	const isHydrated = useIsHydrated();
	const onClick = () => {
		// logout user
		logoutUser();
		// Redirect to the login page
		router.push('/login');
	}

	return (
		isHydrated ? <button className="theme-nav-logout-btn btn btn-sm rounded-xl px-4" onClick={onClick}>Logout</button> : null
	);
}
