import { useEffect, useMemo, useState, useContext, useRef } from "react";
import { ProviderSetting, MultiProviderConfig } from "@/api/provider";
import { get_provider_settings, get_multi_providers, update_multi_providers, add_provider_config, delete_provider_config } from "@/api/provider";
import { CommonContext } from "@/components/user";
import { useRouter } from "next/navigation";
import SearchableSelect from "./searchable-select";
import classNames from "classnames";
import { toast } from "react-toastify";
import { PlusIcon, TrashIcon, GearIcon } from "@/components/icons";

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
			const success = await add_provider_config(credentials, providerName, config);
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
						<label className={classNames("input input-bordered flex items-center gap-2", {
							'input-error': !formData.email
						})}>
							Email
							<input
								type="text"
								className="grow"
								placeholder="Input the email"
								value={formData.email}
								onChange={(e) => setFormData(prev => ({ ...prev, email: e.target.value }))}
							/>
						</label>
					)}

					{currentProviderSettings.password && (
						<label className={classNames("input input-bordered flex items-center gap-2", {
							'input-error': !formData.password
						})}>
							Password
							<input
								type="password"
								className="grow"
								placeholder="Input the password"
								value={formData.password}
								onChange={(e) => setFormData(prev => ({ ...prev, password: e.target.value }))}
							/>
						</label>
					)}

					{currentProviderSettings.login_token && (
						<label className={classNames("input input-bordered flex items-center gap-2", {
							'input-error': !formData.loginToken
						})}>
							Login Token
							<input
								type="text"
								className="grow"
								placeholder="Input the token"
								value={formData.loginToken}
								onChange={(e) => setFormData(prev => ({ ...prev, loginToken: e.target.value }))}
							/>
						</label>
					)}

					{currentProviderSettings.app_key && (
						<label className={classNames("input input-bordered flex items-center gap-2", {
							'input-error': !formData.appKey
						})}>
							App Key
							<input
								type="text"
								className="grow"
								placeholder="Input the app key"
								value={formData.appKey}
								onChange={(e) => setFormData(prev => ({ ...prev, appKey: e.target.value }))}
							/>
						</label>
					)}

					{currentProviderSettings.app_secret && (
						<label className={classNames("input input-bordered flex items-center gap-2", {
							'input-error': !formData.appSecret
						})}>
							App Secret
							<input
								type="text"
								className="grow"
								placeholder="Input the app secret"
								value={formData.appSecret}
								onChange={(e) => setFormData(prev => ({ ...prev, appSecret: e.target.value }))}
							/>
						</label>
					)}

					{currentProviderSettings.consumer_key && (
						<label className={classNames("input input-bordered flex items-center gap-2", {
							'input-error': !formData.consumerKey
						})}>
							Consumer Key
							<input
								type="text"
								className="grow"
								placeholder="Input the consumer key"
								value={formData.consumerKey}
								onChange={(e) => setFormData(prev => ({ ...prev, consumerKey: e.target.value }))}
							/>
						</label>
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
			<div className="flex items-center justify-between mb-4">
				<h2 className="text-xl font-semibold text-neutral-500">Provider Settings</h2>
				<button className="btn btn-primary btn-sm" onClick={openAddModal}>
					<PlusIcon />
					Add Provider
				</button>
			</div>

			{Object.keys(providers).length > 0 ? (
				<div className="flex flex-wrap gap-2">
					{Object.keys(providers).map((providerName, index) => (
						<div
							key={providerName}
							className={classNames("card w-full bg-primary-content shadow-xl mb-1", {
								"md:w-1/3": (index + 1) % 3 !== 0,
								"md:flex-1": (index + 1) % 3 === 0
							})}
						>
							<div className="card-body">
								<h2 className="card-title">
									{providerName}
								</h2>
								<div className="flex flex-wrap justify-start gap-2">
									{/* Show configured credential types as badges */}
									{providers[providerName]?.email && (
										<div className="badge badge-primary">Email</div>
									)}
									{providers[providerName]?.password && (
										<div className="badge badge-primary">Password</div>
									)}
									{providers[providerName]?.login_token && (
										<div className="badge badge-primary">Token</div>
									)}
									{providers[providerName]?.app_key && (
										<div className="badge badge-primary">App Key</div>
									)}
									{providers[providerName]?.app_secret && (
										<div className="badge badge-primary">App Secret</div>
									)}
									{providers[providerName]?.consumer_key && (
										<div className="badge badge-primary">Consumer Key</div>
									)}
								</div>
								<div className="card-actions justify-end">
									<button
										className="btn btn-secondary btn-sm"
										onClick={() => openEditModal(providerName)}
									>
										Edit
									</button>
									<button
										className="btn btn-error btn-sm"
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
				<div className="text-center py-8">
					<p className="text-neutral-500 mb-4">No providers configured yet.</p>
					<button className="btn btn-primary" onClick={openAddModal}>
						<PlusIcon />
						Add Your First Provider
					</button>
				</div>
			)}

			{/* Add Provider Modal */}
			<dialog id="add_provider_modal" className="modal" ref={addModalRef}>
				<div className="modal-box max-w-lg">
					<h3 className="font-bold text-lg mb-4">Add New Provider</h3>
					<div className="flex flex-col gap-4">
						{renderProviderForm()}
					</div>
					<div className="modal-action">
						<button className="btn" onClick={closeAddModal}>
							Cancel
						</button>
						<button
							className="btn btn-primary"
							onClick={handleSaveProvider}
							disabled={!selectedProvider || !currentProviderSettings}
						>
							Add Provider
						</button>
					</div>
				</div>
			</dialog>

			{/* Edit Provider Modal */}
			<dialog id="edit_provider_modal" className="modal" ref={editModalRef}>
				<div className="modal-box max-w-lg">
					<h3 className="font-bold text-lg mb-4">Edit {editingProvider}</h3>
					<div className="flex flex-col gap-4">
						{renderProviderForm()}
					</div>
					<div className="modal-action">
						<button className="btn" onClick={closeEditModal}>
							Cancel
						</button>
						<button className="btn btn-primary" onClick={handleSaveProvider}>
							Update Provider
						</button>
					</div>
				</div>
			</dialog>

			{/* Delete Confirmation Modal */}
			<dialog id="delete_provider_modal" className="modal" ref={deleteModalRef}>
				<div className="modal-box">
					<h3 className="font-bold text-lg">Remove Provider</h3>
					<p className="py-4">
						Are you sure you want to remove <strong>{providerToDelete}</strong> provider? 
						This action cannot be undone and will affect any domains using this provider.
					</p>
					<div className="modal-action">
						<button className="btn" onClick={closeDeleteModal}>
							Cancel
						</button>
						<button className="btn btn-error" onClick={handleDeleteProvider}>
							Remove Provider
						</button>
					</div>
				</div>
			</dialog>
		</div>
	);
};