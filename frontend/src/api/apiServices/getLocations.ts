import api from '../index';

export const getLocations = async () => {
    try {
        const { status, data } = await api.getLocations();

        if (status !== 200) {
            return null;
        }

        return data.data; // TODO errors?
    } catch (e) {
        return null;
    }
};
