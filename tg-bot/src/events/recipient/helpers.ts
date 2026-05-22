export const generateRecipientMessage = (params: {
  donorName: string;
  donorBloodGroup: string;
}): string => {
  const donorName = params.donorName?.trim() || "Анонимный донор";
  const donorBloodGroup =
    params.donorBloodGroup === "UNKNOWN"
      ? "не определена"
      : params.donorBloodGroup;

  return `\nКонтакты хозяина донора - ${donorName} (группа ${donorBloodGroup})\n\nНаправляем контакты хозяина донора - обсудите возможность донации.\nБудьте вежливы и доброжелательны в общении!\nЕсли не получится договориться о донации, можете продолжить поиск в приложении.\n`;
};

export const generateDonorMessage = (params: {
  recipientName: string;
  recipientBloodGroup: string;
  recipientVolume: number;
}): string => {
  const recipientBloodGroup =
    params.recipientBloodGroup === "UNKNOWN"
      ? "не определена"
      : params.recipientBloodGroup;

  return `На ваше предложение откликнулся реципиент - ${params.recipientName} (${params.recipientVolume} мл, группа ${recipientBloodGroup})\n\nХозяин реципиента получил Ваши контакты. Дождитесь, пока с Вами свяжутся, или напишите хозяину реципиента`;
};

export const generateDonationMessage = (params: {
  volume: number;
  recipientPetName: string;
  recipientBloodGroup: string;
}): string => {
  const recipientBloodGroup =
    params.recipientBloodGroup === "UNKNOWN"
      ? "не определена"
      : params.recipientBloodGroup;

  return `Донация подтверждена (реципиент ${params.recipientPetName}, группа ${recipientBloodGroup}). Спасибо за Вашу помощь! Вам начислены бонусы – посмотрите их на Портале.`;
};
