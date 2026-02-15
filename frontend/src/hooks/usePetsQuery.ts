import { useQuery } from '@tanstack/react-query';

import { getPets } from 'api/apiServices/getPets';

export const usePetsQuery = (userId: string) =>
    useQuery({
        queryKey: ['pets', userId],
        queryFn: () => getPets(userId),
        enabled: !!userId,
        staleTime: 5 * 60 * 1000, // данные "свежие" 5 минут
        gcTime: 10 * 60 * 1000, // хранится в кэше 10 минут
    });
