'use client'
import { useContext } from "react";
import { CommonContext } from '@/components/user';
import { useState, useEffect } from "react";

export const LogoutBtn = () => {
	const { logoutUser } = useContext(CommonContext);
	const [isClient, setIsClient] = useState(false);

	useEffect(() => {
		// Set isClient to true once the component has mounted
		setIsClient(typeof window !== 'undefined');
	}, []);
	const onClick = () => {
		// logout user
		logoutUser();
		// Redirect to the login page
		window.location.href = '/login';
	}

	return (
		isClient ? <button className="btn btn-outline btn-sm" onClick={onClick}>Logout</button> : null
	);
}