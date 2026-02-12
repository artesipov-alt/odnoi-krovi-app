import api from '../index';

export const getBloodComponents = async () => {
    const { status, data } = await api.getBloodComponents();

    if (status !== 200) {
        throw new Error(`Ошибка сервера: ${status}`);
    }

    return data.data;
};
