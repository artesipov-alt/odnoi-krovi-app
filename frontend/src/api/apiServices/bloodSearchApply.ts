import { BloodSearchApplyRequest } from '../donor';
import api from '../index';

export const bloodSearchApply = async (params: BloodSearchApplyRequest) => {
    try {
        const { data } = await api.bloodSearchApply(params);

        return { data }; // TODO errors?
    } catch (e: any) {
        return null;
    }
};
