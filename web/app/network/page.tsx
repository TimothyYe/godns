import { IpMode } from "@/components/ip-mode";

export default function Network() {
	return (
		<main className="flex min-h-screen flex-col items-center justify-between pt-10 max-w-screen-xl">
			<div className="p-10">
				<div className="flex flex-col max-w-screen-lg">
					<IpMode />
				</div>
			</div>
		</main>
	);
}