import api from '../index';

export const getHealthStatuses = async () => {
    const { status, data } = await api.getHealthStatuses();

    if (status !== 200) {
        throw new Error(`Ошибка сервера: ${status}`);
    }

    return data.data;
};
