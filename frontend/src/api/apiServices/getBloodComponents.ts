import api from '../index';

export const getBloodComponents = async () => {
    try {
        const { status, data } = await api.getBloodComponents();

        if (status !== 200) {
            return null;
        }

        return data.data; // TODO errors?
    } catch (e) {
        return null;
    }
};
