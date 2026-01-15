import { bloodRequestApi, IBloodRequestApi } from './bloodRequest';
import { IPetsApi, petsApi } from './pets';
import { IReferenceApi, referenceApi } from './reference';
import { IUserApi, userApi } from './user';

export interface IApi extends IUserApi, IPetsApi, IReferenceApi, IBloodRequestApi {}

const api = (): IApi => ({
    ...userApi(),
    ...petsApi(),
    ...referenceApi(),
    ...bloodRequestApi(),
});

export default api();
