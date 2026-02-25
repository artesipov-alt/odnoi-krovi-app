import { useQuery } from '@tanstack/react-query';

import { getPoolRequest } from 'api/apiServices/getPoolRequest';

export const useGetPoolRequestByPetId = (id?: string) =>
    useQuery({
        queryKey: ['poolRequestByPetId', id],
        queryFn: () => getPoolRequest(id),
        enabled: !!id,
        staleTime: 5 * 60 * 1000, // данные "свежие" 5 минут
        gcTime: 10 * 60 * 1000, // хранится в кэше 10 минут
    });
