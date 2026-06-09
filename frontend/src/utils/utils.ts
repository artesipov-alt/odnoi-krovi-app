import { format } from 'date-fns';
import { ru } from 'date-fns/locale';

import { Identities } from '../api/user';

export enum Variants {
    DAYS = 'days',
    YEARS = 'years',
    MONTHS = 'months',
}

const declension = {
    [Variants.DAYS]: ['день', 'дня', 'дней'],
    [Variants.YEARS]: ['год', 'года', 'лет'],
    [Variants.MONTHS]: ['месяц', 'месяца', 'месяцев'],
};

export const getCorrectDeclension = (type: Variants, number: number): string => {
    const lastDigit = number % 10;
    const lastTwoDigits = number % 100;

    if (lastTwoDigits >= 11 && lastTwoDigits <= 19) {
        return declension[type][2];
    }

    if (lastDigit === 1) {
        return declension[type][0];
    }

    if (lastDigit >= 2 && lastDigit <= 4) {
        return declension[type][1];
    }

    return declension[type][2];
};

export const getDateFormat = (date: Date) => format(date, 'dd.MM.yyyy', { locale: ru });

export const isExpiredDate = (expiresAt: string): boolean => {
    const expirationTime = new Date(expiresAt).getTime();
    const currentTime = Date.now();

    return currentTime > expirationTime;
};

export const isWithinHours = (startTime: string, hours: number): boolean => {
    const updateDate = new Date(startTime).getTime();
    const now = new Date().getTime();
    const diffInMs = now - updateDate;
    const thresholdInMs = hours * 60 * 60 * 1000;

    return diffInMs < thresholdInMs;
};

export const matchIdentities = (first: Identities[], second?: Identities[]) => {
    if (!second) {
        return first;
    }

    const result: Identities[] = [];

    first.forEach((item) => {
        if (second.some(({ providerName }) => providerName === item.providerName)) {
            result.push(item);
        }
    });

    return result;
};
