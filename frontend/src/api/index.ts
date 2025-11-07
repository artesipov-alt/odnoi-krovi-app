import { IUserApi, userApi } from './user';

export interface IApi extends IUserApi {}

const api = (): IApi => ({
    ...userApi(),
});

export default api();
