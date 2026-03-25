import api from '../index';

export const getUserIdentities = async (id: string) => {
    try {
        const { status, data } = await api.getUserIdentities(id);

        if (status !== 200) {
            return null;
        }

        return { data }; // TODO errors?
    } catch (e: any) {
        return null;
    }
};
