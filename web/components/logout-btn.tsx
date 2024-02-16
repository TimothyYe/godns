'use client'
import { useContext } from "react";
import { useRouter } from "next/navigation";
import { CommonContext } from '@/components/user';
import { useState, useEffect } from "react";

export const LogoutBtn = () => {
	const router = useRouter();
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
		router.push('/login');
	}

	return (
		isClient ? <button className="btn btn-outline btn-sm" onClick={onClick}>Logout</button> : null
	);
}