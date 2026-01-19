import { AddToPoolRequest } from '../bloodRequest';
import api from '../index';
import { CreatePetRequest } from '../pets';
import { Role } from '../user';

type Args = CreatePetRequest & {
    photo: File | null;
    poolInfo: Omit<AddToPoolRequest, 'petId'>;
};

export const createPet = async ({ photo, poolInfo, ...params }: Args) => {
    try {
        debugger;
        const { data } = await api.createPet(params);

        if (photo) {
            const {
                data: { url, path },
            } = await api.getPhotoLink(data.id);

            await fetch(url, {
                method: 'PUT',
                headers: {
                    'Content-Type': photo.type,
                },
                body: photo,
            });

            await api.confirmUploadPhoto(path);
        }

        if (params.petStatus === Role.DONOR) {
            await api.addToPool({ ...poolInfo, petId: data.id });
        }

        return { success: true };
    } catch (e) {
        return { success: false };
    }
};
