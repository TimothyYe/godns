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
			if (!domains) {
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
			<div className="flex items-center justify-between mb-4">
				<h2 className="text-xl font-semibold text-neutral-500">Domain Settings</h2>
				<button className="btn btn-primary btn-sm" onClick={openModal}>
					<PlusIcon />
					Add Domain
				</button>
			</div>
			{
				showAlert ? (
					<div role="alert" className="alert alert-warning">
						<WarningIcon />
						<span>Warning: No domains configured, please add a domain first!</span>
					</div>
				) : (

					<div className="flex flex-wrap gap-2">
						{domains.map((domain, index) => (
							<DomainCard key={index} domain={domain} index={index} showActionBtn={true} onRemove={onRemove} />
						))}
					</div>
				)
			}
			<dialog id="modal_add" className="modal" ref={modalRef}>
				<div className="modal-box max-w-lg">
					<h3 className="font-bold text-lg">Add Domain</h3>
					<p className="py-4">Add a new domain to the configuration.</p>
					<form method="dialog">
						<label className="form-control w-full mb-4">
							<div className="label">
								<span className="label-text font-bold">Provider</span>
							</div>
							<select
								className="select select-primary select-bordered w-full"
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
						<label className="form-control w-full mb-4">
							<div className="label">
								<span className="label-text font-bold">Domain</span>
							</div>
							<input
								type="text"
								id="domain"
								placeholder="Input the domain name"
								className="input input-primary input-bordered w-full"
								value={domainName}
								onChange={(e) => setDomainName(e.target.value)}
							/>
						</label>
						<label className="form-control w-full mb-4">
							<div className="label">
								<span className="label-text font-bold">Subdomains</span>
							</div>
							<textarea
								className="textarea textarea-primary h-36"
								placeholder={`subdomain1\nsubdomain2\nsubdomain3`}
								value={subDomains.join('\n')}
								onChange={(e) => setSubDomains(e.target.value.split('\n').filter(s => s.trim()))}
							/>
							<div className="label">
								<span className="label-text-alt">Enter each subdomain on a new line</span>
							</div>
						</label>
						<div className="modal-action">
							<button 
								className="btn mr-2" 
								type="button"
								onClick={() => modalRef.current?.close()}
							>
								Close
							</button>
							<button 
								className="btn btn-primary" 
								type="button"
								onClick={addNewDomain}
								disabled={!domainName || !subDomains.length || !selectedProvider}
							>
								Add Domain
							</button>
						</div>
					</form>
				</div>
			</dialog>
		</div>
	);
};

