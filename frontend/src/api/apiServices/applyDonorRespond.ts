import api from '../index';

export const applyDonorRespond = async (id: string) => {
    try {
        const { status, data } = await api.applyDonorRespond(id);

        if (status !== 200) {
            return null;
        }

        return { data }; // TODO errors?
    } catch (e: any) {
        return null;
    }
};
