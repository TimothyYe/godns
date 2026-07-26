// components/TabControl.tsx
import React, { useState, useEffect, useRef } from 'react';
import { useRouter } from 'next/navigation';
import { Domain } from '@/api/domain';
import { DomainCard } from '@/components/domain-card';
import { useContext } from 'react';
import { CommonContext } from '@/components/user';
import { get_domains, add_domain, remove_domain } from '@/api/domain';
import { get_multi_providers } from '@/api/provider';
import { toast } from 'react-toastify';
import { WarningIcon, PlusIcon } from './icons';

export const DomainControl = () => {
	const router = useRouter();
	const userStore = useContext(CommonContext);
	const { credentials } = userStore;
	const [domains, setDomains] = useState<Domain[]>([]);
	const [showAlert, setShowAlert] = useState(false);
	const modalRef = useRef<HTMLDialogElement | null>(null);
	const [domainName, setDomainName] = useState<string>('');
	const [subDomains, setSubDomains] = useState<string[]>([]);
	const [selectedProvider, setSelectedProvider] = useState<string>('');
	const [availableProviders, setAvailableProviders] = useState<string[]>([]);

	const openModal = () => {
		// Reset form fields
		setDomainName('');
		setSubDomains([]);
		setSelectedProvider('');
		
		if (modalRef.current) {
			modalRef.current.showModal();
		}
	};

	useEffect(() => {
		if (!credentials) {
			router.push('/login');
			return;
		}

		// Load domains
		get_domains(credentials).then((domains) => {
			if (!domains || !Array.isArray(domains) || domains.length === 0) {
				setShowAlert(true);
			} else {
				setShowAlert(false);
				setDomains(domains.sort((a, b) => a.domain_name.localeCompare(b.domain_name)));
			}
		});

		// Load available providers
		get_multi_providers(credentials).then((providers) => {
			const providerNames = Object.keys(providers);
			setAvailableProviders(providerNames);
		});
	}, [credentials, router]);

	const onRemove = (domain: string) => {
		if (credentials) {
			remove_domain(credentials, domain).then((success) => {
				if (success) {
					toast.success('Domain removed successfully');
					const newDomains = domains.filter((d) => d.domain_name !== domain).sort((a, b) => a.domain_name.localeCompare(b.domain_name));
					setDomains(newDomains);
				} else {
					toast.error('Failed to remove domain');
				}
			});
		} else {
			toast.error('Invalid credentials');
		}
	}

	const addNewDomain = async () => {
		if (!domainName || !subDomains.length) {
			toast.error('Invalid domain or subdomain');
			return;
		}

		if (!selectedProvider) {
			toast.error('Please select a provider');
			return;
		}

		const newDomain: Domain = {
			domain_name: domainName,
			sub_domains: subDomains,
			provider: selectedProvider
		};

		if (credentials) {
			add_domain(credentials, newDomain).then((success) => {
				if (success) {
					setDomains([...domains, newDomain].sort((a, b) => a.domain_name.localeCompare(b.domain_name)));
					// reset the form fields
					setDomainName('');
					setSubDomains([]);
					setSelectedProvider('');
					// close the modal
					modalRef.current?.close();
					toast.success('Domain added successfully');
				} else {
					toast.error('Failed to add domain');
				}
			});
		} else {
			toast.error('Invalid credentials');
		}
	};

	return (
		<div className="w-full">
			<div className="section-header">
				<div className="section-copy">
					<div className="section-label ml-0">Domains</div>
					<h2 className="text-2xl font-semibold tracking-tight theme-heading">Managed domains</h2>
					<p className="mt-2 text-sm leading-7 theme-muted">
						Review the currently tracked root domains and add new ones without leaving the panel.
					</p>
				</div>
				<div className="flex items-center gap-3">
					{domains.length > 0 ? (
						<div className="theme-chip-sky rounded-full px-3 py-1 text-xs font-medium">
							{domains.length} domains
						</div>
					) : null}
					<button className="theme-primary-sky btn btn-sm rounded-xl border-none px-4 shadow-lg shadow-sky-900/20" onClick={openModal}>
						<PlusIcon />
						Add Domain
					</button>
				</div>
			</div>
			{
				showAlert ? (
					<div role="alert" className="mt-5 rounded-2xl border border-amber-500/20 bg-amber-500/10 px-4 py-4 text-amber-100">
						<WarningIcon />
						<span>Warning: No domains configured, please add a domain first!</span>
					</div>
				) : (

					<div className="mt-5 grid gap-4 lg:grid-cols-3">
						{domains.map((domain, index) => (
							<DomainCard key={index} domain={domain} index={index} showActionBtn={true} onRemove={onRemove} />
						))}
					</div>
				)
			}
			<dialog id="modal_add" className="modal modal-bottom sm:modal-middle" ref={modalRef}>
				<div className="theme-modal modal-box max-w-xl rounded-[1.5rem]">
					<h3 className="text-xl font-semibold tracking-tight theme-heading">Add Domain</h3>
					<p className="py-4 text-sm leading-7 theme-muted">Add a new domain to the configuration and map it to an existing provider profile.</p>
					<form method="dialog">
						<label className="theme-field mb-4">
							<span className="theme-field-label">Provider</span>
							<select
								className="select theme-input w-full rounded-2xl"
								value={selectedProvider}
								onChange={(e) => setSelectedProvider(e.target.value)}
							>
								<option value="">Select a provider</option>
								{availableProviders.map((provider) => (
									<option key={provider} value={provider}>
										{provider}
									</option>
								))}
							</select>
						</label>
						<label className="theme-field mb-4">
							<span className="theme-field-label">Domain</span>
							<input
								type="text"
								id="domain"
								placeholder="example.com"
								className="input theme-input w-full rounded-2xl"
								value={domainName}
								onChange={(e) => setDomainName(e.target.value)}
							/>
						</label>
						<label className="theme-field mb-4">
							<span className="theme-field-label">Subdomains</span>
							<textarea
								className="textarea theme-input h-36 rounded-2xl"
								placeholder={`subdomain1\nsubdomain2\nsubdomain3`}
								value={subDomains.join('\n')}
								onChange={(e) => setSubDomains(e.target.value.split('\n').filter(s => s.trim()))}
							/>
							<span className="theme-field-hint">Enter each subdomain on a new line.</span>
						</label>
						<div className="modal-action">
							<button 
								className="theme-subtle-btn btn mr-2 rounded-xl" 
								type="button"
								onClick={() => modalRef.current?.close()}
							>
								Cancel
							</button>
							<button 
								className="theme-primary-sky btn rounded-xl border-none" 
								type="button"
								onClick={addNewDomain}
								disabled={!domainName || !subDomains.length || !selectedProvider}
							>
								Add Domain
							</button>
						</div>
					</form>
				</div>
				<form method="dialog" className="modal-backdrop">
					<button aria-label="Close add domain dialog">close</button>
				</form>
			</dialog>
		</div>
	);
};
