import { RejectDonationRequest } from '../bloodRequest';
import api from '../index';

export const cancelDonation = async (params: RejectDonationRequest) => {
    try {
        const { data } = await api.cancelDonation(params);

        return { data }; // TODO errors?
    } catch (e: any) {
        return null;
    }
};
