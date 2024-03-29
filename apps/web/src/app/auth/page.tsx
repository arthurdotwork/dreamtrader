import AuthenticationScreen from '@/features/auth/auth';
import { cookies } from 'next/headers';

const Page = () => {
	return (
		<div
			className="h-screen flex flex-col justify-center items-center px-4 lg:px-0 bg-neutral-50">
			<AuthenticationScreen />
		</div>
	);
};

export default Page;