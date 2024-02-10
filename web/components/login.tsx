// components/Login.js
'use client';
import React from 'react';
import { Button, Input, Card, CardHeader, CardBody, Divider, Spacer } from '@nextui-org/react';
import { EyeFilledIcon } from "@/components/EyeFilledIcon";
import { EyeSlashFilledIcon } from "@/components/EyeSlashFilledIcon";

interface LoginProps {
	onLogin: (username: string, password: string) => void;
}

export const Login = (props: LoginProps) => {
	const [username, setUsername] = React.useState('');
	const [password, setPassword] = React.useState('');
	const [isVisible, setIsVisible] = React.useState(false);

	const handleLogin = (e: React.FormEvent) => {
		e.preventDefault();
		props.onLogin(username, password);
	};

	const toggleVisibility = () => setIsVisible(!isVisible);

	return (
		<div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '80vh' }}>
			<Card className=''>
				<CardHeader>
					<h2 className='text-xl font-bold text-center'>Login</h2>
				</CardHeader>
				<Divider />
				<CardBody>
					<form onSubmit={handleLogin}>
						<Input
							label="Username"
							isClearable={true}
							variant="bordered"
							fullWidth
							color="primary"
							size="lg"
							placeholder="Enter your username"
							value={username}
							onChange={(e) => setUsername(e.target.value)}
						/>
						<Spacer y={2} />
						<Input
							label="Password"
							variant="bordered"
							placeholder="Enter your password"
							color="primary"
							endContent={
								<button className="focus:outline-none" type="button" onClick={toggleVisibility}>
									{isVisible ? (
										<EyeSlashFilledIcon className="text-2xl text-default-400 pointer-events-none" />
									) : (
										<EyeFilledIcon className="text-2xl text-default-400 pointer-events-none" />
									)}
								</button>
							}
							type={isVisible ? "text" : "password"}
							className="max-w-xs"
							onChange={(e) => setPassword(e.target.value)}
						/>
						<Spacer y={2} />
						<Button type="submit" fullWidth size="lg" color="primary">Login</Button>
					</form>
				</CardBody>
			</Card>
		</div>
	);
};

