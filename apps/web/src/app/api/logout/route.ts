import { NextResponse } from "next/server";

export async function POST(req: Request){
	const res = new NextResponse(JSON.stringify({}));
	res.cookies.delete('accessToken')
	res.cookies.delete('refreshToken')

	return res;
}