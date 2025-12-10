import { FC, useState } from 'react';
import { TelegramUser } from 'types';

import Layout from 'components/Layout';

import Recipient from './Recipient';
import Start from './Start';

enum View {
    START = 'start',
    DONOR = 'donor',
    RECIPIENT = 'recipient',
}

type Props = {
    user: TelegramUser;
};

const Adding: FC<Props> = ({ user }) => {
    const [view, setView] = useState<View>(View.START);

    const onAddRecipientPetClickHandler = () => {
        setView(View.RECIPIENT);
    };

    const onAddDonorPetClickHandler = () => {
        setView(View.DONOR);
    };

    const onBackToStartClickHandler = () => {
        setView(View.START);
    };

    const renderContent = () => {
        switch (view) {
            case View.RECIPIENT: {
                return <Recipient onBackToStart={onBackToStartClickHandler} />;
            }
            default: {
                return (
                    <Start
                        onAddDonorClick={onAddDonorPetClickHandler}
                        onAddRecipientClick={onAddRecipientPetClickHandler}
                    />
                );
            }
        }
    };

    return <Layout>{renderContent()}</Layout>;
};

export default Adding;
