import api from '../index';
import { CreatePetRequest } from '../pets';

type Args = CreatePetRequest & {
    photo: File | null;
};

export const createDonor = async ({ photo, ...params }: Args) => {
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

        return { success: true };
    } catch (e) {
        return { success: false };
    }
};
