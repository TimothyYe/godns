import { useEffect, useMemo, useState, useContext, useRef } from "react";
import { ProviderSetting, MultiProviderConfig } from "@/api/provider";
import { get_provider_settings, get_multi_providers, update_multi_providers, add_provider_config, delete_provider_config } from "@/api/provider";
import { CommonContext } from "@/components/user";
import { useRouter } from "next/navigation";
import SearchableSelect from "./searchable-select";
import classNames from "classnames";
import { toast } from "react-toastify";
import { PlusIcon, TrashIcon } from "@/components/icons";

interface ProviderConfigForm {
	email: string;
	password: string;
	loginToken: string;
	appKey: string;
	appSecret: string;
	consumerKey: string;
}

export const MultiProviderControl = () => {
	const { credentials } = useContext(CommonContext);
	const [providers, setProviders] = useState<MultiProviderConfig>({});
	const [providerSettings, setProviderSettings] = useState<ProviderSetting[]>([]);
	const [selectedProvider, setSelectedProvider] = useState<string>('');
	const [formData, setFormData] = useState<ProviderConfigForm>({
		email: '',
		password: '',
		loginToken: '',
		appKey: '',
		appSecret: '',
		consumerKey: ''
	});
	const [editingProvider, setEditingProvider] = useState<string | null>(null);
	const addModalRef = useRef<HTMLDialogElement | null>(null);
	const editModalRef = useRef<HTMLDialogElement | null>(null);
	const deleteModalRef = useRef<HTMLDialogElement | null>(null);
	const [providerToDelete, setProviderToDelete] = useState<string>('');
	const router = useRouter();

	const availableProviders = useMemo(() => {
		if (providerSettings) {
			return providerSettings
				.filter(setting => !providers[setting.name])
				.map((setting) => ({
					value: setting.name,
					label: setting.name
				}));
		}
		return [];
	}, [providerSettings, providers]);

	const currentProviderSettings = useMemo(() => {
		const providerName = editingProvider || selectedProvider;
		if (providerName) {
			const settings = providerSettings.filter((setting) => setting.name === providerName);
			if (settings.length > 0) {
				return settings[0];
			}
		}
		return null;
	}, [editingProvider, selectedProvider, providerSettings]);

	const resetForm = () => {
		setFormData({
			email: '',
			password: '',
			loginToken: '',
			appKey: '',
			appSecret: '',
			consumerKey: ''
		});
		setSelectedProvider('');
		setEditingProvider(null);
	};

	const openAddModal = () => {
		resetForm();
		addModalRef.current?.showModal();
	};

	const closeAddModal = () => {
		addModalRef.current?.close();
		resetForm();
	};

	const openEditModal = (providerName: string) => {
		setEditingProvider(providerName);
		loadProviderConfig(providerName);
		editModalRef.current?.showModal();
	};

	const closeEditModal = () => {
		editModalRef.current?.close();
		resetForm();
	};

	const openDeleteModal = (providerName: string) => {
		setProviderToDelete(providerName);
		deleteModalRef.current?.showModal();
	};

	const closeDeleteModal = () => {
		deleteModalRef.current?.close();
		setProviderToDelete('');
	};

	const loadProviderConfig = (providerName: string) => {
		const config = providers[providerName];
		if (config) {
			setFormData({
				email: config.email || '',
				password: config.password || '',
				loginToken: config.login_token || '',
				appKey: config.app_key || '',
				appSecret: config.app_secret || '',
				consumerKey: config.consumer_key || ''
			});
		}
	};

	const handleDeleteProvider = async () => {
		if (credentials && providerToDelete) {
			const success = await delete_provider_config(credentials, providerToDelete);
			if (success) {
				const newProviders = { ...providers };
				delete newProviders[providerToDelete];
				setProviders(newProviders);
				toast.success(`${providerToDelete} provider removed successfully`);
				closeDeleteModal();
			} else {
				toast.error(`Failed to remove ${providerToDelete} provider`);
			}
		}
	};

	const validateForm = (): boolean => {
		if (!currentProviderSettings) return false;

		const errors: string[] = [];

		if (currentProviderSettings.email && !formData.email) {
			errors.push('Email is required');
		}
		if (currentProviderSettings.password && !formData.password) {
			errors.push('Password is required');
		}
		if (currentProviderSettings.login_token && !formData.loginToken) {
			errors.push('Login Token is required');
		}
		if (currentProviderSettings.app_key && !formData.appKey) {
			errors.push('App Key is required');
		}
		if (currentProviderSettings.app_secret && !formData.appSecret) {
			errors.push('App Secret is required');
		}
		if (currentProviderSettings.consumer_key && !formData.consumerKey) {
			errors.push('Consumer Key is required');
		}

		if (errors.length > 0) {
			toast.error(errors[0]);
			return false;
		}

		return true;
	};

	const handleSaveProvider = async () => {
		if (!validateForm()) return;

		const providerName = editingProvider || selectedProvider;
		if (!providerName) return;

		const config: any = {};
		if (currentProviderSettings?.email) config.email = formData.email;
		if (currentProviderSettings?.password) config.password = formData.password;
		if (currentProviderSettings?.login_token) config.login_token = formData.loginToken;
		if (currentProviderSettings?.app_key) config.app_key = formData.appKey;
		if (currentProviderSettings?.app_secret) config.app_secret = formData.appSecret;
		if (currentProviderSettings?.consumer_key) config.consumer_key = formData.consumerKey;

		if (credentials) {
			let success = false;

			if (editingProvider) {
				// For editing, update the entire providers map
				const updatedProviders = {
					...providers,
					[providerName]: config
				};
				success = await update_multi_providers(credentials, updatedProviders);
			} else {
				// For adding, use add_provider_config
				success = await add_provider_config(credentials, providerName, config);
			}

			if (success) {
				setProviders(prev => ({
					...prev,
					[providerName]: config
				}));
				toast.success(`${providerName} provider ${editingProvider ? 'updated' : 'added'} successfully`);
				if (editingProvider) {
					closeEditModal();
				} else {
					closeAddModal();
				}
			} else {
				toast.error(`Failed to ${editingProvider ? 'update' : 'add'} ${providerName} provider`);
			}
		}
	};

	const renderTextField = (
		label: string,
		value: string,
		onChange: (next: string) => void,
		placeholder: string,
		type: string = 'text',
		hint?: string
	) => (
		<label className="theme-field">
			<span className="theme-field-label">{label}</span>
			<input
				type={type}
				className="input theme-input w-full rounded-2xl"
				placeholder={placeholder}
				value={value}
				onChange={(e) => onChange(e.target.value)}
			/>
			{hint ? <span className="theme-field-hint">{hint}</span> : null}
		</label>
	);

	const getAccentChipClass = (index: number) => index % 2 === 0 ? 'theme-chip-sky' : 'theme-chip-violet';

	const renderProviderForm = () => (
		<>
			{!editingProvider && (
				<SearchableSelect
					options={availableProviders}
					placeholder="Select Provider to Add"
					defaultValue=""
					onSelected={setSelectedProvider}
				/>
			)}

			{currentProviderSettings && (
				<>
					{currentProviderSettings.email && (
						renderTextField('Email', formData.email, (email) => setFormData(prev => ({ ...prev, email })), 'name@example.com')
					)}

					{currentProviderSettings.password && (
						renderTextField('Password', formData.password, (password) => setFormData(prev => ({ ...prev, password })), 'Provider password', 'password')
					)}

					{currentProviderSettings.login_token && (
						renderTextField('Login Token', formData.loginToken, (loginToken) => setFormData(prev => ({ ...prev, loginToken })), 'Paste the login token')
					)}

					{currentProviderSettings.app_key && (
						renderTextField('App Key', formData.appKey, (appKey) => setFormData(prev => ({ ...prev, appKey })), 'Paste the app key')
					)}

					{currentProviderSettings.app_secret && (
						renderTextField('App Secret', formData.appSecret, (appSecret) => setFormData(prev => ({ ...prev, appSecret })), 'Paste the app secret')
					)}

					{currentProviderSettings.consumer_key && (
						renderTextField('Consumer Key', formData.consumerKey, (consumerKey) => setFormData(prev => ({ ...prev, consumerKey })), 'Paste the consumer key')
					)}
				</>
			)}
		</>
	);

	useEffect(() => {
		if (!credentials) {
			router.push('/login');
			return;
		}

		// Fetch provider settings
		get_provider_settings(credentials).then((settings) => {
			if (settings.length) {
				setProviderSettings(settings);
			}
		});

		// Fetch current multi-provider configuration
		get_multi_providers(credentials).then((multiProviders) => {
			setProviders(multiProviders);
		});
	}, [credentials, router]);

	return (
		<div className="w-full">
			<div className="section-header">
				<div className="section-copy">
					<div className="section-label ml-0">Providers</div>
					<h2 className="text-2xl font-semibold tracking-tight theme-heading">Provider credentials</h2>
					<p className="mt-2 text-sm leading-7 theme-muted">
						Keep provider credentials organized and easy to scan before attaching domains to them.
					</p>
				</div>
				<div className="flex items-center gap-3">
					{Object.keys(providers).length > 0 ? (
						<div className="theme-chip-violet rounded-full px-3 py-1 text-xs font-medium">
							{Object.keys(providers).length} providers
						</div>
					) : null}
					<button className="theme-primary-violet btn btn-sm rounded-xl border-none px-4 shadow-lg shadow-violet-950/20" onClick={openAddModal}>
						<PlusIcon />
						Add Provider
					</button>
				</div>
			</div>

			{Object.keys(providers).length > 0 ? (
				<div className="mt-5 grid gap-4 lg:grid-cols-3">
					{Object.keys(providers).map((providerName) => (
						<div
							key={providerName}
							className="surface-card"
						>
							<div className="card-body gap-4 p-5 sm:p-6">
								<div className="metric-kicker">Provider</div>
								<div className="flex items-start justify-between gap-4">
									<h2 className="text-xl font-semibold tracking-tight theme-heading">
										{providerName}
									</h2>
									<div className="theme-chip-violet rounded-full px-3 py-1 text-xs font-medium">
										{Object.keys(providers[providerName] || {}).length} fields
									</div>
								</div>
								<div className="flex flex-wrap justify-start gap-2">
									{providers[providerName]?.email && (
										<div className={classNames("badge rounded-full px-3 py-3", getAccentChipClass(0))}>Email</div>
									)}
									{providers[providerName]?.password && (
										<div className={classNames("badge rounded-full px-3 py-3", getAccentChipClass(1))}>Password</div>
									)}
									{providers[providerName]?.login_token && (
										<div className={classNames("badge rounded-full px-3 py-3", getAccentChipClass(2))}>Token</div>
									)}
									{providers[providerName]?.app_key && (
										<div className={classNames("badge rounded-full px-3 py-3", getAccentChipClass(3))}>App Key</div>
									)}
									{providers[providerName]?.app_secret && (
										<div className={classNames("badge rounded-full px-3 py-3", getAccentChipClass(4))}>App Secret</div>
									)}
									{providers[providerName]?.consumer_key && (
										<div className={classNames("badge rounded-full px-3 py-3", getAccentChipClass(5))}>Consumer Key</div>
									)}
								</div>
								<p className="text-sm leading-6 theme-muted">
									These credentials are available to any domain mapped to this provider.
								</p>
								<div className="card-actions justify-end">
									<button
										className="theme-subtle-btn btn btn-sm rounded-xl px-4"
										onClick={() => openEditModal(providerName)}
									>
										Edit
									</button>
									<button
										className="theme-danger btn btn-sm rounded-xl border-none px-3"
										onClick={() => openDeleteModal(providerName)}
									>
										<TrashIcon />
									</button>
								</div>
							</div>
						</div>
					))}
				</div>
			) : (
				<div className="surface-panel-soft px-6 py-10 text-center">
					<p className="mb-3 text-lg font-medium theme-heading">No providers configured yet.</p>
					<p className="mx-auto mb-5 max-w-lg text-sm leading-7 theme-muted">
						Add your first provider before creating domains. Credentials added here are reused by the domains below.
					</p>
					<button className="theme-primary-violet btn rounded-xl border-none" onClick={openAddModal}>
						<PlusIcon />
						Add Your First Provider
					</button>
				</div>
			)}

			{/* Add Provider Modal */}
			<dialog id="add_provider_modal" className="modal modal-bottom sm:modal-middle" ref={addModalRef}>
				<div className="theme-modal modal-box max-w-lg rounded-[1.5rem]">
					<h3 className="mb-4 text-xl font-semibold tracking-tight theme-heading">Add New Provider</h3>
					<p className="mb-4 text-sm leading-7 theme-muted">
						Choose a provider first, then fill in only the credentials required by that integration.
					</p>
					<div className="flex flex-col gap-4">
						{renderProviderForm()}
					</div>
					<div className="modal-action">
						<button className="theme-subtle-btn btn rounded-xl" onClick={closeAddModal}>
							Cancel
						</button>
						<button
							className="theme-primary-violet btn rounded-xl border-none"
							onClick={handleSaveProvider}
							disabled={!selectedProvider || !currentProviderSettings}
						>
							Add Provider
						</button>
					</div>
				</div>
				<form method="dialog" className="modal-backdrop">
					<button aria-label="Close add provider dialog">close</button>
				</form>
			</dialog>

			{/* Edit Provider Modal */}
			<dialog id="edit_provider_modal" className="modal modal-bottom sm:modal-middle" ref={editModalRef}>
				<div className="theme-modal modal-box max-w-lg rounded-[1.5rem]">
					<h3 className="mb-4 text-xl font-semibold tracking-tight theme-heading">Edit {editingProvider}</h3>
					<p className="mb-4 text-sm leading-7 theme-muted">
						Update the saved credentials for this provider profile.
					</p>
					<div className="flex flex-col gap-4">
						{renderProviderForm()}
					</div>
					<div className="modal-action">
						<button className="theme-subtle-btn btn rounded-xl" onClick={closeEditModal}>
							Cancel
						</button>
						<button className="theme-primary-violet btn rounded-xl border-none" onClick={handleSaveProvider}>
							Update Provider
						</button>
					</div>
				</div>
				<form method="dialog" className="modal-backdrop">
					<button aria-label="Close edit provider dialog">close</button>
				</form>
			</dialog>

			{/* Delete Confirmation Modal */}
			<dialog id="delete_provider_modal" className="modal modal-bottom sm:modal-middle" ref={deleteModalRef}>
				<div className="theme-modal modal-box rounded-[1.5rem]">
					<h3 className="text-lg font-semibold theme-heading">Remove Provider</h3>
					<p className="py-4 text-sm leading-7 theme-muted">
						Are you sure you want to remove <strong>{providerToDelete}</strong> provider? 
						This action cannot be undone and will affect any domains using this provider.
					</p>
					<div className="modal-action">
						<button className="theme-subtle-btn btn rounded-xl" onClick={closeDeleteModal}>
							Cancel
						</button>
						<button className="theme-danger btn rounded-xl border-none" onClick={handleDeleteProvider}>
							Remove Provider
						</button>
					</div>
				</div>
				<form method="dialog" className="modal-backdrop">
					<button aria-label="Close delete provider dialog">close</button>
				</form>
			</dialog>
		</div>
	);
};
