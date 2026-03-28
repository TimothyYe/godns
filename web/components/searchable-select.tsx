import classNames from 'classnames';
import React, { useEffect, useMemo, useRef, useState } from 'react';
import { createPortal } from 'react-dom';

interface SearchableSelectProps {
	options: { label: string; value: string }[];
	placeholder: string;
	defaultValue: string | undefined;
	onSelected?: (value: string) => void;
}

const SearchableSelect = ({ options, placeholder, defaultValue, onSelected }: SearchableSelectProps) => {
	const [searchTerm, setSearchTerm] = useState(defaultValue || '');
	const [showDropdown, setShowDropdown] = useState(false);
	const [mounted, setMounted] = useState(false);
	const [dropdownPosition, setDropdownPosition] = useState({ top: 0, left: 0, width: 0, maxHeight: 288 });
	const wrapperRef = useRef<HTMLDivElement | null>(null);
	const inputRef = useRef<HTMLInputElement | null>(null);
	const dropdownRef = useRef<HTMLDivElement | null>(null);

	const filteredOptions = useMemo(() => {
		return options.filter(option => option.label.toLowerCase().includes(searchTerm.toLowerCase()));
	}, [options, searchTerm]);

	useEffect(() => {
		setMounted(true);
	}, []);

	useEffect(() => {
		setSearchTerm(defaultValue || '');
	}, [defaultValue]);

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

		setDropdownPosition({
			top: openAbove ? Math.max(viewportPadding, rect.top - maxHeight - offset) : rect.bottom + offset,
			left: Math.max(viewportPadding, rect.left),
			width: rect.width,
			maxHeight,
		});
	};

	useEffect(() => {
		if (!showDropdown) {
			return;
		}

		updateDropdownPosition();
		window.addEventListener('resize', updateDropdownPosition);
		window.addEventListener('scroll', updateDropdownPosition, true);

		return () => {
			window.removeEventListener('resize', updateDropdownPosition);
			window.removeEventListener('scroll', updateDropdownPosition, true);
		};
	}, [showDropdown]);

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

	const portalTarget = mounted
		? ((wrapperRef.current?.closest('dialog') as HTMLElement | null) ?? document.body)
		: null;

	return (
		<div className="relative" ref={wrapperRef}>
			<input
				ref={inputRef}
				type="text"
				className="input theme-input w-full rounded-2xl"
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
				value={searchTerm}
				aria-autocomplete="list"
			/>
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
					<ul className="theme-dropdown-list">
						{filteredOptions.length > 0 ? filteredOptions.map((option) => (
							<li key={option.value}>
								<button
									type="button"
									className={classNames("theme-dropdown-item w-full text-left", {
										"theme-dropdown-item-active": option.value === searchTerm
									})}
									onMouseDown={(event) => {
										event.preventDefault();
										setSearchTerm(option.value);
										onSelected?.(option.value);
										setShowDropdown(false);
									}}
								>
									{option.label}
								</button>
							</li>
						)) : (
							<li>
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
