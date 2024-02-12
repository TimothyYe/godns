'use client';
import { createContext, useState, useEffect, ReactNode } from 'react';

export const CommonContext = createContext({
	credentials: '',
	loginUser: (_: string) => { },
	logoutUser: () => { },
	currentPage: '',
	setCurrentPage: (_: string) => { },
	version: '',
	setVersion: (_: string) => { },
});

interface UserProviderProps {
	children: ReactNode;
}

// user provider
export const UserProvider = ({ children }: UserProviderProps) => {
	const [credentials, setCredentials] = useState<string>('');
	const [currentPage, setCurrentPage] = useState<string>('Home');
	const [version, setVersion] = useState<string>('');

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
		<CommonContext.Provider value={{
			credentials,
			loginUser,
			logoutUser,
			currentPage,
			setCurrentPage,
			version,
			setVersion
		}}>
			{children}
		</CommonContext.Provider>
	);
};
