import api from '../index';

export const getUserByTelegramId = async (id: number) => {
    try {
        const { status, data } = await api.getUserByTelegramId(id);

        if (status !== 200) {
            return { error: data.message };
        }

        return { data }; // TODO errors?
    } catch (e: any) {
        return { error: e.response.data.message || 'Что-то пошло не так' };
    }
};
