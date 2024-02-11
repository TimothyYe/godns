'use client';
import { createContext, useState, useEffect, ReactNode } from 'react';

export const UserContext = createContext({
	credentials: '',
	currentPage: '',
	setCurrentPage: (_: string) => { },
	loginUser: (_: string) => { },
	logoutUser: () => { },
});

interface UserProviderProps {
	children: ReactNode;
}

// user provider
export const UserProvider = ({ children }: UserProviderProps) => {
	const [credentials, setCredentials] = useState<string>('');
	const [currentPage, setCurrentPage] = useState<string>('Home');

	useEffect(() => {
		const localCredentials = localStorage.getItem('credentials');
		if (localCredentials) {
			setCredentials(localCredentials);
		}
	}, []);

	const loginUser = (credentials: string) => {
		setCredentials(credentials);
		localStorage.setItem('credentials', credentials);
	};

	const logoutUser = () => {
		setCredentials('');
		localStorage.removeItem('credentials');
	};

	return (
		<UserContext.Provider value={{ credentials, loginUser, logoutUser, currentPage, setCurrentPage }}>
			{children}
		</UserContext.Provider>
	);
};
