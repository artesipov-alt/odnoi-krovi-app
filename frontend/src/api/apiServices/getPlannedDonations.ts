import api from '../index';

export const getPlannedDonations = async (id: string) => {
    const { status, data } = await api.getPlannedDonations(id);

    if (status !== 200) {
        throw new Error(`Ошибка сервера: ${status}`);
    }

    return data.items;
};
