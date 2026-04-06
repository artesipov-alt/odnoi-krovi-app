import api from '../index';

export const rejectDonation = async (id: string) => {
    try {
        const { status, data } = await api.rejectDonation(id);

        if (status !== 200) {
            return null;
        }

        return { data }; // TODO errors?
    } catch (e: any) {
        return null;
    }
};
