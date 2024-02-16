'use client'
import { createContext, useState, ReactNode } from 'react';

type UserAction = (_: string) => void;
type UserLogoutAction = () => void;
type PageAction = (_: string) => void;

interface ICommonContext {
	credentials: string | null;
	loginUser: UserAction;
	logoutUser: UserLogoutAction;
	saveVersion: UserAction;
	currentPage: string;
	setCurrentPage: PageAction;
	version: string | null;
}

export const CommonContext = createContext<ICommonContext>({
	credentials: null,
	loginUser: (_: string) => { },
	logoutUser: () => { },
	saveVersion: (_: string) => { },
	currentPage: '',
	setCurrentPage: (_: string) => { },
	version: null,
});

interface UserProviderProps {
	children: ReactNode;
}

// user provider
export const UserProvider = ({ children }: UserProviderProps) => {
	const [credentials, setCredentials] = useState<string | null>(typeof window !== "undefined" ? localStorage.getItem('credentials') : null);
	const [currentPage, setCurrentPage] = useState<string>('');
	const [version, setVersion] = useState<string | null>(typeof window !== "undefined" ? localStorage.getItem('version') : '');

	const loginUser = (credentials: string) => {
		setCredentials(credentials);
		localStorage.setItem('credentials', credentials);
	};

	const logoutUser = () => {
		setCredentials(null);
		localStorage.removeItem('credentials');
	};

	const saveVersion = (version: string) => {
		setVersion(version);
		localStorage.setItem('version', version);
	}

	return (
		<CommonContext.Provider value={{
			credentials,
			loginUser,
			logoutUser,
			saveVersion,
			currentPage,
			setCurrentPage,
			version,
		}}>
			{children}
		</CommonContext.Provider>
	);
};
