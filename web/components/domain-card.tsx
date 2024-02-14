'use client'

import { useState } from "react";
import { Domain } from "@/api/info";

interface DomainControlProps {
	domain: Domain;
	isEdit: boolean;
}

export const DomainCard = (props: DomainControlProps) => {
	const [isEdit, setEdit] = useState(props.isEdit);

	return (
		<div className="card bg-primary-content">
			<div className="card-body">
				<div className="form-control">
					<label className="label">
						<span className="label-text">Domain Name</span>
					</label>
					<input type="text" placeholder="Domain Name" value={props.domain.domain_name} className="input input-bordered" disabled={!isEdit} />
					<label className="label">
						<span className="label-text">Subdomains</span>
						{
							props.domain.sub_domains.map((sub_domain, index) => (
								<input key={index} type="text" placeholder="Subdomain" value={sub_domain} className="input input-bordered" disabled={!isEdit} />
							))
						}
					</label>
				</div>
			</div>
		</div>
	);
}