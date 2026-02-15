import api from '../index';

export const getGenders = async () => {
    const { status, data } = await api.getGenders();

    if (status !== 200) {
        throw new Error(`Ошибка сервера: ${status}`);
    }

    return data.data;
};
