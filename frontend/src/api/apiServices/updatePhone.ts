import api from '../index';
import { UpdatePhoneRequest } from '../user';

export const updatePhone = async (params: UpdatePhoneRequest) => {
    try {
        const { data } = await api.updatePhone(params);

        return { data };
    } catch (e: any) {
        return { error: e.response.data.Message || 'Что-то пошло не так' };
    }
};
