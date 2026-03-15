import { useEffect, useMemo, useState, useContext, useRef } from "react";
import { ProviderSetting, MultiProviderConfig } from "@/api/provider";
import { get_provider_settings, get_multi_providers, update_multi_providers, add_provider_config, delete_provider_config } from "@/api/provider";
import { CommonContext } from "@/components/user";
import { useRouter } from "next/navigation";
import SearchableSelect from "./searchable-select";
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

const fieldMeta = {
	email: {
		label: 'Email',
		placeholder: 'Enter the provider account email',
		hint: 'Used by providers that authenticate with an email address.',
		type: 'text',
	},
	password: {
		label: 'Password',
		placeholder: 'Enter the provider password',
		hint: 'Stored so GoDNS can authenticate against the provider API.',
		type: 'password',
	},
	loginToken: {
		label: 'Login token',
		placeholder: 'Enter the login token',
		hint: 'Some providers use an API or login token instead of a password.',
		type: 'text',
	},
	appKey: {
		label: 'App key',
		placeholder: 'Enter the app key',
		hint: 'Application identifier required by the provider API.',
		type: 'text',
	},
	appSecret: {
		label: 'App secret',
		placeholder: 'Enter the app secret',
		hint: 'Secret paired with the app key.',
		type: 'text',
	},
	consumerKey: {
		label: 'Consumer key',
		placeholder: 'Enter the consumer key',
		hint: 'End-user or account-level access token.',
		type: 'text',
	},
} as const;

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
	const [loading, setLoading] = useState(true);
	const addModalRef = useRef<HTMLDialogElement | null>(null);
	const editModalRef = useRef<HTMLDialogElement | null>(null);
	const deleteModalRef = useRef<HTMLDialogElement | null>(null);
	const [providerToDelete, setProviderToDelete] = useState<string>('');
	const router = useRouter();

	const availableProviders = useMemo(() => {
		return providerSettings
			.filter((setting) => !providers[setting.name])
			.map((setting) => ({
				value: setting.name,
				label: setting.name
			}));
	}, [providerSettings, providers]);

	const currentProviderSettings = useMemo(() => {
		const providerName = editingProvider || selectedProvider;
		if (!providerName) {
			return null;
		}

		return providerSettings.find((setting) => setting.name === providerName) || null;
	}, [editingProvider, selectedProvider, providerSettings]);

	const requiredFields = useMemo(() => {
		if (!currentProviderSettings) {
			return [];
		}

		return [
			currentProviderSettings.email ? 'email' : null,
			currentProviderSettings.password ? 'password' : null,
			currentProviderSettings.login_token ? 'loginToken' : null,
			currentProviderSettings.app_key ? 'appKey' : null,
			currentProviderSettings.app_secret ? 'appSecret' : null,
			currentProviderSettings.consumer_key ? 'consumerKey' : null,
		].filter(Boolean) as Array<keyof ProviderConfigForm>;
	}, [currentProviderSettings]);

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
		if (!config) {
			return;
		}

		setFormData({
			email: config.email || '',
			password: config.password || '',
			loginToken: config.login_token || '',
			appKey: config.app_key || '',
			appSecret: config.app_secret || '',
			consumerKey: config.consumer_key || ''
		});
	};

	const handleDeleteProvider = async () => {
		if (!credentials || !providerToDelete) {
			return;
		}

		const success = await delete_provider_config(credentials, providerToDelete);
		if (success) {
			setProviders((current) => {
				const nextProviders = { ...current };
				delete nextProviders[providerToDelete];
				return nextProviders;
			});
			toast.success(`${providerToDelete} provider removed successfully`);
			closeDeleteModal();
		} else {
			toast.error(`Failed to remove ${providerToDelete} provider`);
		}
	};

	const validateForm = (): boolean => {
		if (!currentProviderSettings) {
			return false;
		}

		for (const field of requiredFields) {
			if (!formData[field]) {
				toast.error(`${fieldMeta[field].label} is required`);
				return false;
			}
		}

		return true;
	};

	const handleSaveProvider = async () => {
		if (!validateForm()) {
			return;
		}

		const providerName = editingProvider || selectedProvider;
		if (!providerName || !credentials) {
			return;
		}

		const config: MultiProviderConfig[string] = {};
		if (currentProviderSettings?.email) config.email = formData.email;
		if (currentProviderSettings?.password) config.password = formData.password;
		if (currentProviderSettings?.login_token) config.login_token = formData.loginToken;
		if (currentProviderSettings?.app_key) config.app_key = formData.appKey;
		if (currentProviderSettings?.app_secret) config.app_secret = formData.appSecret;
		if (currentProviderSettings?.consumer_key) config.consumer_key = formData.consumerKey;

		let success = false;
		if (editingProvider) {
			success = await update_multi_providers(credentials, {
				...providers,
				[providerName]: config
			});
		} else {
			success = await add_provider_config(credentials, providerName, config);
		}

		if (!success) {
			toast.error(`Failed to ${editingProvider ? 'update' : 'add'} ${providerName} provider`);
			return;
		}

		setProviders((current) => ({
			...current,
			[providerName]: config
		}));
		toast.success(`${providerName} provider ${editingProvider ? 'updated' : 'added'} successfully`);
		if (editingProvider) {
			closeEditModal();
		} else {
			closeAddModal();
		}
	};

	const renderProviderForm = () => (
		<div className="grid gap-4">
			{!editingProvider ? (
				<label className="field-stack">
					<span className="field-label">Provider name</span>
					<SearchableSelect
						options={availableProviders}
						placeholder="Search and select a provider"
						defaultValue=""
						onSelected={setSelectedProvider}
					/>
					<span className="field-hint">Only providers not yet configured are shown here.</span>
				</label>
			) : null}

			{requiredFields.length > 0 ? requiredFields.map((field) => (
				<label className="field-stack" key={field}>
					<span className="field-label">{fieldMeta[field].label}</span>
					<input
						type={fieldMeta[field].type}
						className="input input-bordered h-12 w-full rounded-2xl"
						placeholder={fieldMeta[field].placeholder}
						value={formData[field]}
						onChange={(e) => setFormData((prev) => ({ ...prev, [field]: e.target.value }))}
					/>
					<span className="field-hint">{fieldMeta[field].hint}</span>
				</label>
			)) : currentProviderSettings ? (
				<div className="panel-muted">
					<p className="font-semibold">No additional credentials required</p>
					<p className="mt-2 text-sm text-base-content/60">
						This provider profile can be added without extra secrets in the current GoDNS configuration schema.
					</p>
				</div>
			) : null}
		</div>
	);

	useEffect(() => {
		if (!credentials) {
			router.push('/login');
			return;
		}

		setLoading(true);
		Promise.all([
			get_provider_settings(credentials),
			get_multi_providers(credentials),
		]).then(([settings, multiProviders]) => {
			setProviderSettings(settings || []);
			setProviders(multiProviders || {});
		}).finally(() => {
			setLoading(false);
		});
	}, [credentials, router]);

	return (
		<section className="panel">
			<div className="section-header">
				<div className="space-y-1">
					<h2 className="section-title">Provider profiles</h2>
					<p className="section-subtitle">
						Define one credential profile per DNS provider so domains can be mapped explicitly and audited later.
					</p>
				</div>
				<button className="btn btn-primary rounded-full px-5" onClick={openAddModal}>
					<PlusIcon />
					Add provider
				</button>
			</div>

			<div className="mb-5 flex flex-wrap items-center gap-3 text-sm text-base-content/60">
				<span className="badge badge-ghost badge-lg">{Object.keys(providers).length} configured</span>
				<span className="badge badge-ghost badge-lg">{availableProviders.length} available to add</span>
			</div>

			{loading ? (
				<div className="flex min-h-48 items-center justify-center">
					<span className="loading loading-spinner loading-lg" />
				</div>
			) : Object.keys(providers).length > 0 ? (
				<div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
					{Object.keys(providers).map((providerName) => {
						const configuredFields = [
							providers[providerName]?.email ? 'Email' : null,
							providers[providerName]?.password ? 'Password' : null,
							providers[providerName]?.login_token ? 'Token' : null,
							providers[providerName]?.app_key ? 'App key' : null,
							providers[providerName]?.app_secret ? 'App secret' : null,
							providers[providerName]?.consumer_key ? 'Consumer key' : null,
						].filter(Boolean);

						return (
							<article key={providerName} className="flex flex-col rounded-[1.5rem] border border-base-300/70 bg-gradient-to-br from-base-100 to-base-200/60 p-5 shadow-sm">
								<div className="flex items-start justify-between gap-4">
									<div className="space-y-2">
										<h3 className="text-lg font-semibold tracking-tight">{providerName}</h3>
										<div className="badge badge-ghost badge-lg">{configuredFields.length} credentials configured</div>
									</div>
									<button className="btn btn-ghost btn-sm rounded-full text-error hover:bg-error/10" onClick={() => openDeleteModal(providerName)}>
										<TrashIcon />
									</button>
								</div>

								<div className="chip-list mt-5 flex-1">
									{configuredFields.length > 0 ? configuredFields.map((field) => (
										<div key={field} className="badge badge-primary badge-outline h-auto px-3 py-3">
											{field}
										</div>
									)) : (
										<p className="text-sm text-base-content/55">No secrets stored for this profile.</p>
									)}
								</div>

								<div className="mt-5 flex justify-end">
									<button className="btn btn-secondary rounded-full px-5" onClick={() => openEditModal(providerName)}>
										Edit profile
									</button>
								</div>
							</article>
						);
					})}
				</div>
			) : (
				<div className="panel-muted">
					<p className="text-lg font-semibold">No providers configured yet</p>
					<p className="mt-2 text-sm text-base-content/60">
						Add at least one provider profile before creating domains. This page supports multi-provider GoDNS setups, so keep profiles focused and named by provider.
					</p>
				</div>
			)}

			<dialog className="modal" ref={addModalRef}>
				<div className="modal-box max-w-2xl rounded-[1.75rem]">
					<h3 className="text-xl font-semibold tracking-tight">Add provider profile</h3>
					<p className="pb-5 pt-2 text-sm text-base-content/65">
						Only the credentials required by the selected provider will be shown here.
					</p>
					{renderProviderForm()}
					<div className="modal-action">
						<button className="btn rounded-full" onClick={closeAddModal}>
							Cancel
						</button>
						<button
							className="btn btn-primary rounded-full px-5"
							onClick={handleSaveProvider}
							disabled={!selectedProvider || !currentProviderSettings}
						>
							Add provider
						</button>
					</div>
				</div>
			</dialog>

			<dialog className="modal" ref={editModalRef}>
				<div className="modal-box max-w-2xl rounded-[1.75rem]">
					<h3 className="text-xl font-semibold tracking-tight">Edit {editingProvider}</h3>
					<p className="pb-5 pt-2 text-sm text-base-content/65">
						Update only the credentials GoDNS needs for this provider profile.
					</p>
					{renderProviderForm()}
					<div className="modal-action">
						<button className="btn rounded-full" onClick={closeEditModal}>
							Cancel
						</button>
						<button className="btn btn-primary rounded-full px-5" onClick={handleSaveProvider}>
							Update provider
						</button>
					</div>
				</div>
			</dialog>

			<dialog className="modal" ref={deleteModalRef}>
				<div className="modal-box rounded-[1.75rem]">
					<h3 className="text-lg font-semibold">Delete provider profile</h3>
					<p className="py-3 text-sm text-base-content/65">
						Remove <strong>{providerToDelete}</strong> from the configuration. Any domains assigned to this provider will need to be re-mapped before GoDNS can manage them safely.
					</p>
					<div className="modal-action">
						<button className="btn rounded-full" onClick={closeDeleteModal}>
							Cancel
						</button>
						<button className="btn btn-error rounded-full" onClick={handleDeleteProvider}>
							Delete provider
						</button>
					</div>
				</div>
			</dialog>
		</section>
	);
};
