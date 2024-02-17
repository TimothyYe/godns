import { useEffect, useMemo, useState } from "react";
import { Provider, ProviderSetting } from "@/api/provider";
import { get_provider_settings, get_provider } from "@/api/provider";
import { useContext } from "react";
import { CommonContext } from "@/components/user";
import { useRouter } from "next/navigation";
import SearchableSelect from "./searchable-select";

export const ProviderControl = () => {
	const { credentials } = useContext(CommonContext);
	const [provider, setProvider] = useState<Provider>();
	const [providerSettings, setProviderSettings] = useState<ProviderSetting[]>([]);
	const router = useRouter();

	const options = useMemo(() => {
		if (providerSettings) {
			return providerSettings.map((setting) => {
				return {
					value: setting.name,
					label: setting.name
				};
			});
		}
		return [];
	}, [providerSettings]);

	useEffect(() => {
		if (!credentials) {
			router.push('/login');
			return;
		}
		// fetch provider settings
		get_provider_settings(credentials).then((settings) => {
			if (settings.length) {
				setProviderSettings(settings);
			}
		});
		// fetch provider
		get_provider(credentials).then((provider) => {
			setProvider(provider);
		});
	}, [credentials, router]);
	return (
		<div className="card w-96 sm:w-2/3 bg-primary-content shadow-xl">
			<div className="card-body">
				<h2 className="card-title">Provider Settings</h2>
				<SearchableSelect options={options} placeholder="Select Provider" defaultValue="" />
				<label className="input input-bordered flex items-center gap-2">
					Username
					<input type="text" className="grow" placeholder="Input the username" />
				</label>
				<label className="input input-bordered flex items-center gap-2">
					Email
					<input type="text" className="grow" placeholder="Input the email" />
				</label>
				<label className="input input-bordered flex items-center gap-2">
					Password
					<input type="password" className="grow" placeholder="Input the password" />
				</label>
				<div className="card-actions justify-end">
					<button className="btn btn-primary">Save</button>
				</div>
			</div>
		</div>
	);
}