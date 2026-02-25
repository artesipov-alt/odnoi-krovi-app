import { useQuery } from '@tanstack/react-query';

import { getUserById } from 'api/apiServices/getUserById';

export const useGetUserById = (id: string) =>
    useQuery({
        queryKey: ['userById', id],
        queryFn: () => getUserById(id),
        enabled: !!id,
        staleTime: 5 * 60 * 1000, // данные "свежие" 5 минут
        gcTime: 10 * 60 * 1000, // хранится в кэше 10 минут
    });
