import api from '../index';

export const getLocations = async () => {
    const { status, data } = await api.getLocations();

    if (status !== 200) {
        throw new Error(`Ошибка сервера: ${status}`);
    }

    return data.data;
};
