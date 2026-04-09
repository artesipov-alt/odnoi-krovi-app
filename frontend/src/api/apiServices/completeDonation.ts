import { CompleteDonationRequest } from '../donor';
import api from '../index';

export const completeDonation = async (params: CompleteDonationRequest) => {
    try {
        const { data } = await api.completeDonation(params);

        return { data }; // TODO errors?
    } catch (e: any) {
        return null;
    }
};
