import api from '../index';
import { GetUserContactsRequest } from '../user';

export const getUserContacts = async (params: GetUserContactsRequest) => {
    try {
        const { status, data } = await api.getUserContacts(params);

        if (status !== 200) {
            return null;
        }

        return { data }; // TODO errors?
    } catch (e: any) {
        return null;
    }
};
