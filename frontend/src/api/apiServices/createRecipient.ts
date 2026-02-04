import { AddToPoolRequest } from '../bloodRequest';
import api from '../index';
import { CreatePetRequest } from '../pets';
import { addPhoto } from './addPhoto';

type Args = CreatePetRequest & {
    photo: File | null;
    poolInfo: Omit<AddToPoolRequest, 'petId'>;
};

export const createRecipient = async ({ photo, poolInfo, ...params }: Args) => {
    try {
        const { data } = await api.createPet(params);

        await api.addToPool({ ...poolInfo, petId: data.id });

        if (!photo) {
            return { success: true };
        }

        const { success } = await addPhoto({ photo, id: data.id, isAvatar: true });

        return { success };
    } catch (e) {
        return { success: false };
    }
};
