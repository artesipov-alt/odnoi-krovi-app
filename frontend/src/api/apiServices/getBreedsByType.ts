import api from '../index';
import { PetType } from '../types';

export const getBreedsByType = async (pet: PetType) => {
    const { status, data } = await api.getBreedsByType(pet);

    if (status !== 200) {
        throw new Error(`Ошибка сервера: ${status}`);
    }

    return data.data;
};
