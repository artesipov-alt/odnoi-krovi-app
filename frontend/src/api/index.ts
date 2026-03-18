import { authApi, IAuthApi } from './auth';
import { bloodRequestApi, IBloodRequestApi } from './bloodRequest';
import { donorApi, IDonorApi } from './donor';
import { IPetsApi, petsApi } from './pets';
import { IPhotoApi, photoApi } from './photo';
import { IReferenceApi, referenceApi } from './reference';
import { IUserApi, userApi } from './user';

export interface IApi extends IUserApi, IAuthApi, IPetsApi, IPhotoApi, IDonorApi, IReferenceApi, IBloodRequestApi {}

const api = (): IApi => ({
    ...userApi(),
    ...petsApi(),
    ...authApi(),
    ...photoApi(),
    ...donorApi(),
    ...referenceApi(),
    ...bloodRequestApi(),
});

export default api();
