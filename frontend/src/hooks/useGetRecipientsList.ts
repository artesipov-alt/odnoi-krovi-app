import { useQuery } from '@tanstack/react-query';

import { getRecipientsList } from 'api/apiServices/getRecipientsList';

export const useGetRecipientsList = (id: string) =>
    useQuery({
        queryKey: ['recipientsList', id],
        queryFn: () => getRecipientsList(id),
        enabled: !!id,
        staleTime: 5 * 60 * 1000, // данные "свежие" 5 минут
        gcTime: 10 * 60 * 1000, // хранится в кэше 10 минут
    });
