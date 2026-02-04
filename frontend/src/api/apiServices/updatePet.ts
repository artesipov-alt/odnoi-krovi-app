import api from '../index';
import { Pet } from '../pets';

export const updatePet = async (pet: Pet) => {
    try {
        await api.updatePet(pet);

        return { success: true };
    } catch (e) {
        return { success: false };
    }
};
