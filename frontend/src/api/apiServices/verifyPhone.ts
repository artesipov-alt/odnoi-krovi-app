import api from '../index';
import { VerifyPhoneRequest } from '../user';

export const verifyPhone = async (params: VerifyPhoneRequest) => {
    try {
        const { status, data } = await api.verifyPhone(params);

        if (status !== 200) {
            return { error: data.message };
        }

        return { data };
    } catch (e: any) {
        return { error: e.response.data.Message || 'Что-то пошло не так' };
    }
};
