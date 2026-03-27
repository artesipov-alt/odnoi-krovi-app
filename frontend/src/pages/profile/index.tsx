import Button from '@mui/material/Button';
import cn from 'classnames';
import useBodyScrollLock from 'hooks/useBodyScrollLock';
import { useGetUserById } from 'hooks/useGetUserById';
import bonusBg from 'imgs/bonusBg.png';
import profileBonus from 'imgs/profileBonus.png';
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
import { ChangeEvent, FC, useEffect, useRef, useState } from 'react';
import InputMask from 'react-input-mask';
import { useNavigate } from 'react-router';
import { toast } from 'react-toastify';

import { addPhoto } from 'api/apiServices/addPhoto';
import { updateUser } from 'api/apiServices/updateUser';
import { queryClient } from 'api/queryClient';
import Curtain from 'components/Curtain';
import Layout from 'components/Layout';

// import PromoSlider from 'components/PromoSlider';
import styles from './Profile.module.less';

type Props = {
    userId: string;
};

type SocialRow = {
    title: string;
    value: string;
    type: 'telegram' | 'max';
    isAction?: boolean;
};

const socialRows: SocialRow[] = [
    { title: 'Telegram', value: '@superdaschale', type: 'telegram' as const },
    { title: 'MAX', value: 'id384843', type: 'max' as const },
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
    const [isAvatarUploading, setIsAvatarUploading] = useState(false);
    const avatarInputRef = useRef<HTMLInputElement | null>(null);

    const [editFullName, setEditFullName] = useState('');
    const [editPhone, setEditPhone] = useState('');
    const [editEmail, setEditEmail] = useState('');

    const [isEditFullNameFocused, setIsEditFullNameFocused] = useState(false);
    const [isEditEmailFocused, setIsEditEmailFocused] = useState(false);

    useBodyScrollLock(isInvitePopupOpen || isEditCurtainOpen);

    const userInitial = userData?.fullName?.charAt(0).toUpperCase() || '?';
    const avatarUrl = userData?.photoUrls?.[0];

    const isPhoneValid = !!editPhone.trim() && /^\+?\d{11,15}$/.test(normalizePhone(editPhone));
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

    const onEditAvatarClickHandler = () => {
        if (isAvatarUploading) {
            return;
        }

        avatarInputRef.current?.click();
    };

    const onAvatarSelectHandler = async (e: ChangeEvent<HTMLInputElement>) => {
        const newPhoto = e.target.files?.[0];

        if (!newPhoto || isAvatarUploading) {
            return;
        }

        setIsAvatarUploading(true);

        const { success } = await addPhoto({ id: userId, photo: newPhoto, isUserAvatar: true });

        if (!success) {
            toast.warn('Не удалось обновить фотографию, попробуйте еще раз');
            setIsAvatarUploading(false);
            e.target.value = '';

            return;
        }

        await queryClient.invalidateQueries({ queryKey: ['userById', userId] });
        setIsAvatarUploading(false);
        e.target.value = '';
    };

    return (
        <Layout>
            <div className={styles.page}>
                <div className={styles.header}>
                    <button type='button' className={styles.backButton} onClick={() => navigate('/owner')}>
                        <BackAngularArrow />
                    </button>
                    <div className={styles.avatar}>
                        {avatarUrl ? (
                            <img src={avatarUrl} alt='Фото профиля' className={styles.avatarImage} />
                        ) : (
                            userInitial
                        )}
                    </div>
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

                    <div className={styles.bonusCardNew}>
                        <div className={styles.bonusCardNewTitle}>
                            Спасайте жизни
                            <br />
                            вместе
                        </div>
                        <Button
                            variant='contained'
                            className={styles.bonusCardNewButton}
                            onClick={() => setIsInvitePopupOpen(true)}
                        >
                            Пригласить друга
                        </Button>
                        <img src={profileBonus} alt='Питомцы-доноры' className={styles.bonusCardNewImage} />
                    </div>

                    {/* <PromoSlider /> */}

                    <div className={styles.infoButtons}>
                        <button type='button' className={cn(styles.infoButton, styles.infoButton_disabled)} disabled>
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
                        {socialRows.map(({ title, value, type, isAction }) => (
                            <div key={title} className={styles.socialRow}>
                                <div className={styles.socialLeft}>
                                    <div className={cn(styles.socialIcon, styles[`socialIcon_${type}`])}>
                                        {type === 'telegram' && <Tg />}
                                        {type === 'max' && <Max />}
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
                            <div className={styles.editAvatarCircle}>
                                {avatarUrl ? (
                                    <img src={avatarUrl} alt='Фото профиля' className={styles.editAvatarImage} />
                                ) : (
                                    userInitial
                                )}
                                <button
                                    type='button'
                                    className={styles.editHeaderIcon}
                                    onClick={onEditAvatarClickHandler}
                                    aria-label='Изменить фото профиля'
                                >
                                    <Edit />
                                </button>
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
                            <InputMask
                                mask='+79999999999'
                                value={editPhone}
                                onChange={(e) => setEditPhone(normalizePhone(e.target.value))}
                            >
                                {(inputProps) => (
                                    <input
                                        {...inputProps}
                                        id='editPhone'
                                        className={cn(styles.editInput, {
                                            [styles.editInputError]: editPhone.trim() && !isPhoneValid,
                                        })}
                                    />
                                )}
                            </InputMask>

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
            <input
                ref={avatarInputRef}
                type='file'
                accept='image/*'
                style={{ display: 'none' }}
                onChange={onAvatarSelectHandler}
            />
        </Layout>
    );
};

export default Profile;
