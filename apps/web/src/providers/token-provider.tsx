'use client';

import {QueryClient} from "@tanstack/react-query";
import { ReactNode, useEffect } from 'react';

type TokenProviderProps = {
	children: ReactNode;
	accessToken: string;
	refreshToken: string;
}

const TokenProvider = ({children, accessToken, refreshToken}: TokenProviderProps) => {
	useEffect(() => {
		localStorage.setItem('accessToken', accessToken);
		localStorage.setItem('refreshToken', refreshToken);
	}, [accessToken, refreshToken]);

	return (
		<>
			{children}
		</>
	);
}

export default TokenProvider;