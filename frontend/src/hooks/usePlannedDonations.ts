import { useQuery } from '@tanstack/react-query';

import { getPlannedDonations } from '../api/apiServices/getPlannedDonations';

export const usePlannedDonations = (userId: string) =>
    useQuery({
        queryKey: ['plannedDonations', userId],
        queryFn: () => getPlannedDonations(userId),
        enabled: !!userId,
        staleTime: 5 * 60 * 1000, // данные "свежие" 5 минут
        gcTime: 10 * 60 * 1000, // хранится в кэше 10 минут
    });
