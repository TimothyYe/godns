import { Domain } from "@/api/info";
import classNames from "classnames";

interface DomainControlProps {
	domain: Domain;
	index: number;
	showActionBtn?: boolean;
}

export const DomainCard = (props: DomainControlProps) => {

	return (
		<div key="value" className={classNames("card w-full bg-primary-content shadow-xl mb-1",
			{
				"md:w-1/3": (props.index + 1) % 3 !== 0,
				"md:flex-1": (props.index + 1) % 3 === 0
			})}>
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
		</div >
	);
}