import api from '../index';

export const getReproductiveStatuses = async () => {
    const { status, data } = await api.getReproductiveStatuses();

    if (status !== 200) {
        throw new Error(`Ошибка сервера: ${status}`);
    }

    return data.data;
};
