import classNames from 'classnames';
import type { ReactNode } from 'react';

interface PageShellProps {
	eyebrow?: string;
	title: string;
	description?: string;
	actions?: ReactNode;
	children: ReactNode;
	className?: string;
}

interface SectionCardProps {
	title: string;
	description?: string;
	actions?: ReactNode;
	children: ReactNode;
	className?: string;
}

export const PageShell = ({
	eyebrow,
	title,
	description,
	actions,
	children,
	className,
}: PageShellProps) => {
	return (
		<div className={classNames('page-shell', className)}>
			<header className="page-hero">
				<div className="space-y-3">
					{eyebrow ? <span className="hero-badge">{eyebrow}</span> : null}
					<div className="space-y-2">
						<h1 className="page-title">{title}</h1>
						{description ? <p className="page-subtitle">{description}</p> : null}
					</div>
				</div>
				{actions ? <div className="flex flex-wrap items-center gap-3">{actions}</div> : null}
			</header>
			<div className="space-y-6">{children}</div>
		</div>
	);
};

export const SectionCard = ({
	title,
	description,
	actions,
	children,
	className,
}: SectionCardProps) => {
	return (
		<section className={classNames('panel', className)}>
			<div className="section-header">
				<div className="space-y-1">
					<h2 className="section-title">{title}</h2>
					{description ? <p className="section-subtitle">{description}</p> : null}
				</div>
				{actions ? <div className="flex flex-wrap items-center gap-2">{actions}</div> : null}
			</div>
			{children}
		</section>
	);
};
