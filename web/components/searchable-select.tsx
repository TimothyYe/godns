import classNames from 'classnames';
import React, { useState, useEffect, useId } from 'react';

interface SearchableSelectProps {
	options: { label: string; value: string }[];
	placeholder: string;
	defaultValue: string | undefined;
	onSelected?: (value: string) => void;
}

const SearchableSelect = ({ options, placeholder, defaultValue, onSelected }: SearchableSelectProps) => {
	const [searchTerm, setSearchTerm] = useState(defaultValue || '');
	const [showDropdown, setShowDropdown] = useState(false);
	const listboxId = useId();

	const filteredOptions = options.filter((option) =>
		option.label.toLowerCase().includes(searchTerm.toLowerCase())
	);

	useEffect(() => {
		setSearchTerm(defaultValue || '');
	}, [defaultValue]);

	const handleSelect = (value: string) => {
		setSearchTerm(value);
		onSelected?.(value);
		setShowDropdown(false);
	};

	return (
		<div className="relative">
			<input
				type="text"
				role="combobox"
				aria-controls={listboxId}
				aria-expanded={showDropdown}
				aria-autocomplete="list"
				className={classNames("input input-bordered h-12 w-full rounded-2xl", {
					'input-error': searchTerm.length === 0
				})}
				placeholder={placeholder}
				onFocus={() => setShowDropdown(true)}
				onBlur={() => setTimeout(() => setShowDropdown(false), 120)}
				onChange={(e) => {
					setSearchTerm(e.target.value);
					setShowDropdown(true);
				}}
				value={searchTerm}
			/>
			{showDropdown ? (
				<ul id={listboxId} className="menu dropdown-content absolute z-10 mt-2 max-h-64 w-full overflow-y-auto rounded-2xl border border-base-300 bg-base-100 p-2 shadow-xl">
					{filteredOptions.length > 0 ? (
						filteredOptions.map((option) => (
							<li key={option.value} onMouseDown={(event) => event.preventDefault()}>
								<button type="button" onClick={() => handleSelect(option.value)}>
									{option.label}
								</button>
							</li>
						))
					) : (
						<li className="pointer-events-none rounded-xl px-3 py-2 text-sm text-base-content/55">
							No matching options
						</li>
					)}
				</ul>
			) : null}
		</div>
	);
};

export default SearchableSelect;
