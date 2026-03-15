import { Domain } from "@/api/domain";
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
		modalRef.current?.showModal();
	};

	const removeDomain = () => {
		props.onRemove?.(props.domain.domain_name);
		modalRef.current?.close();
	};

	return (
		<>
			<article className="flex h-full flex-col rounded-[1.5rem] border border-base-300/70 bg-gradient-to-br from-base-100 to-base-200/60 p-5 shadow-sm">
				<div className="flex items-start justify-between gap-4">
					<div className="space-y-2">
						<h3 className="text-lg font-semibold tracking-tight">{props.domain.domain_name}</h3>
						<div className="chip-list">
							{props.domain.provider ? <div className="badge badge-secondary badge-lg">{props.domain.provider}</div> : null}
							<div className="badge badge-ghost badge-lg">{props.domain.sub_domains?.length || 0} subdomains</div>
						</div>
					</div>
					{props.showActionBtn ? (
						<button className="btn btn-ghost btn-sm rounded-full text-error hover:bg-error/10" onClick={openModal}>
							Remove
						</button>
					) : null}
				</div>

				<div className="mt-5 flex flex-1 flex-wrap gap-2">
					{props.domain.sub_domains && props.domain.sub_domains.length > 0 ? (
						props.domain.sub_domains.map((subDomain) => (
							<div key={subDomain} className="badge badge-primary badge-outline h-auto px-3 py-3">
								{subDomain}
							</div>
						))
					) : (
						<p className="text-sm text-base-content/55">No subdomains configured.</p>
					)}
				</div>
			</article>

			<dialog className="modal" ref={modalRef}>
				<div className="modal-box rounded-[1.75rem]">
					<h3 className="text-lg font-semibold">Delete domain</h3>
					<p className="py-3 text-sm text-base-content/65">
						Remove <strong>{props.domain.domain_name}</strong> from the GoDNS configuration. This does not delete records at the provider, but GoDNS will stop managing them.
					</p>
					<div className="modal-action">
						<button className="btn rounded-full" onClick={() => modalRef.current?.close()}>Cancel</button>
						<button className="btn btn-error rounded-full" onClick={removeDomain}>Delete domain</button>
					</div>
				</div>
			</dialog>
		</>
	);
};
