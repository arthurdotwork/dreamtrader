import { NextRequest, NextResponse } from 'next/server';

export async function POST(req: NextRequest){
	const res = new NextResponse(JSON.stringify({}));

	const requestBody: {accessToken: string, refreshToken: string} = JSON.parse(await req.text());

	res.cookies.set('accessToken', requestBody.accessToken, {
		maxAge: 60 * 60 * 24 * 7,
		httpOnly: true,
		sameSite: 'strict',
		secure: true
	})

	res.cookies.set('refreshToken', requestBody.refreshToken, {
		maxAge: 60 * 60 * 24 * 7,
		httpOnly: true,
		sameSite: 'strict',
		secure: true
	})

	return res;
}