'use client';
import { UserIcon } from "@/components/icons";
import { button as buttonStyles } from "@nextui-org/theme";
import { Link } from "@nextui-org/link";

export const LogoutButton = () => {
	return (
		<Link
			isExternal
			className={buttonStyles({ variant: "bordered", radius: "full" })}
			onClick={() => {
				localStorage.removeItem('token');
				window.location.href = '/';
			}}
		>
			<UserIcon size={20} />
			Logout
		</Link>
	);
}