'use client';
// pages/login.js or another page
import React from 'react';
import { Login } from '@/components/login';

const LoginPage = () => {
	const handleLogin = (username: String, password: String) => {
		console.log('Login attempt with:', username, password);
		// Here you would typically handle the login logic,
		// such as setting user state or redirecting upon successful login.
	};

	return <Login onLogin={handleLogin} />;
};

export default LoginPage;
