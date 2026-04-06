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

        return { data: { ...data, userId: 'USR-A9fK3LQ2Xp' } };
        // return { data: { ...data, userId: 'USR-HdUVYbZA_3' } };
    } catch (e: any) {
        return null;
    }
};
