import api from '../index';

export const getPetsTypes = async () => {
    const { status, data } = await api.getPetsTypes();

    if (status !== 200) {
        throw new Error(`Ошибка сервера: ${status}`);
    }

    return data.data;
};
