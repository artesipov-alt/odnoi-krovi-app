import Button from '@mui/material/Button';
import cn from 'classnames';
import useBodyScrollLock from 'hooks/useBodyScrollLock';
import { useGetUserById } from 'hooks/useGetUserById';
import bonusBg from 'imgs/bonusBg.png';
import profilePhoto from 'imgs/profilePhoto.png';
import BackAngularArrow from 'imgs/svg/backAngularArrow';
import ChatBubble from 'imgs/svg/chatBubble';
import Edit from 'imgs/svg/edit';
import Mail from 'imgs/svg/mail';
import MainLogo from 'imgs/svg/mainLogo';
import Max from 'imgs/svg/max';
import Phone from 'imgs/svg/phone';
import PrioritySearch from 'imgs/svg/prioritySearch';
import Tg from 'imgs/svg/tg';
import Vk from 'imgs/svg/vk';
import { FC, useEffect, useState } from 'react';
import { useNavigate } from 'react-router';
import { toast } from 'react-toastify';

import { updateUser } from 'api/apiServices/updateUser';
import { queryClient } from 'api/queryClient';
import Curtain from 'components/Curtain';
import Layout from 'components/Layout';

import styles from './Profile.module.less';

type Props = {
    userId: string;
};

const socialRows = [
    { title: 'Telegram', value: '@superdaschale', type: 'telegram' as const },
    { title: 'MAX', value: 'id384843', type: 'max' as const, trailing: 'x' },
    { title: 'ВКонтакте', value: 'Привязать', type: 'vk' as const, isAction: true },
];

const emailRegexp = /^\w+([+.-]?\w+)*@\w+([.-]?\w+)*(\.\w+)+$/i;

const normalizePhone = (value: string) => {
    const trimmed = value.trim();
    const hasPlus = trimmed.startsWith('+');
    const digits = trimmed.replace(/\D/g, '');

    return hasPlus ? `+${digits}` : digits;
};

const Profile: FC<Props> = ({ userId }) => {
    const navigate = useNavigate();
    const { data: userData, isLoading } = useGetUserById(userId);
    const [isInvitePopupOpen, setIsInvitePopupOpen] = useState(false);
    const [isEditCurtainOpen, setIsEditCurtainOpen] = useState(false);
    const [isEditLoading, setIsEditLoading] = useState(false);

    const [editFullName, setEditFullName] = useState('');
    const [editPhone, setEditPhone] = useState('');
    const [editEmail, setEditEmail] = useState('');

    const [isEditFullNameFocused, setIsEditFullNameFocused] = useState(false);
    const [isEditEmailFocused, setIsEditEmailFocused] = useState(false);

    useBodyScrollLock(isInvitePopupOpen || isEditCurtainOpen);

    const userInitial = userData?.fullName?.charAt(0).toUpperCase() || '?';

    const isPhoneValid = !!editPhone.trim() && /^\+?\d{10,15}$/.test(normalizePhone(editPhone));
    const isEmailValid = !!editEmail.trim() && !!editEmail.trim().match(emailRegexp);

    useEffect(() => {
        if (!isEditCurtainOpen || !userData) {
            return;
        }

        setEditFullName(userData.fullName || '');
        setEditPhone(userData.phone || '');
        setEditEmail(userData.email || '');
    }, [isEditCurtainOpen, userData]);

    if (isLoading || !userData) {
        return null;
    }

    const onSaveProfileClickHandler = async () => {
        if (
            !editFullName.trim() ||
            !editPhone.trim() ||
            !editEmail.trim() ||
            !isPhoneValid ||
            !isEmailValid ||
            isEditLoading
        ) {
            return;
        }

        setIsEditLoading(true);

        const { error } = await updateUser({
            id: userId,
            fullName: editFullName,
            phone: editPhone,
            email: editEmail,
        });

        if (error) {
            toast.warn(error);
            setIsEditLoading(false);

            return;
        }

        await queryClient.invalidateQueries({ queryKey: ['userById', userId] });
        setIsEditCurtainOpen(false);
        setIsEditLoading(false);
    };

    const isEditSaveDisabled =
        isEditLoading ||
        !editFullName.trim() ||
        !editPhone.trim() ||
        !editEmail.trim() ||
        !isPhoneValid ||
        !isEmailValid;

    return (
        <Layout className={styles.layout}>
            <div className={styles.page}>
                <div className={styles.header}>
                    <button type='button' className={styles.backButton} onClick={() => navigate('/owner')}>
                        <BackAngularArrow />
                    </button>
                    <div className={styles.avatar}>{userInitial}</div>
                    <h1 className={styles.fullName}>{userData.fullName}</h1>
                </div>

                <div className={styles.content}>
                    <div className={styles.contactCard}>
                        <div className={styles.contactRow}>
                            <span className={styles.contactIcon}>
                                <Phone />
                            </span>
                            <span className={styles.contactValue}>{userData.phone || 'Не указан'}</span>
                        </div>
                        <div className={styles.contactRow}>
                            <span className={styles.contactIcon}>
                                <Mail />
                            </span>
                            <span className={styles.contactValue}>{userData.email || 'Не указан'}</span>
                        </div>
                        <button
                            type='button'
                            className={styles.editButton}
                            aria-label='Редактировать профиль'
                            onClick={() => setIsEditCurtainOpen(true)}
                        >
                            <Edit />
                        </button>
                    </div>

                    <div className={styles.bonusCard}>
                        <div className={styles.bonusText}>
                            Пригласи друзей
                            <br />и получи бонус
                        </div>
                        <Button
                            variant='contained'
                            className={styles.detailsButton}
                            onClick={() => setIsInvitePopupOpen(true)}
                        >
                            Подробнее
                        </Button>
                        <img src={profilePhoto} alt='Питомцы' className={styles.bonusImage} />
                    </div>

                    <div className={styles.infoButtons}>
                        <button type='button' className={styles.infoButton}>
                            <span className={styles.infoIcon}>
                                <ChatBubble />
                            </span>
                            <span className={styles.infoText}>У меня проблема</span>
                        </button>
                        <button type='button' className={styles.infoButton} onClick={() => navigate('/about')}>
                            <span className={cn(styles.infoIcon, styles.infoIcon_app)}>
                                <MainLogo />
                            </span>
                            <span className={styles.infoText}>О приложении</span>
                        </button>
                    </div>

                    <div className={styles.socials}>
                        {socialRows.map(({ title, value, type, trailing, isAction }) => (
                            <div key={title} className={styles.socialRow}>
                                <div className={styles.socialLeft}>
                                    <div className={cn(styles.socialIcon, styles[`socialIcon_${type}`])}>
                                        {type === 'telegram' && <Tg />}
                                        {type === 'max' && <Max />}
                                        {type === 'vk' && <Vk />}
                                    </div>
                                    <span className={styles.socialTitle}>{title}</span>
                                </div>
                                <div className={styles.socialRight}>
                                    {isAction ? (
                                        <button type='button' className={styles.bindButton}>
                                            {value}
                                        </button>
                                    ) : (
                                        <>
                                            <span className={styles.socialValue}>{value}</span>
                                            {trailing && <span className={styles.trailing}>{trailing}</span>}
                                        </>
                                    )}
                                </div>
                            </div>
                        ))}
                    </div>
                </div>
            </div>
            {isInvitePopupOpen && (
                <Curtain
                    title={<span style={{ display: 'none' }} />}
                    onClose={() => setIsInvitePopupOpen(false)}
                    shouldCloseByWrapperClick
                    noRednerButtons
                    backgroundImage={bonusBg}
                    contentBorderRadius={0}
                    contentOverflow='visible'
                >
                    <div className={styles.popupShell}>
                        <div className={styles.popupIcon}>
                            <PrioritySearch />
                        </div>
                        <div className={styles.popupCard}>
                            <h3 className={styles.popupTitle}>Помогайте вместе!</h3>
                            <p className={styles.popupText}>
                                Пригласите друга в приложение - когда он проведет донацию, вы оба получите приоритетный
                                поиск
                            </p>
                            <div className={styles.popupDivider} />
                            <button type='button' className={styles.popupShareButton}>
                                Поделиться
                            </button>
                            <div className={styles.popupSocials}>
                                <button type='button' className={styles.popupSocialButton} aria-label='Telegram'>
                                    <Tg />
                                </button>
                                <button type='button' className={styles.popupSocialButton} aria-label='MAX'>
                                    <Max />
                                </button>
                            </div>
                            <button
                                type='button'
                                className={styles.popupBackButton}
                                onClick={() => setIsInvitePopupOpen(false)}
                            >
                                Вернуться
                            </button>
                        </div>
                    </div>
                </Curtain>
            )}

            {isEditCurtainOpen && (
                <Curtain
                    title={<span style={{ display: 'none' }} />}
                    onClose={() => setIsEditCurtainOpen(false)}
                    shouldCloseByWrapperClick
                    noRednerButtons
                    contentBorderRadius={0}
                >
                    <div className={styles.editCurtainWrapper}>
                        <div className={styles.editCurtainHeader}>
                            <div className={styles.editAvatarCircle}>{userInitial}</div>
                            <div className={styles.editHeaderIcon}>
                                <Edit />
                            </div>
                        </div>

                        <form
                            className={styles.editForm}
                            onSubmit={(e) => {
                                e.preventDefault();
                                onSaveProfileClickHandler();
                            }}
                        >
                            <label className={styles.editLabel} htmlFor='editFullName'>
                                Как вас зовут?
                            </label>
                            <div className={styles.editInputWrapper}>
                                <input
                                    id='editFullName'
                                    className={styles.editInput}
                                    value={editFullName}
                                    onChange={(e) => setEditFullName(e.target.value)}
                                    onFocus={() => setIsEditFullNameFocused(true)}
                                    onBlur={() => setIsEditFullNameFocused(false)}
                                />
                                {isEditFullNameFocused && (
                                    <button
                                        type='button'
                                        className={styles.editClearButton}
                                        aria-label='Очистить поле'
                                        onMouseDown={(e) => e.preventDefault()}
                                        onClick={() => setEditFullName('')}
                                    >
                                        ×
                                    </button>
                                )}
                            </div>

                            <label className={styles.editLabel} htmlFor='editPhone'>
                                Телефон
                            </label>
                            <input
                                id='editPhone'
                                className={cn(styles.editInput, {
                                    [styles.editInputError]: editPhone.trim() && !isPhoneValid,
                                })}
                                value={editPhone}
                                onChange={(e) => setEditPhone(e.target.value)}
                            />

                            <label className={styles.editLabel} htmlFor='editEmail'>
                                E-mail
                            </label>
                            <div className={styles.editInputWrapper}>
                                <input
                                    id='editEmail'
                                    className={cn(styles.editInput, {
                                        [styles.editInputError]: editEmail.trim() && !isEmailValid,
                                    })}
                                    value={editEmail}
                                    onChange={(e) => setEditEmail(e.target.value)}
                                    onFocus={() => setIsEditEmailFocused(true)}
                                    onBlur={() => setIsEditEmailFocused(false)}
                                />
                                {isEditEmailFocused && (
                                    <button
                                        type='button'
                                        className={styles.editClearButton}
                                        aria-label='Очистить поле'
                                        onMouseDown={(e) => e.preventDefault()}
                                        onClick={() => setEditEmail('')}
                                    >
                                        ×
                                    </button>
                                )}
                            </div>

                            <button type='submit' className={styles.editSaveButton} disabled={isEditSaveDisabled}>
                                Сохранить изменения
                            </button>
                        </form>
                    </div>
                </Curtain>
            )}
        </Layout>
    );
};

export default Profile;
