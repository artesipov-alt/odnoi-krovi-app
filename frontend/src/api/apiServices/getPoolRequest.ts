import api from '../index';

export const getPoolRequest = async (id?: string) => {
    if (!id) {
        throw new Error('Нет id питомца');
    }

    const { status, data } = await api.getPoolRequest(id);

    if (status !== 200) {
        throw new Error(`Ошибка сервера: ${status}`);
    }

    return data;
};
