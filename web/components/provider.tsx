import { useEffect, useMemo, useState } from "react";
import { Provider, ProviderSetting } from "@/api/provider";
import { get_provider_settings, get_provider, update_provider } from "@/api/provider";
import { useContext } from "react";
import { CommonContext } from "@/components/user";
import { useRouter } from "next/navigation";
import SearchableSelect from "./searchable-select";
import classNames from "classnames";
import { toast } from "react-toastify";

export const ProviderControl = () => {
	const { credentials } = useContext(CommonContext);
	const [currentProvider, setCurrentProvider] = useState<Provider>();
	const [providerSettings, setProviderSettings] = useState<ProviderSetting[]>([]);
	const router = useRouter();

	const [username, setUsername] = useState(currentProvider?.username || '');
	const [email, setEmail] = useState(currentProvider?.email || '');
	const [password, setPassword] = useState(currentProvider?.password || '');
	const [loginToken, setLoginToken] = useState(currentProvider?.login_token || '');
	const [appKey, setAppKey] = useState(currentProvider?.app_key || '');
	const [appSecret, setAppSecret] = useState(currentProvider?.app_secret || '');
	const [consumerKey, setConsumerKey] = useState(currentProvider?.consumer_key || '');

	const options = useMemo(() => {
		if (providerSettings) {
			return providerSettings.map((setting) => {
				return {
					value: setting.name,
					label: setting.name
				};
			});
		}
		return [];
	}, [providerSettings]);

	const currentProviderSettings = useMemo(() => {
		if (currentProvider) {
			const settings = providerSettings.filter((setting) => setting.name === currentProvider.provider);
			if (settings.length > 0) {
				return settings[0];
			}
		}
		return null;
	}, [currentProvider, providerSettings]);

	const onProviderSelected = (provider: string) => {
		if (currentProvider) {
			// update provider
			setCurrentProvider({
				...currentProvider,
				provider
			});
		}
	}

	const onSaveProviderSettings = () => {
		if (currentProvider) {
			const newProvider: Provider = {
				provider: currentProvider.provider,
				username: '',
				email: '',
				password: '',
				login_token: '',
				app_key: '',
				app_secret: '',
				consumer_key: ''
			};

			// update values according to the current provider settings
			if (currentProviderSettings?.username) {
				if (username) {
					newProvider.username = username;
				} else {
					toast.error('Username is required');
					return;
				}
			}
			if (currentProviderSettings?.email) {
				if (email) {
					newProvider.email = email;
				} else {
					toast.error('Email is required');
					return;
				}
			}
			if (currentProviderSettings?.password) {
				if (password) {
					newProvider.password = password;
				} else {
					toast.error('Password is required');
					return;
				}
			}
			if (currentProviderSettings?.login_token) {
				if (loginToken) {
					newProvider.login_token = loginToken;
				} else {
					toast.error('Login Token is required');
					return;
				}
			}
			if (currentProviderSettings?.app_key) {
				if (appKey) {
					newProvider.app_key = appKey;
				} else {
					toast.error('App Key is required');
					return;
				}
			}
			if (currentProviderSettings?.app_secret) {
				if (appSecret) {
					newProvider.app_secret = appSecret;
				} else {
					toast.error('App Secret is required');
					return;
				}
			}
			if (currentProviderSettings?.consumer_key) {
				if (consumerKey) {
					newProvider.consumer_key = consumerKey;
				} else {
					toast.error('Consumer Key is required');
					return;
				}
			}

			// save provider
			if (credentials) {
				update_provider(credentials, newProvider).then((success) => {
					if (success) {
						toast.success('Provider settings saved successfully');
						// update the current provider
						setCurrentProvider(newProvider);
					} else {
						toast.error('Failed to save provider settings');
					}
				});
			} else {
				toast.error('Invalid credentials');
			}
		}
	}

	useEffect(() => {
		if (!credentials) {
			router.push('/login');
			return;
		}
		// fetch provider settings
		get_provider_settings(credentials).then((settings) => {
			if (settings.length) {
				setProviderSettings(settings);
			}
		});
		// fetch provider
		get_provider(credentials).then((provider) => {
			setCurrentProvider(provider);
			// set the state according to the current provider settings
			if (provider?.username) {
				setUsername(provider.username);
			}
			if (provider?.email) {
				setEmail(provider.email);
			}
			if (provider?.password) {
				setPassword(provider.password);
			}
			if (provider?.login_token) {
				setLoginToken(provider.login_token);
			}
			if (provider?.app_key) {
				setAppKey(provider.app_key);
			}
			if (provider?.app_secret) {
				setAppSecret(provider.app_secret);
			}
			if (provider?.consumer_key) {
				setConsumerKey(provider.consumer_key);
			}
		});
	}, [credentials, router]);

	return (
		<div className="card w-96 sm:w-2/3 bg-primary-content shadow-xl">
			<div className="card-body">
				<h2 className="card-title">Provider Settings</h2>
				<SearchableSelect
					options={options}
					placeholder="Select Provider"
					defaultValue={currentProvider?.provider}
					onSelected={onProviderSelected}
				/>
				{
					currentProviderSettings && currentProviderSettings.username ? (
						<label className={classNames("input input-bordered flex items-center gap-2", {
							'input-error': !username
						})}>
							Username
							<input type="text" className="grow" placeholder="Input the username" value={username}
								onChange={
									(e) => {
										setUsername(e.target.value);
									}
								}
							/>
						</label>
					) : null
				}
				{
					currentProviderSettings && currentProviderSettings.email ? (
						<label className={classNames("input input-bordered flex items-center gap-2", {
							'input-error': !email
						})}>
							Email
							<input type="text" className="grow" placeholder="Input the email" value={email}
								onChange={
									(e) => {
										setEmail(e.target.value);
									}
								}
							/>
						</label>
					) : null
				}
				{
					currentProviderSettings && currentProviderSettings.password ? (
						<label className={classNames("input input-bordered flex items-center gap-2", {
							'input-error': !password
						})}>
							Password
							<input type="password" className="grow" placeholder="Input the password" value={password}
								onChange={
									(e) => {
										setPassword(e.target.value);
									}
								}
							/>
						</label>
					) : null
				}
				{
					currentProviderSettings && currentProviderSettings.login_token ? (
						<label className={classNames("input input-bordered flex items-center gap-2", {
							'input-error': !loginToken
						})}>
							Login Token
							<input type="text" className="grow" placeholder="Input the token" value={loginToken}
								onChange={
									(e) => {
										setLoginToken(e.target.value);
									}
								}
							/>
						</label>
					) : null
				}
				{
					currentProviderSettings && currentProviderSettings.app_key ? (
						<label className={classNames("input input-bordered flex items-center gap-2", {
							'input-error': !appKey
						})}>
							App Key
							<input type="text" className="grow" placeholder="Input the app key" value={appKey}
								onChange={
									(e) => {
										setAppKey(e.target.value);
									}
								}
							/>
						</label>
					) : null
				}
				{
					currentProviderSettings && currentProviderSettings.app_secret ? (
						<label className={classNames("input input-bordered flex items-center gap-2", {
							'input-error': !appSecret
						})}>
							App Secret
							<input type="text" className="grow" placeholder="Input the app secret" value={appSecret}
								onChange={
									(e) => {
										setAppSecret(e.target.value);
									}
								}
							/>
						</label>
					) : null
				}
				{
					currentProviderSettings && currentProviderSettings.consumer_key ? (
						<label className={classNames("input input-bordered flex items-center gap-2", {
							'input-error': !consumerKey
						})}>
							Consumer Key
							<input type="text" className="grow" placeholder="Input the consumer key" value={consumerKey}
								onChange={
									(e) => {
										setConsumerKey(e.target.value);
									}
								} />
						</label>
					) : null
				}
				<div className="card-actions justify-end">
					<button className="btn btn-primary" onClick={onSaveProviderSettings}>Save</button>
				</div>
			</div>
		</div >
	);
}