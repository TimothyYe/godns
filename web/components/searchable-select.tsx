import classNames from 'classnames';
import React, { useState, useEffect } from 'react';

interface SearchableSelectProps {
	options: { label: string; value: string }[];
	placeholder: string;
	defaultValue: string | undefined;
	onSelected?: (value: string) => void;
}

const SearchableSelect = ({ options, placeholder, defaultValue, onSelected }: SearchableSelectProps) => {
	const [searchTerm, setSearchTerm] = useState(defaultValue || '');
	const [showDropdown, setShowDropdown] = useState(false);

	const filteredOptions = options.filter(option =>
		option.label.toLowerCase().includes(searchTerm.toLowerCase())
	);

	useEffect(() => {
		setSearchTerm(defaultValue || '');
	}, [defaultValue]);

	return (
		<div className="relative">
			<input
				type="text"
				className={classNames("input input-bordered w-full", {
					'input-error': searchTerm.length === 0
				})}
				placeholder={placeholder}
				onFocus={() => setShowDropdown(true)}
				onBlur={() => setTimeout(() => setShowDropdown(false), 200)}
				onChange={(e) => {
					setSearchTerm(e.target.value);
				}}
				value={searchTerm}
			/>
			{showDropdown && (
				<ul className="menu dropdown-content p-2 shadow bg-base-100 rounded-box w-full absolute z-10">
					{filteredOptions.map((option) => (
						<li key={option.value} onClick={() => {
							setSearchTerm(option.value);
							if (onSelected) {
								onSelected(option.value);
							}
						}}>
							<a>{option.label}</a>
						</li>
					))}
				</ul>
			)}
		</div>
	);
};

export default SearchableSelect;
