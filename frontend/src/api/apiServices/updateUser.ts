import api from '../index';
import { UpdateUserRequest } from '../user';

const ERROR = 'Что-то пошло не так';

export const updateUser = async (params: UpdateUserRequest) => {
    try {
        const { status, data } = await api.updateUser(params);

        if (status !== 200) {
            return { error: data.message };
        }

        return { data };
    } catch (e: any) {
        if (e.response.status === 409) {
            return { error: e.response.data.Message || ERROR };
        }

        return { error: ERROR };
    }
};
