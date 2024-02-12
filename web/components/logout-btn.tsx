'use client';
import { useContext } from "react";
import { CommonContext } from '@/components/user';

export const LogoutBtn = () => {
	const { logoutUser } = useContext(CommonContext);
	const onClick = () => {
		// logout user
		logoutUser();
		// Redirect to the login page
		window.location.href = '/login';
	}

	return (
		<button className="btn btn-outline btn-sm" onClick={onClick}>Logout</button>
	);
}