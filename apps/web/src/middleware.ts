import { NextRequest, NextResponse } from 'next/server';
import camelcaseKeys from 'camelcase-keys';

const authenticatedMiddleware = async (request: NextRequest) => {
	const accessToken = request.cookies.get('accessToken');
	const refreshToken = request.cookies.get('refreshToken');

	if (!accessToken || !refreshToken) {
		return NextResponse.redirect(new URL('/auth', request.url));
	}

	try {
		const req = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/v1/auth/verify`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
				'Authorization': `Bearer ${accessToken.value}`,
			},
		});

		if (req.status === 401) {
			const refreshReq = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/v1/auth/refresh`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application',
					'Authorization': `Bearer ${refreshToken.value}`,
				}
			});

			if (refreshReq.status !== 200) {
				throw new Error('Failed to refresh token');
			}

			const refreshRes = await refreshReq.json();
			const {accessToken} = camelcaseKeys(refreshRes);

			const req = NextResponse.redirect(new URL('/', request.url));
			req.cookies.set('accessToken', accessToken.token, {
				maxAge: 60 * 60 * 24 * 7,
				httpOnly: true,
				sameSite: 'strict',
				secure: true
			});

			req.cookies.set('refreshToken', refreshToken.value, {
				maxAge: 60 * 60 * 24 * 7,
				httpOnly: true,
				sameSite: 'strict',
				secure: true
			});

			return req;
		}
	} catch (e) {
		console.error({e})

		const req = NextResponse.redirect(new URL('/auth', request.url))
		req.cookies.delete('accessToken');
		req.cookies.delete('refreshToken');
		return req;
	}


	return NextResponse.next();
}

const isAuthenticatedMiddleware = (request: NextRequest) => {
	const accessToken = request.cookies.get('accessToken');
	const refreshToken = request.cookies.get('refreshToken');

	if (accessToken && refreshToken) {
		return NextResponse.redirect(new URL('/', request.url));
	}

	return NextResponse.next();
}

const nonMiddlewarePaths = ["/api", "/_next", "/favicon.ico", "/manifest/webmanifest"];

export async function middleware(request: NextRequest) {
	const incomingURL = new URL(request.url);

	if (['/auth', '/auth/reset', '/register'].includes(incomingURL.pathname)) {
		return isAuthenticatedMiddleware(request);
	}

	if (nonMiddlewarePaths.some((path) => incomingURL.pathname.startsWith(path))) {
		return NextResponse.next();
	}

	return authenticatedMiddleware(request);
}

export const config = {};
