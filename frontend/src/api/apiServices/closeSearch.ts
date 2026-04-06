import api from '../index';

export const closeSearch = async (id: string) => {
    try {
        const { status, data } = await api.closeSearch(id);

        if (status !== 200) {
            return null;
        }

        return { data }; // TODO errors?
    } catch (e: any) {
        return null;
    }
};
