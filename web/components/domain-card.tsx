import { Domain } from "@/api/info";

interface DomainControlProps {
	domain: Domain;
	index: number;
	showActionBtn?: boolean;
}

export const DomainCard = (props: DomainControlProps) => {

	return (
		<div key="value" className={(props.index + 1) % 3 !== 0 ? "card w-full md:w-1/3 bg-primary-content shadow-xl" : "card w-full md:flex-1 bg-primary-content shadow-xl"}>
			<div className="card-body">
				<h2 className="card-title">
					{props.domain.domain_name}
				</h2>
				<div className="flex flex-wrap justify-start gap-2">
					{
						props.domain.sub_domains ? props.domain.sub_domains.map((sub_domain) => {
							return (
								<div key={sub_domain} className="badge badge-primary">{sub_domain}</div>
							);
						}) : null
					}
				</div>
				{
					props.showActionBtn ? (
						<div className="card-actions justify-end">
							<button className="btn btn-secondary btn-sm">Remove</button>
						</div>
					) : null
				}
			</div>
		</div>
	);
}