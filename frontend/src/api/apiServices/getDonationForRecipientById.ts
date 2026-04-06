import api from '../index';

export const getDonationForRecipientById = async (id: string) => {
    try {
        const { status, data } = await api.getDonationForRecipientById(id);

        if (status !== 200) {
            return null;
        }

        return { data }; // TODO errors?
    } catch (e: any) {
        return null;
    }
};
