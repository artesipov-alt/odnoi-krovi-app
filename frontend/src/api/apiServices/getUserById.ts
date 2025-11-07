import api from '../index';

export const getUserById = async (id: number) => {
    try {
        const { status, data } = await api.getUser(id);

        if (status !== 200) {
            return null;
        }

        return data; // TODO errors?
    } catch (e) {
        return null;
    }
};
