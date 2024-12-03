import { Domain } from "@/api/domain";
import classNames from "classnames";
import { useRef } from "react";

interface DomainControlProps {
	domain: Domain;
	index: number;
	showActionBtn?: boolean;
	onRemove?: (domain: string) => void;
}

export const DomainCard = (props: DomainControlProps) => {
	const modalRef = useRef<HTMLDialogElement | null>(null);

	const openModal = () => {
		if (modalRef.current) {
			modalRef.current.showModal();
		}
	};

	const removeDomain = () => {
		if (props.onRemove) {
			props.onRemove(props.domain.domain_name);
		}
	};

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
							<button className="btn btn-secondary btn-sm" onClick={openModal}>Remove</button>
						</div>
					) : null
				}
			</div>
			<dialog id="modal_remove" className="modal" ref={modalRef}>
				<div className="modal-box">
					<h3 className="font-bold text-lg">Remove this domain?</h3>
					<p className="py-4">You will permanently remove this domain from the configuration.</p>
					<div className="modal-action">
						<form method="dialog">
							<button className="btn mr-2">Now now</button>
							<button className="btn btn-secondary" onClick={removeDomain} >Remove domain</button>
						</form>
					</div>
				</div>
			</dialog>
		</div >
	);
}