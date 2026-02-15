import api from '../index';

export const getLivingConditions = async () => {
    const { status, data } = await api.getLivingConditions();

    if (status !== 200) {
        throw new Error(`Ошибка сервера: ${status}`);
    }

    return data.data;
};
