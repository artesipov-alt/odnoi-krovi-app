import { AxiosPromise } from 'axios';

import { instance } from './instance';

export type SigninRequest = {
    appInitData: string;
};

export type SigninResponse = {
    userId: string;
    expiresAt: string;
    tokenType: string;
    accessToken: string;
};

export type ExternalServiceRequest = {
    providerId: string;
    providerName: string;
};

export interface IAuthApi {
    signinMax(params: SigninRequest): AxiosPromise<SigninResponse>;
    signinTg(params: SigninRequest): AxiosPromise<SigninResponse>;
    externalService(params: ExternalServiceRequest): AxiosPromise<SigninResponse>;
}

export const AUTH_URL = '/v1/auth/signin';

export const authApi = (): IAuthApi => ({
    signinMax(params) {
        return instance.post(`${AUTH_URL}/max`, params);
    },
    signinTg(params) {
        return instance.post(`${AUTH_URL}/telegram`, params);
    },
    externalService(params) {
        return instance.post(`${AUTH_URL}/service`, params, {
            headers: { 'X-Internal-Key': 'api-frimdsj764nsmksnj8x77dsjdns' },
        });
    },
});
