export type AuthenticationResponse = {
	accessToken: {
		token: string;
		expiresAt: Date;
	},
	refreshToken: {
		token: string;
		expiresAt: Date;
	},
};