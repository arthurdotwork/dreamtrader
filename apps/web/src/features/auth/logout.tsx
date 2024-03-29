'use client';

import { useRouter } from 'next/navigation';
import { cn } from '@/lib/utils';
import { LogOutIcon } from 'lucide-react';
import Routes from '@/routes';

const Logout = ({ className }: { className?: string }) => {
	const router = useRouter();

	const handleLogout = () => {
		localStorage.removeItem('accessToken');
		localStorage.removeItem('refreshToken');

		fetch(Routes.API.Logout, {
			method: 'POST',
		}).then(() => {
			router.push(Routes.Auth);
		});
	};

	return (
		<>
			<p
				className={cn(className, 'text-muted-foreground text-sm cursor-pointer hover:text-red-600 transition flex items-center')}
				onClick={handleLogout}
			>
				<LogOutIcon className={"h-3 w-3 block mr-2"} />
				Log out
			</p>
		</>
	);
};

export default Logout;