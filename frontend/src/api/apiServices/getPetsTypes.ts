import api from '../index';

export const getPetsTypes = async () => {
    try {
        const { status, data } = await api.getPetsTypes();

        if (status !== 200) {
            return null;
        }

        return data.data; // TODO errors?
    } catch (e) {
        return null;
    }
};
