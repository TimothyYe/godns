// components/TabControl.tsx
import React, { useState, useEffect } from 'react';
import { Domain } from '@/api/domain';
import { DomainCard } from '@/components/domain-card';
import { useContext } from 'react';
import { CommonContext } from '@/components/user';
import { get_domains } from '@/api/domain';

interface DomainControlProps {
	domains: Domain[];
}

export const DomainControl = () => {
	const userStore = useContext(CommonContext);
	const { credentials } = userStore;
	const [domains, setDomains] = useState<Domain[]>([]);

	useEffect(() => {
		if (!credentials) {
			window.location.href = '/login';
			return;
		}

		get_domains(credentials).then((domains) => {
			setDomains(domains);
		});
	}, [credentials, setDomains]);

	const onRemove = (domain: string) => {
		const newDomains = domains.filter((d) => d.domain_name !== domain).sort((a, b) => a.domain_name.localeCompare(b.domain_name));
		setDomains(newDomains);
	}

	const addNewDomain = () => {
		const newDomain: Domain = {
			domain_name: `Domain ${domains.length + 1}`,
			sub_domains: [`Subdomain ${domains.length + 1}`, `Subdomain ${domains.length + 1}`, "ipv6", "ipv4", "ddns"]
		};
		setDomains([...domains, newDomain]);
	};

	return (
		<div className="flex flex-col">
			<div className="flex flex-row justify-start">
				<button className="btn btn-primary btn-sm mb-5" onClick={addNewDomain}>Add Domain</button>
			</div>
			<div className="flex flex-wrap gap-2">
				{domains.map((domain, index) => (
					<DomainCard key={index} domain={domain} index={index} showActionBtn={true} onRemove={onRemove} />
				))}
			</div>
		</div>
	);
};

