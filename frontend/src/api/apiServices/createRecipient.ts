import { AddToPoolRequest } from '../bloodRequest';
import api from '../index';
import { CreatePetRequest } from '../pets';

type Args = CreatePetRequest & {
    photo: File | null;
    poolInfo: Omit<AddToPoolRequest, 'petId'>;
};

export const createRecipient = async ({ photo, poolInfo, ...params }: Args) => {
    try {
        const { data } = await api.createPet(params);

        if (photo) {
            const { data: photoLink } = await api.getPhotoLink({ id: data.id, photos_count: 1, for_pet_avatar: true });

            await fetch(photoLink.items[0].url, {
                method: 'PUT',
                headers: {
                    'Content-Type': photo.type,
                },
                body: photo,
            });

            await api.confirmUploadPhoto({ entityId: data.id, paths: [photoLink.items[0].path] });
        }

        await api.addToPool({ ...poolInfo, petId: data.id });

        return { success: true };
    } catch (e) {
        return { success: false };
    }
};
