import { title, subtitle } from "@/components/primitives";
import { Card, CardHeader, CardBody, Button, Slider } from "@nextui-org/react";
import { DotIcon } from "@/components/icons";

export default function Home() {
	return (
		<section className="flex flex-col items-center justify-center gap-4 py-8 md:py-10">
			{/* <div className=" max-w-lg text-center justify-center">
				<h1 className={title()}>Make&nbsp;</h1>
				<h1 className={title({ color: "violet" })}>beautiful&nbsp;</h1>
				<br />
				<h1 className={title()}>
					websites regardless of your design experience.
				</h1>
				<h2 className={subtitle({ class: "mt-4" })}>
					Beautiful, fast and modern React UI library.
				</h2>
			</div> */}
			<div className="w-1/2 text-center justify-center">
				<Card>
					<CardHeader>
						<div className="flex flex-col">
							<p className="text-md">NextUI</p>
							<DotIcon className="text-green-500" />
							<p className="text-small text-default-500">nextui.org</p>
						</div>
					</CardHeader>
					<CardBody>
						<h2 className='text-md font-bold text-center'></h2>
					</CardBody>
				</Card>
			</div>
			<div className="w-1/2 text-center justify-center">
				<Card
					className="border-none bg-background/60 dark:bg-default-100/50 m-20"
					shadow="sm"
				>
					<CardBody>
						<div className="grid grid-cols-6 md:grid-cols-12 gap-6 md:gap-4 items-center justify-center">
							<div className="relative col-span-6 md:col-span-4">
							</div>

							<div className="flex flex-col col-span-6 md:col-span-8">
								<div className="flex justify-between items-start">
									<div className="flex flex-col gap-0">
										<h3 className="font-semibold text-foreground/90">Daily Mix</h3>
										<p className="text-small text-foreground/80">12 Tracks</p>
										<h1 className="text-large font-medium mt-2">Frontend Radio</h1>
									</div>
									<Button
										isIconOnly
										className="text-default-900/60 data-[hover]:bg-foreground/10 -translate-y-2 translate-x-2"
										radius="full"
										variant="light"
									// onPress={() => setLiked((v) => !v)}
									>
										BTN
									</Button>
								</div>

								<div className="flex flex-col mt-3 gap-1">
									<Slider
										aria-label="Music progress"
										classNames={{
											track: "bg-default-500/30",
											thumb: "w-2 h-2 after:w-2 after:h-2 after:bg-foreground",
										}}
										color="foreground"
										defaultValue={33}
										size="sm"
									/>
									<div className="flex justify-between">
										<p className="text-small">1:23</p>
										<p className="text-small text-foreground/50">4:32</p>
									</div>
								</div>

								<div className="flex w-full items-center justify-center">
									<Button
										isIconOnly
										className="data-[hover]:bg-foreground/10"
										radius="full"
										variant="light"
									>
										{/* <RepeatOneIcon className="text-foreground/80" /> */}
									</Button>
									<Button
										isIconOnly
										className="data-[hover]:bg-foreground/10"
										radius="full"
										variant="light"
									>
										{/* <PreviousIcon /> */}
									</Button>
									<Button
										isIconOnly
										className="w-auto h-auto data-[hover]:bg-foreground/10"
										radius="full"
										variant="light"
									>
										{/* <PauseCircleIcon size={54} /> */}
									</Button>
									<Button
										isIconOnly
										className="data-[hover]:bg-foreground/10"
										radius="full"
										variant="light"
									>
										{/* <NextIcon /> */}
									</Button>
									<Button
										isIconOnly
										className="data-[hover]:bg-foreground/10"
										radius="full"
										variant="light"
									>
										{/* <ShuffleIcon className="text-foreground/80" /> */}
									</Button>
								</div>
							</div>
						</div>
					</CardBody>
				</Card>
			</div>
		</section>
	);
}
