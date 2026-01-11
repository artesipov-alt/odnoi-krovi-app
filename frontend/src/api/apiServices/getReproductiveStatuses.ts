import api from '../index';

export const getReproductiveStatuses = async () => {
    try {
        const { status, data } = await api.getReproductiveStatuses();

        if (status !== 200) {
            return null;
        }

        return data.data; // TODO errors?
    } catch (e) {
        return null;
    }
};
