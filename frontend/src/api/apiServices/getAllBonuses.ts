import api from '../index';

export const getAllBonuses = async (id: string) => {
    try {
        const { status, data } = await api.getAllBonuses(id);

        if (status !== 200) {
            return null;
        }

        return { data }; // TODO errors?
    } catch (e: any) {
        return null;
    }
};
