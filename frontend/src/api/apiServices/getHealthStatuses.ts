import api from '../index';

export const getHealthStatuses = async () => {
    try {
        const { status, data } = await api.getHealthStatuses();

        if (status !== 200) {
            return null;
        }

        return data.data; // TODO errors?
    } catch (e) {
        return null;
    }
};
