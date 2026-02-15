import { format } from 'date-fns';
import { ru } from 'date-fns/locale';

export enum Variants {
    YEARS = 'years',
    MONTHS = 'months',
}

const declension = {
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
