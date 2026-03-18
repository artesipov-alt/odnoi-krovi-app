import api from '../index';

export const getRecipientsList = async (id: string) => {
    const { status, data } = await api.getRecipientsList(id);

    if (status !== 200) {
        throw new Error(`Ошибка сервера: ${status}`);
    }

    return data;
};
