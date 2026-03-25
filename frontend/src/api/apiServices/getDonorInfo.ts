import api from '../index';

export const getDonorInfo = async (id: string) => {
    try {
        const { status, data } = await api.getDonorInfo(id);

        if (status !== 200) {
            return null;
        }

        return { data }; // TODO errors?
    } catch (e: any) {
        return null;
    }
};
