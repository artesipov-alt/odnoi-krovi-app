import api from '../index';
import { CreatePetRequest } from '../pets';
import { addPhoto } from './addPhoto';

type Args = CreatePetRequest & {
    photo: File | null;
    photoBuffer?: ArrayBuffer;
};

export const createDonor = async ({ photo, photoBuffer, ...params }: Args) => {
    try {
        const { data } = await api.createPet(params);

        if (!photo) {
            return { success: true };
        }

        const { success } = await addPhoto({ photo, buffer: photoBuffer, id: data.id, isAvatar: true });

        return { success };
    } catch (e) {
        return { success: false };
    }
};
