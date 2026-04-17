import api from '../index';

export const getCompletedDonations = async (id: string) => {
    try {
        const { status, data } = await api.getCompletedDonations(id);

        if (status !== 200) {
            return null;
        }

        return { data }; // TODO errors?
    } catch (e: any) {
        return null;
    }
};
