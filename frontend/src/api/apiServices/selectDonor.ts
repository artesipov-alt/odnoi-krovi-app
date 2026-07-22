import { SelectDonorRequest } from '../bloodRequest';
import api from '../index';

export const selectDonor = async (params: SelectDonorRequest) => {
    try {
        const { status, data } = await api.selectDonor(params);

        if (status !== 200) {
            return null;
        }

        return { data }; // TODO errors?
    } catch (e: any) {
        return null;
    }
};
