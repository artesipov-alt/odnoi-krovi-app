import { RejectDonationRequest } from '../bloodRequest';
import api from '../index';

export const rejectDonation = async (params: RejectDonationRequest) => {
    try {
        const { status, data } = await api.rejectDonation(params);

        if (status !== 200) {
            return null;
        }

        return { data }; // TODO errors?
    } catch (e: any) {
        return null;
    }
};
