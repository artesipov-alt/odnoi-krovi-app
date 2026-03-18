import api from '../index';

export const getRecipientDetail = async (id: string) => {
    try {
        const { status, data } = await api.getRecipientDetails(id);

        if (status !== 200) {
            return null;
        }

        return data; // TODO errors?
    } catch (e) {
        return null;
    }
};
