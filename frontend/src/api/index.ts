import { bloodRequestApi, IBloodRequestApi } from './bloodRequest';
import { IPetsApi, petsApi } from './pets';
import { IPhotoApi, photoApi } from './photo';
import { IReferenceApi, referenceApi } from './reference';
import { IUserApi, userApi } from './user';

export interface IApi extends IUserApi, IPetsApi, IPhotoApi, IReferenceApi, IBloodRequestApi {}

const api = (): IApi => ({
    ...userApi(),
    ...petsApi(),
    ...photoApi(),
    ...referenceApi(),
    ...bloodRequestApi(),
});

export default api();
