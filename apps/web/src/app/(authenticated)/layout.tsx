import Logout from '@/features/auth/logout';
import TokenProvider from '@/providers/token-provider';
import { cookies } from 'next/headers';

const Layout = ({ children }: Readonly<{ children: React.ReactNode }>) => {
	const accessToken = cookies().get('accessToken');
	const refreshToken = cookies().get('refreshToken');

	return (
		<TokenProvider accessToken={accessToken?.value ?? ''} refreshToken={refreshToken?.value ?? ''}>
			<div className="h-screen flex w-full">
				<aside className="border-r border-gray-200 w-72 bg-gray-100 p-6 flex flex-col">
					<Logout className={"mt-auto"} />
				</aside>
				<main className="p-6">
					{children}
				</main>
			</div>
		</TokenProvider>
	)
}

export default Layout;