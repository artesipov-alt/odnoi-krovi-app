import api from '../index';
import { UpdateUserRequest } from '../user';

export const updateUser = async (params: UpdateUserRequest) => {
    try {
        const { status, data } = await api.updateUser(params);

        if (status !== 200) {
            return { error: data.message };
        }

        return { data }; // TODO errors?
    } catch (e: any) {
        return { error: e.response.data.message || 'Что-то пошло не так' };
    }
};
