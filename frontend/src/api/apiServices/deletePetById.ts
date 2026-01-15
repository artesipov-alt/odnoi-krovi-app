import api from '../index';

export const deletePetById = async (id: string) => {
    try {
        await api.deletePetById(id);

        return 'ok';
    } catch (e) {
        return null;
    }
};
