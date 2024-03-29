import axios, { HeadersDefaults, AxiosHeaderValue, AxiosRequestConfig } from 'axios';
import camelcaseKeys from 'camelcase-keys';
import snakecaseKeys from 'snakecase-keys';

type APIResponseType = object | object[] | string | boolean;

type APIResponse<T extends APIResponseType, K extends string> = {
	[key in K]: T;
};

const axiosClient = axios.create();

const BASE_API_URL = process.env.NEXT_PUBLIC_API_URL;

const getRefreshToken = () => {
	return localStorage.getItem('refreshToken');
};

const getAccessToken = () => {
	return localStorage.getItem('accessToken');
};

const setTokens = async (tokens: { accessToken: string; refreshToken: string } | null) => {
	if (!tokens) {
		await axios.post('/api/logout');
		localStorage.removeItem('accessToken');
		localStorage.removeItem('refreshToken');
		return;
	}

	await axios.post('/api/auth', tokens);
	localStorage.setItem('accessToken', tokens.accessToken);
	localStorage.setItem('refreshToken', tokens.refreshToken);
};

// Replace this with our own backend base URL
axiosClient.defaults.baseURL = BASE_API_URL;

type headers = {
	'Content-Type': string;
	Accept: string;
	Authorization: string;
};

axiosClient.defaults.headers = {
	'Content-Type': 'application/json',
	Accept: 'application/json',
} as headers & HeadersDefaults & { [key: string]: AxiosHeaderValue };

axiosClient.interceptors.request.use(
	(config) => {
		const token = getAccessToken();
		if (token) {
			config.headers['Authorization'] = `Bearer ${token}`;
		}

		if (config.data) {
			config.data = snakecaseKeys(config.data, { deep: true });
		}

		return config;
	},
	(error) => {
		return Promise.reject(error);
	},
);

axiosClient.interceptors.response.use(
	(res) => {
		return {
			...res,
			data: camelcaseKeys(res.data, { deep: true }),
		};
	},
	async (err) => {
		const originalConfig = err.config;

		if (err.response) {
			if (err.response.status === 401 && !originalConfig._retry) {
				originalConfig._retry = true;

				try {
					const refreshToken = getRefreshToken();
					console.log({refreshToken})
					const { data } = await axios.post(`${BASE_API_URL}/api/v1/auth/refresh`, null, {
						headers: {
							Authorization: `Bearer ${refreshToken}`,
						},
					});

					await setTokens(camelcaseKeys(data, { deep: true }));

					originalConfig.headers['Authorization'] = `Bearer ${getAccessToken()}`;

					return axiosClient(originalConfig);
				} catch (_error) {
					await setTokens(null);
					window.location.href = window.location.origin;
					return Promise.reject(_error);
				}
			}
		}

		return Promise.reject(err);
	},
);

const get = async <T extends APIResponseType>(
	url: string,
	options?: AxiosRequestConfig,
): Promise<APIResponse<T, string>> => {
	const { data } = await axiosClient.get<APIResponse<T, string>>(url, { ...options });
	return data;
};

const post = async <T extends APIResponseType>(
	url: string,
	body: object,
	options?: AxiosRequestConfig,
): Promise<T> => {
	const { data } = await axiosClient.post<T>(url, body, { ...options });
	return data;
};

const put = async <T extends APIResponseType>(
	url: string,
	body: object,
	options?: AxiosRequestConfig,
): Promise<APIResponse<T, string>> => {
	const { data } = await axiosClient.put<APIResponse<T, string>>(url, body, { ...options });
	return data;
};

const del = async <T extends APIResponseType>(
	url: string,
	options?: AxiosRequestConfig,
): Promise<APIResponse<T, string>> => {
	const { data } = await axiosClient.delete<APIResponse<T, string>>(url, { ...options });
	return data;
};

export const client = { get, post, put, delete: del };
