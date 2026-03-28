import classNames from 'classnames';
import React, { useEffect, useMemo, useRef, useState } from 'react';
import { createPortal } from 'react-dom';

interface SearchableSelectProps {
	options: { label: string; value: string }[];
	placeholder: string;
	defaultValue: string | undefined;
	onSelected?: (value: string) => void;
}

type DropdownPosition = {
	top: number;
	left: number;
	width: number;
};

const SearchableSelect = ({ options, placeholder, defaultValue, onSelected }: SearchableSelectProps) => {
	const [searchTerm, setSearchTerm] = useState(defaultValue || '');
	const [showDropdown, setShowDropdown] = useState(false);
	const [mounted, setMounted] = useState(false);
	const [dropdownPosition, setDropdownPosition] = useState<DropdownPosition>({ top: 0, left: 0, width: 0 });
	const wrapperRef = useRef<HTMLDivElement | null>(null);
	const inputRef = useRef<HTMLInputElement | null>(null);

	const filteredOptions = useMemo(() => {
		return options.filter(option => option.label.toLowerCase().includes(searchTerm.toLowerCase()));
	}, [options, searchTerm]);

	const updateDropdownPosition = () => {
		const input = inputRef.current;
		if (!input) {
			return;
		}

		const rect = input.getBoundingClientRect();
		setDropdownPosition({
			top: rect.bottom + window.scrollY + 8,
			left: rect.left + window.scrollX,
			width: rect.width,
		});
	};

	useEffect(() => {
		setMounted(true);
	}, []);

	useEffect(() => {
		setSearchTerm(defaultValue || '');
	}, [defaultValue]);

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
			if (wrapperRef.current?.contains(target)) {
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
				className={classNames("input theme-input w-full rounded-2xl", {
					'border-rose-400/60': searchTerm.length === 0
				})}
				placeholder={placeholder}
				onFocus={() => {
					updateDropdownPosition();
					setShowDropdown(true);
				}}
				onChange={(e) => {
					setSearchTerm(e.target.value);
					if (!showDropdown) {
						updateDropdownPosition();
						setShowDropdown(true);
					}
				}}
				value={searchTerm}
			/>
			{portalTarget && showDropdown ? createPortal(
				<ul
					className="theme-dropdown menu fixed z-[120]"
					style={{
						top: `${dropdownPosition.top}px`,
						left: `${dropdownPosition.left}px`,
						width: `${dropdownPosition.width}px`,
					}}
				>
					{filteredOptions.length > 0 ? filteredOptions.map((option) => (
						<li key={option.value}>
							<button
								type="button"
								className="theme-dropdown-item w-full text-left"
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
				</ul>,
				portalTarget
			) : null}
		</div>
	);
};

export default SearchableSelect;
