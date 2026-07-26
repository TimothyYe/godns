import { Domain } from "@/api/domain";
import classNames from "classnames";
import { Fragment, useRef } from "react";

interface DomainControlProps {
	domain: Domain;
	index: number;
	showActionBtn?: boolean;
	onRemove?: (domain: string) => void;
}

export const DomainCard = (props: DomainControlProps) => {
	const modalRef = useRef<HTMLDialogElement | null>(null);
	const accentClass = 'surface-card-accent-sky';

	const openModal = () => {
		modalRef.current?.showModal();
	};

	const removeDomain = () => {
		props.onRemove?.(props.domain.domain_name);
	};

	return (
		<Fragment>
			<div className={classNames("surface-card flex h-full flex-col", accentClass)}>
				<div className="card-body gap-4 p-5 sm:p-6">
					<div className="flex items-start justify-between gap-4">
						<div className="space-y-2">
							<div className="metric-kicker">Domain</div>
							<h2 className="text-xl font-semibold tracking-tight theme-heading">
								{props.domain.domain_name}
							</h2>
							{props.domain.provider ? (
								<div className={classNames(
									"badge px-3 py-3",
									"theme-badge-sky"
								)}>
									{props.domain.provider}
								</div>
							) : null}
						</div>
						{props.showActionBtn ? (
							<button className="theme-danger btn btn-sm rounded-xl border-none px-4" onClick={openModal}>
								Remove
							</button>
						) : null}
					</div>

					<div className="space-y-3">
						<div className="metric-kicker">Subdomains</div>
						<div className="flex flex-wrap gap-2">
						{props.domain.sub_domains ? props.domain.sub_domains.map((sub_domain) => (
							<div key={sub_domain} className="theme-chip badge rounded-full px-3 py-3">
								{sub_domain}
							</div>
						)) : null}
						</div>
					</div>
				</div>
			</div>

			<dialog id="modal_remove" className="modal modal-bottom sm:modal-middle" ref={modalRef}>
				<div className="theme-modal modal-box rounded-[1.5rem]">
					<h3 className="text-lg font-semibold theme-heading">Remove this domain?</h3>
					<p className="py-4 text-sm leading-7 theme-muted">
						You will permanently remove <span className="font-medium theme-heading">{props.domain.domain_name}</span> from the current configuration.
					</p>
					<div className="modal-action">
						<form method="dialog" className="flex gap-2">
							<button className="theme-subtle-btn btn rounded-xl">Cancel</button>
							<button className="theme-danger btn rounded-xl border-none" onClick={removeDomain}>Remove domain</button>
						</form>
					</div>
				</div>
				<form method="dialog" className="modal-backdrop">
					<button aria-label="Close remove domain dialog">close</button>
				</form>
			</dialog>
		</Fragment>
	);
}
