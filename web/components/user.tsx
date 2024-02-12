'use client';
import { createContext, useState, useEffect, ReactNode } from 'react';

type UserAction = (_: string) => void;
type PageAction = (_: string) => void;

interface ICommonContext {
	credentials: string | null;
	loginUser: UserAction;
	logoutUser: UserAction;
	currentPage: string;
	setCurrentPage: PageAction;
	version: string;
	setVersion: PageAction;
}

export const CommonContext = createContext<ICommonContext>({
	credentials: null,
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
	const [credentials, setCredentials] = useState<string | null>(localStorage.getItem('credentials'));
	const [currentPage, setCurrentPage] = useState<string>('Home');
	const [version, setVersion] = useState<string>('');

	const loginUser = (credentials: string) => {
		setCredentials(credentials);
		localStorage.setItem('credentials', credentials);
	};

	const logoutUser = () => {
		setCredentials(null);
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
