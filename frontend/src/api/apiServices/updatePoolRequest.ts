import { UpdatePoolRequestRequest } from '../bloodRequest';
import api from '../index';

export const updatePoolRequest = async (params: UpdatePoolRequestRequest) => {
    try {
        await api.updatePoolRequest(params);

        return { success: true };
    } catch (e) {
        return { success: false };
    }
};
