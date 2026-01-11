import api from '../index';
import { PetType } from '../types';

export const getBreedsByType = async (pet: PetType) => {
    try {
        const { status, data } = await api.getBreedsByType(pet);

        if (status !== 200) {
            return null;
        }

        return data.data; // TODO errors?
    } catch (e) {
        return null;
    }
};
