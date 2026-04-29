import { ExternalServiceRequest } from '../auth';
import api from '../index';
import { instance } from '../instance';

export const signinExtServ = async (params: ExternalServiceRequest) => {
    try {
        const { status, data } = await api.externalService(params);

        if (status !== 200) {
            return null;
        }

        if (data.accessToken && data.tokenType) {
            instance.defaults.headers.common.Authorization = `${data.tokenType} ${data.accessToken}`;
        }

        return { data: { ...data } };
    } catch (e: any) {
        return null;
    }
};
