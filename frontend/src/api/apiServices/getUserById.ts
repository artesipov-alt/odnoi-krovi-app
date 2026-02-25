import api from '../index';

export const getUserById = async (id: string) => {
    const { status, data } = await api.getUser(id);

    if (status !== 200) {
        throw new Error(`Ошибка сервера: ${status}`);
    }

    return data;
};
