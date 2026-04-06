import { ConfirmDonationRequest } from '../bloodRequest';
import api from '../index';

export const confirmDonation = async (params: ConfirmDonationRequest) => {
    try {
        const { status, data } = await api.confirmDonation(params);

        if (status !== 200) {
            return null;
        }

        return { data }; // TODO errors?
    } catch (e: any) {
        return null;
    }
};
