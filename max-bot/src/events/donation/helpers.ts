export const generateDonationMessage = (params: {
  initiator: 'donor' | 'recipient';
  volume: number;
}): string => {
  const who = params.initiator === 'donor' ? 'Донор' : 'Реципиент';
  return `${who} сообщил о донации в ${params.volume} мл.\n\nЧтобы подтвердить донацию, перейдите в приложение.`;
};
