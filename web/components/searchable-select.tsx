import classNames from 'classnames';
import React, { useEffect, useId, useMemo, useRef, useState } from 'react';
import { createPortal } from 'react-dom';
import { useIsHydrated } from './use-is-hydrated';

interface SearchableSelectProps {
	options: { label: string; value: string }[];
	placeholder: string;
	defaultValue: string | undefined;
	onSelected?: (value: string) => void;
}

const SearchableSelect = ({ options, placeholder, defaultValue, onSelected }: SearchableSelectProps) => {
	const [searchTerm, setSearchTerm] = useState(defaultValue || '');
	const [showDropdown, setShowDropdown] = useState(false);
	const [activeIndex, setActiveIndex] = useState(0);
	const [dropdownPosition, setDropdownPosition] = useState({ top: 0, left: 0, width: 0, maxHeight: 288 });
	const [portalTarget, setPortalTarget] = useState<HTMLElement | null>(null);
	const wrapperRef = useRef<HTMLDivElement | null>(null);
	const inputRef = useRef<HTMLInputElement | null>(null);
	const dropdownRef = useRef<HTMLDivElement | null>(null);
	const inputId = useId();
	const listboxId = `${inputId}-listbox`;
	const isHydrated = useIsHydrated();

	const filteredOptions = useMemo(() => {
		return options.filter(option => option.label.toLowerCase().includes(searchTerm.toLowerCase()));
	}, [options, searchTerm]);

	useEffect(() => {
		setSearchTerm(defaultValue || '');
	}, [defaultValue]);

	useEffect(() => {
		if (!showDropdown) {
			return;
		}

		const exactMatchIndex = filteredOptions.findIndex((option) => option.value === searchTerm);
		setActiveIndex(exactMatchIndex >= 0 ? exactMatchIndex : 0);
	}, [showDropdown, filteredOptions, searchTerm]);

	const updateDropdownPosition = () => {
		const input = inputRef.current;
		if (!input) {
			return;
		}

		const rect = input.getBoundingClientRect();
		const viewportPadding = 16;
		const offset = 8;
		const availableBelow = window.innerHeight - rect.bottom - viewportPadding - offset;
		const availableAbove = rect.top - viewportPadding - offset;
		const openAbove = availableBelow < 220 && availableAbove > availableBelow;
		const maxHeight = Math.max(160, Math.min(320, openAbove ? availableAbove : availableBelow));
		const width = Math.min(rect.width, window.innerWidth - viewportPadding * 2);
		const left = Math.min(
			Math.max(viewportPadding, rect.left),
			window.innerWidth - viewportPadding - width
		);

		setDropdownPosition({
			top: openAbove ? Math.max(viewportPadding, rect.top - maxHeight - offset) : rect.bottom + offset,
			left,
			width,
			maxHeight,
		});
	};

	useEffect(() => {
		if (!showDropdown || !isHydrated) {
			return;
		}

		setPortalTarget((wrapperRef.current?.closest('dialog') as HTMLElement | null) ?? document.body);
		updateDropdownPosition();
		window.addEventListener('resize', updateDropdownPosition);
		window.addEventListener('scroll', updateDropdownPosition, true);

		return () => {
			window.removeEventListener('resize', updateDropdownPosition);
			window.removeEventListener('scroll', updateDropdownPosition, true);
		};
	}, [showDropdown, isHydrated]);

	useEffect(() => {
		if (!showDropdown) {
			return;
		}

		const handlePointerDown = (event: MouseEvent) => {
			const target = event.target as Node;
			if (wrapperRef.current?.contains(target) || dropdownRef.current?.contains(target)) {
				return;
			}

			setShowDropdown(false);
		};

		document.addEventListener('mousedown', handlePointerDown);
		return () => document.removeEventListener('mousedown', handlePointerDown);
	}, [showDropdown]);

	const activeOption = filteredOptions[activeIndex];

	const commitSelection = (value: string) => {
		setSearchTerm(value);
		onSelected?.(value);
		setShowDropdown(false);
	};

	return (
		<div className="relative" ref={wrapperRef}>
			<div className="relative">
				<input
					ref={inputRef}
					id={inputId}
					type="text"
					role="combobox"
					className="input theme-input w-full rounded-2xl pr-12"
					placeholder={placeholder}
					onFocus={() => {
						setShowDropdown(true);
					}}
					onChange={(e) => {
						setSearchTerm(e.target.value);
						if (!showDropdown) {
							setShowDropdown(true);
						}
					}}
					onKeyDown={(event) => {
						if (event.key === 'ArrowDown') {
							event.preventDefault();
							if (!showDropdown) {
								setShowDropdown(true);
								return;
							}
							setActiveIndex((current) => Math.min(current + 1, Math.max(filteredOptions.length - 1, 0)));
						}

						if (event.key === 'ArrowUp') {
							event.preventDefault();
							if (!showDropdown) {
								setShowDropdown(true);
								return;
							}
							setActiveIndex((current) => Math.max(current - 1, 0));
						}

						if (event.key === 'Enter' && showDropdown && activeOption) {
							event.preventDefault();
							commitSelection(activeOption.value);
						}

						if (event.key === 'Escape') {
							setShowDropdown(false);
						}

						if (event.key === 'Tab') {
							setShowDropdown(false);
						}
					}}
					value={searchTerm}
					aria-autocomplete="list"
					aria-expanded={showDropdown}
					aria-controls={listboxId}
					aria-activedescendant={showDropdown && activeOption ? `${listboxId}-${activeOption.value}` : undefined}
				/>
				<span className="theme-combobox-icon pointer-events-none absolute inset-y-0 right-4 flex items-center">
					<svg aria-hidden="true" className={classNames("h-4 w-4 transition-transform", { "rotate-180": showDropdown })} viewBox="0 0 20 20" fill="currentColor">
						<path fillRule="evenodd" d="M5.23 7.21a.75.75 0 0 1 1.06.02L10 11.168l3.71-3.938a.75.75 0 1 1 1.08 1.04l-4.25 4.5a.75.75 0 0 1-1.08 0l-4.25-4.5a.75.75 0 0 1 .02-1.06Z" clipRule="evenodd" />
					</svg>
				</span>
			</div>
			{portalTarget && showDropdown ? createPortal(
				<div
					ref={dropdownRef}
					className="theme-dropdown fixed z-[260]"
					style={{
						top: `${dropdownPosition.top}px`,
						left: `${dropdownPosition.left}px`,
						width: `${dropdownPosition.width}px`,
						maxHeight: `${dropdownPosition.maxHeight}px`,
					}}
				>
					<ul className="theme-dropdown-list" id={listboxId} role="listbox" aria-labelledby={inputId}>
						{filteredOptions.length > 0 ? filteredOptions.map((option, index) => (
							<li key={option.value} role="presentation">
								<button
									type="button"
									id={`${listboxId}-${option.value}`}
									role="option"
									aria-selected={index === activeIndex}
									className={classNames("theme-dropdown-item w-full text-left", {
										"theme-dropdown-item-active": index === activeIndex
									})}
									onMouseEnter={() => setActiveIndex(index)}
									onMouseDown={(event) => {
										event.preventDefault();
										commitSelection(option.value);
									}}
								>
									{option.label}
								</button>
							</li>
						)) : (
							<li role="presentation">
								<div className="theme-dropdown-item cursor-default text-sm theme-faint">No matching providers</div>
							</li>
						)}
					</ul>
				</div>,
				portalTarget
			) : null}
		</div>
	);
};

export default SearchableSelect;
