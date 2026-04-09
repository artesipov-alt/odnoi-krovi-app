import api from '../index';

export const cancelDonation = async (id: string) => {
    try {
        const { data } = await api.cancelDonation(id);

        return { data }; // TODO errors?
    } catch (e: any) {
        return null;
    }
};
