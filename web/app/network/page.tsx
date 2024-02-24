import { IpMode } from "@/components/ip-mode";
import { Proxy } from "@/components/proxy";
import { WebHook } from "@/components/webhook";
import { Resolver } from "@/components/resolver";
import { IPInterface } from "@/components/ip-interface";

export default function Network() {
	return (
		<main className="flex min-h-screen flex-col items-center justify-between pt-10 max-w-screen-xl">
			<div className="p-5">
				<div className="flex flex-col max-w-screen-lg gap-5">
					<IpMode />
					<Proxy />
					<WebHook />
					<Resolver />
					<IPInterface />
					<div className="flex justify-center">
						<button className="flex btn btn-primary">Save</button>
					</div>
				</div>
			</div>
		</main>
	);
}