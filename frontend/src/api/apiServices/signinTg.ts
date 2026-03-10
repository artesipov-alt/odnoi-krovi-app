import { SigninRequest } from '../auth';
import api from '../index';
import { instance } from '../instance';

export const signinTg = async (params: SigninRequest) => {
    try {
        const { status, data } = await api.signinTg(params);

        if (status !== 200) {
            return null;
        }

        if (data.accessToken && data.tokenType) {
            instance.defaults.headers.common.Authorization = `${data.tokenType} ${data.accessToken}`;
        }

        return { data }; // TODO errors?
    } catch (e: any) {
        return null;
    }
};
