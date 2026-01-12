import api from '../index';

export const getLivingConditions = async () => {
    try {
        const { status, data } = await api.getLivingConditions();

        if (status !== 200) {
            return null;
        }

        return data.data; // TODO errors?
    } catch (e) {
        return null;
    }
};
