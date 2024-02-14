// components/TabControl.tsx
import React, { useState } from 'react';
import { Domain } from '@/api/info';
import { DomainCard } from '@/components/domain-card';

export const DomainControl = () => {
	const [domains, setDomains] = useState<Domain[]>([]);

	const addNewTab = () => {
		const newDomain: Domain = {
			domain_name: `Domain ${domains.length + 1}`,
			sub_domains: [`Subdomain 1 of domain ${domains.length + 1}`, `Subdomain 2 of domain ${domains.length + 1}`]
		};
		setDomains([...domains, newDomain]);
	};

	return (
		<div className="flex flex-col">
			<div className="flex flex-row justify-start">
				<button className="btn btn-primary btn-sm" onClick={addNewTab}>Add Domain</button>
			</div>
			<div className="flex flex-wrap gap-2">
				{domains.map((domain, index) => (
					<DomainCard key={index} domain={domain} isEdit={false} />
				))}
			</div>
		</div>
	);
};

