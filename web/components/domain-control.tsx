import React, { useState, useEffect, useRef, useMemo } from 'react';
import { useRouter } from 'next/navigation';
import { Domain } from '@/api/domain';
import { DomainCard } from '@/components/domain-card';
import { useContext } from 'react';
import { CommonContext } from '@/components/user';
import { get_domains, add_domain, remove_domain } from '@/api/domain';
import { get_multi_providers } from '@/api/provider';
import { toast } from 'react-toastify';
import { PlusIcon, WarningIcon } from './icons';

export const DomainControl = () => {
	const router = useRouter();
	const userStore = useContext(CommonContext);
	const { credentials } = userStore;
	const [domains, setDomains] = useState<Domain[]>([]);
	const [loading, setLoading] = useState(true);
	const modalRef = useRef<HTMLDialogElement | null>(null);
	const [domainName, setDomainName] = useState<string>('');
	const [subDomainText, setSubDomainText] = useState<string>('');
	const [selectedProvider, setSelectedProvider] = useState<string>('');
	const [availableProviders, setAvailableProviders] = useState<string[]>([]);
	const [searchTerm, setSearchTerm] = useState('');

	const sortedDomains = useMemo(
		() => [...domains].sort((a, b) => a.domain_name.localeCompare(b.domain_name)),
		[domains]
	);

	const filteredDomains = sortedDomains.filter((domain) => {
		const query = searchTerm.toLowerCase().trim();
		if (!query) {
			return true;
		}

		return domain.domain_name.toLowerCase().includes(query)
			|| domain.provider?.toLowerCase().includes(query)
			|| domain.sub_domains.some((subDomain) => subDomain.toLowerCase().includes(query));
	});

	const subDomains = subDomainText
		.split('\n')
		.map((item) => item.trim())
		.filter(Boolean);

	const openModal = () => {
		setDomainName('');
		setSubDomainText('');
		setSelectedProvider('');
		modalRef.current?.showModal();
	};

	useEffect(() => {
		if (!credentials) {
			router.push('/login');
			return;
		}

		setLoading(true);
		Promise.all([
			get_domains(credentials),
			get_multi_providers(credentials),
		]).then(([nextDomains, providers]) => {
			setDomains(nextDomains || []);
			setAvailableProviders(Object.keys(providers || {}));
		}).finally(() => {
			setLoading(false);
		});
	}, [credentials, router]);

	const onRemove = (domain: string) => {
		if (!credentials) {
			toast.error('Invalid credentials');
			return;
		}

		remove_domain(credentials, domain).then((success) => {
			if (success) {
				toast.success('Domain removed successfully');
				setDomains((current) => current.filter((item) => item.domain_name !== domain));
			} else {
				toast.error('Failed to remove domain');
			}
		});
	};

	const addNewDomain = async () => {
		if (!domainName || !subDomains.length) {
			toast.error('Enter a domain and at least one subdomain.');
			return;
		}

		if (!selectedProvider) {
			toast.error('Select a provider profile for this domain.');
			return;
		}

		const newDomain: Domain = {
			domain_name: domainName.trim(),
			sub_domains: subDomains,
			provider: selectedProvider
		};

		if (!credentials) {
			toast.error('Invalid credentials');
			return;
		}

		const success = await add_domain(credentials, newDomain);
		if (!success) {
			toast.error('Failed to add domain');
			return;
		}

		setDomains((current) => [...current, newDomain]);
		modalRef.current?.close();
		toast.success('Domain added successfully');
	};

	return (
		<section className="panel">
			<div className="section-header">
				<div className="space-y-1">
					<h2 className="section-title">Domains</h2>
					<p className="section-subtitle">
						Attach each domain to a provider profile and list every subdomain that GoDNS should keep updated.
					</p>
				</div>
				<div className="flex flex-wrap gap-2">
					<input
						type="search"
						className="input input-bordered input-sm h-10 rounded-full"
						placeholder="Search domains or subdomains"
						value={searchTerm}
						onChange={(e) => setSearchTerm(e.target.value)}
					/>
					<button className="btn btn-primary rounded-full px-5" onClick={openModal} disabled={availableProviders.length === 0}>
						<PlusIcon />
						Add domain
					</button>
				</div>
			</div>

			<div className="mb-5 flex flex-wrap items-center gap-3 text-sm text-base-content/60">
				<span className="badge badge-ghost badge-lg">{sortedDomains.length} configured</span>
				<span className="badge badge-ghost badge-lg">{availableProviders.length} providers available</span>
			</div>

			{loading ? (
				<div className="flex min-h-48 items-center justify-center">
					<span className="loading loading-spinner loading-lg" />
				</div>
			) : availableProviders.length === 0 ? (
				<div role="alert" className="panel-muted flex items-start gap-3">
					<WarningIcon />
					<div className="space-y-1">
						<p className="font-semibold">Add a provider before creating domains</p>
						<p className="text-sm text-base-content/65">
							Each domain must be mapped to a provider profile so GoDNS knows where updates should be sent.
						</p>
					</div>
				</div>
			) : filteredDomains.length > 0 ? (
				<div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
					{filteredDomains.map((domain, index) => (
						<DomainCard key={`${domain.domain_name}-${index}`} domain={domain} index={index} showActionBtn={true} onRemove={onRemove} />
					))}
				</div>
			) : sortedDomains.length > 0 ? (
				<div className="panel-muted">
					<p className="font-semibold">No domains match the current filter</p>
					<p className="mt-2 text-sm text-base-content/60">Try a different domain name, provider, or subdomain term.</p>
				</div>
			) : (
				<div className="panel-muted">
					<p className="font-semibold">No domains configured yet</p>
					<p className="mt-2 text-sm text-base-content/60">
						Add your first domain once provider credentials are in place.
					</p>
				</div>
			)}

			<dialog className="modal" ref={modalRef}>
				<div className="modal-box max-w-2xl rounded-[1.75rem]">
					<h3 className="text-xl font-semibold tracking-tight">Add domain</h3>
					<p className="pb-5 pt-2 text-sm text-base-content/65">
						Create a domain entry and list the subdomains GoDNS should manage. One subdomain per line keeps the configuration readable.
					</p>
					<div className="grid gap-5">
						<label className="field-stack">
							<span className="field-label">Provider profile</span>
							<select
								className="select select-bordered h-12 w-full rounded-2xl"
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

						<label className="field-stack">
							<span className="field-label">Domain name</span>
							<input
								type="text"
								id="domain"
								placeholder="example.com"
								className="input input-bordered h-12 w-full rounded-2xl"
								value={domainName}
								onChange={(e) => setDomainName(e.target.value)}
							/>
							<span className="field-hint">Use the root domain, not a full record name.</span>
						</label>

						<label className="field-stack">
							<span className="field-label">Subdomains</span>
							<textarea
								className="textarea textarea-bordered h-40 w-full rounded-[1.25rem]"
								placeholder={`@\nwww\nhome`}
								value={subDomainText}
								onChange={(e) => setSubDomainText(e.target.value)}
							/>
							<span className="field-hint">Enter one subdomain per line. Use <code>@</code> for the apex/root record.</span>
						</label>
					</div>
					<div className="modal-action">
						<button className="btn rounded-full" type="button" onClick={() => modalRef.current?.close()}>
							Cancel
						</button>
						<button
							className="btn btn-primary rounded-full px-5"
							type="button"
							onClick={addNewDomain}
							disabled={!domainName || !subDomains.length || !selectedProvider}
						>
							Add domain
						</button>
					</div>
				</div>
			</dialog>
		</section>
	);
};
