// components/TabControl.tsx
import React, { useState, useEffect, useRef } from 'react';
import { useRouter } from 'next/navigation';
import { Domain } from '@/api/domain';
import { DomainCard } from '@/components/domain-card';
import { useContext } from 'react';
import { CommonContext } from '@/components/user';
import { get_domains, add_domain, remove_domain } from '@/api/domain';
import { toast } from 'react-toastify';
import { WarningIcon } from './icons';

export const DomainControl = () => {
	const router = useRouter();
	const userStore = useContext(CommonContext);
	const { credentials } = userStore;
	const [domains, setDomains] = useState<Domain[]>([]);
	const [showAlert, setShowAlert] = useState(false);
	const modalRef = useRef<HTMLDialogElement | null>(null);

	const openModal = () => {
		if (modalRef.current) {
			modalRef.current.showModal();
		}
	};

	useEffect(() => {
		if (!credentials) {
			router.push('/login');
			return;
		}

		get_domains(credentials).then((domains) => {
			if (!domains) {
				setShowAlert(true);
			} else {
				setShowAlert(false);
				setDomains(domains.sort((a, b) => a.domain_name.localeCompare(b.domain_name)));
			}
		});
	}, [credentials, setDomains, router]);

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
		const newDomain: Domain = {
			domain_name: `sample.com`,
			sub_domains: ["ipv6", "ipv4", "ddns"]
		};

		if (credentials) {
			add_domain(credentials, newDomain).then((success) => {
				if (success) {
					setDomains([...domains, newDomain]);
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
		<div className="flex flex-col">
			<div className="flex flex-row justify-start">
				<button className="btn btn-primary btn-sm mb-5" onClick={openModal}>Add Domain</button>
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
				<div className="modal-box">
					<h3 className="font-bold text-lg">Add domain</h3>
					<p className="py-4">Add a new domain to the configuration.</p>
					<div className="modal-action">
						<form method="dialog">
							<button className="btn mr-2">Close</button>
							<button className="btn btn-primary" onClick={addNewDomain} >Save</button>
						</form>
					</div>
				</div>
			</dialog>
		</div>
	);
};

