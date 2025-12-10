import api from '../index';

export const getPets = async (id: string) => {
    try {
        const { status, data } = await api.getPets(id);

        if (status !== 200) {
            return null;
        }

        return data; // TODO errors?
    } catch (e) {
        return null;
    }
};
