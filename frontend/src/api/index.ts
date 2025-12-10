import { IPetsApi, petsApi } from './pets';
import { IReferenceApi, referenceApi } from './reference';
import { IUserApi, userApi } from './user';

export interface IApi extends IUserApi, IPetsApi, IReferenceApi {}

const api = (): IApi => ({
    ...userApi(),
    ...petsApi(),
    ...referenceApi(),
});

export default api();
