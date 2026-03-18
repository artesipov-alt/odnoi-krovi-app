import { BloodSearchApplyRequest } from '../donor';
import api from '../index';

export const bloodSearchApply = async (params: BloodSearchApplyRequest) => {
    try {
        const { status, data } = await api.bloodSearchApply(params);

        if (status !== 200) {
            return null;
        }

        return { data }; // TODO errors?
    } catch (e: any) {
        return null;
    }
};
