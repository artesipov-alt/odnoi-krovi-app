export const generateRecipientMessage = (params: {
  donorName: string;
  donorBloodGroup: string;
}): string => {
  return `\nКонтакты хозяина донора - ${params.donorName} (${params.donorBloodGroup})\n\nНаправляем контакты хозяина донора - обсудите возможность донации.\nБудьте вежливы и доброжелательны в общении!\nЕсли не получится договориться о донации, можете продолжить поиск в приложении.\n`;
};

export const generateDonorMessage = (params: {
  recipientName: string;
  recipientBloodGroup: string;
  recipientVolume: number;
}): string => {
  return `На ваше предложение откликнулся реципиент - ${params.recipientName} (${params.recipientVolume} мл, группа ${params.recipientBloodGroup})\n\nХозяин реципиента получил Ваши контакты. Дождитесь, пока с Вами свяжутся, или напишите хозяину реципиента`;
};

export const generateDonationMessage = (params: { volume: number }): string => {
  return `Реципиент подтвердил донацию в ${params.volume} мл.\n\nДонация состоялась успешно! Спасибо за вашу помощь.`;
};

export const generateVCF = (name: string, phone: string): string => {
  return `BEGIN:VCARD\r\nVERSION:3.0\r\nFN:${name}\r\nTEL:${phone}\r\nEND:VCARD`;
};
