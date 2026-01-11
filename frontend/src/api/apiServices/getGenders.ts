import api from '../index';

export const getGenders = async () => {
    try {
        const { status, data } = await api.getGenders();

        if (status !== 200) {
            return null;
        }

        return data.data; // TODO errors?
    } catch (e) {
        return null;
    }
};
