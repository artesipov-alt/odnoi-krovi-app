import { FC, useLayoutEffect, useState } from 'react';
import { TelegramUser } from 'types';

import Layout from 'components/Layout';

import Donor from './Donor';
import Recipient from './Recipient';
import Search from './Recipient/Search';
import Start from './Start';

enum View {
    START = 'start',
    DONOR = 'donor',
    RECIPIENT = 'recipient',
    START_SEARCH = '#startSearch',
}

type Props = {
    user: TelegramUser;
};

const Adding: FC<Props> = ({ user }) => {
    const [view, setView] = useState<View>(View.START);
    const [petIdForSearch, setPetIdForSearch] = useState<string | undefined>();

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
                return <Recipient onBackToStart={onBackToStartClickHandler} userId={user.id} />;
            }
            case View.DONOR: {
                return <Donor onBackToStart={onBackToStartClickHandler} userId={user.id} />;
            }
            case View.START_SEARCH: {
                return <Search petId={petIdForSearch} userId={user.id} />;
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

    useLayoutEffect(() => {
        if (window.location.hash) {
            const [hash, id] = window.location.hash.split('_');

            if (hash === View.START_SEARCH) {
                setView(View.START_SEARCH);
                setPetIdForSearch(id);
            }
        }
    }, []);

    return <Layout>{renderContent()}</Layout>;
};

export default Adding;
