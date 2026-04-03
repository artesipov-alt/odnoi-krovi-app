import Button from '@mui/material/Button';
import cn from 'classnames';
import useBodyScrollLock from 'hooks/useBodyScrollLock';
import { useGetUserById } from 'hooks/useGetUserById';
import bonusBg from 'imgs/bonusBg.png';
// import profileBonus from 'imgs/profileBonus.png';
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
import { getUserIdentities } from 'api/apiServices/getUserIdentities';
import { updateUser } from 'api/apiServices/updateUser';
import { queryClient } from 'api/queryClient';
import Curtain from 'components/Curtain';
import Layout from 'components/Layout';

// import PromoSlider from 'components/PromoSlider';
import styles from './Profile.module.less';

type Props = {
    userId: string;
};

type InvitePopupVariant = 'default' | 'bonusReceived';

type SocialRow = {
    title: string;
    value: string;
    type: 'telegram' | 'max';
    isAction?: boolean;
};

type UserIdentity = {
    providerName?: string;
    providerId?: string | number;
    refUrl?: string;
};

const MAX_BIND_BOT_URL = 'https://max.ru/c/-72684925241873/AZ0vv6wDGc8';
const TELEGRAM_BIND_BOT_URL = 'https://t.me/Odnakrovbot';
const REFERRAL_UTM_CAMPAIGN = 'help_together';
const REFERRAL_UTM_MEDIUM = 'referral';

const getProviderType = (providerName?: string): SocialRow['type'] | null => {
    const normalizedName = providerName?.trim().toLowerCase();

    if (!normalizedName) {
        return null;
    }

    if (normalizedName.includes('telegram') || normalizedName === 'tg') {
        return 'telegram';
    }

    if (normalizedName.includes('max')) {
        return 'max';
    }

    return null;
};

const getTelegramValue = (identity: UserIdentity) => {
    const providerId = identity.providerId?.toString().trim();

    if (providerId && /[a-z_]/i.test(providerId)) {
        return providerId.startsWith('@') ? providerId : `@${providerId}`;
    }

    if (providerId || identity.refUrl?.trim()) {
        return 'Привязан';
    }

    return null;
};

const getMaxValue = (identity: UserIdentity) => {
    if (identity.providerId !== undefined && identity.providerId !== null) {
        return 'Привязан';
    }

    return null;
};

const withReferralUtm = (url: string, type: SocialRow['type'], userId: string) => {
    try {
        const normalizedUrl = /^https?:\/\//i.test(url) ? url : `https://${url}`;
        const parsedUrl = new URL(normalizedUrl);

        parsedUrl.searchParams.set('utm_source', type);
        parsedUrl.searchParams.set('utm_medium', REFERRAL_UTM_MEDIUM);
        parsedUrl.searchParams.set('utm_campaign', REFERRAL_UTM_CAMPAIGN);
        parsedUrl.searchParams.set('utm_content', userId);

        return parsedUrl.toString();
    } catch {
        return url;
    }
};

const getTelegramReferralUrl = (identity: UserIdentity) => {
    const rawValue = identity.refUrl?.trim();

    if (!rawValue) {
        return null;
    }

    if (rawValue.startsWith('@')) {
        return `https://t.me/${rawValue.slice(1)}`;
    }

    if (/^(https?:\/\/)?(t\.me|telegram\.me)\//i.test(rawValue)) {
        return rawValue;
    }

    if (!rawValue.includes('/') && !rawValue.startsWith('http')) {
        return `https://t.me/${rawValue}`;
    }

    return null;
};

const getMaxReferralUrl = (identity: UserIdentity) => {
    const rawValue = identity.refUrl?.trim();

    if (!rawValue) {
        return MAX_BIND_BOT_URL;
    }

    return rawValue;
};

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
    const [invitePopupVariant, setInvitePopupVariant] = useState<InvitePopupVariant>('default');
    const [isEditCurtainOpen, setIsEditCurtainOpen] = useState(false);
    const [isEditLoading, setIsEditLoading] = useState(false);
    const [isAvatarUploading, setIsAvatarUploading] = useState(false);
    const [pendingAvatarFile, setPendingAvatarFile] = useState<File | null>(null);
    const [pendingAvatarPreviewUrl, setPendingAvatarPreviewUrl] = useState<string | null>(null);
    const [inviteIdentities, setInviteIdentities] = useState<UserIdentity[]>([]);
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
        setPendingAvatarFile(null);
        setPendingAvatarPreviewUrl(null);
    }, [isEditCurtainOpen, userData]);

    useEffect(() => {
        if (!pendingAvatarFile) {
            setPendingAvatarPreviewUrl(null);

            return;
        }

        const objectUrl = URL.createObjectURL(pendingAvatarFile);
        setPendingAvatarPreviewUrl(objectUrl);

        return () => {
            URL.revokeObjectURL(objectUrl);
        };
    }, [pendingAvatarFile]);

    useEffect(() => {
        if (!isInvitePopupOpen) {
            return;
        }

        let isCancelled = false;

        const fetchInviteIdentities = async () => {
            const response = await getUserIdentities(userId);

            if (isCancelled) {
                return;
            }

            if (response?.data?.identities?.length) {
                setInviteIdentities(response.data.identities);

                return;
            }

            setInviteIdentities(userData?.identities || []);
        };

        fetchInviteIdentities();

        return () => {
            isCancelled = true;
        };
    }, [isInvitePopupOpen, userId, userData?.identities]);

    if (isLoading || !userData) {
        return null;
    }
    const identitiesByType = (userData.identities || []).reduce<Partial<Record<SocialRow['type'], UserIdentity>>>(
        (acc, identity) => {
            const providerType = getProviderType(identity.providerName);

            if (!providerType || acc[providerType]) {
                return acc;
            }

            acc[providerType] = identity;

            return acc;
        },
        {},
    );

    const telegramValue = identitiesByType.telegram ? getTelegramValue(identitiesByType.telegram) : null;
    const maxValue = identitiesByType.max ? getMaxValue(identitiesByType.max) : null;

    const socialRows: SocialRow[] = [
        telegramValue
            ? { title: 'Telegram', value: telegramValue, type: 'telegram' }
            : { title: 'Telegram', value: 'Привязать', type: 'telegram', isAction: true },
        maxValue
            ? { title: 'MAX', value: maxValue, type: 'max' }
            : { title: 'MAX', value: 'Привязать', type: 'max', isAction: true },
    ];

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

        const hasProfileChanges =
            editFullName.trim() !== (userData.fullName || '').trim() ||
            normalizePhone(editPhone) !== normalizePhone(userData.phone || '') ||
            editEmail.trim() !== (userData.email || '').trim();
        const hasAvatarChanges = !!pendingAvatarFile;

        if (!hasProfileChanges && !hasAvatarChanges) {
            setIsEditCurtainOpen(false);

            return;
        }

        setIsEditLoading(true);
        setIsAvatarUploading(hasAvatarChanges);

        const [updateResult, photoResult] = await Promise.all([
            hasProfileChanges
                ? updateUser({
                      id: userId,
                      fullName: editFullName,
                      phone: editPhone,
                      email: editEmail,
                  })
                : Promise.resolve(null),
            hasAvatarChanges && pendingAvatarFile
                ? addPhoto({ id: userId, photo: pendingAvatarFile, isUserAvatar: true })
                : Promise.resolve(null),
        ]);

        if (updateResult?.error) {
            toast.warn(updateResult.error);
        }

        if (hasAvatarChanges && photoResult && !photoResult.success) {
            toast.warn('Не удалось обновить фотографию, попробуйте еще раз');
        }

        const hasUpdateSuccess = hasProfileChanges && !updateResult?.error;
        const hasPhotoSuccess = hasAvatarChanges && !!photoResult?.success;

        if (hasUpdateSuccess || hasPhotoSuccess) {
            await queryClient.invalidateQueries({ queryKey: ['userById', userId] });
        }

        if ((hasProfileChanges && updateResult?.error) || (hasAvatarChanges && photoResult && !photoResult.success)) {
            setIsAvatarUploading(false);
            setIsEditLoading(false);

            return;
        }

        setPendingAvatarFile(null);
        setPendingAvatarPreviewUrl(null);
        setIsEditCurtainOpen(false);
        setIsAvatarUploading(false);
        setIsEditLoading(false);
    };

    const isEditSaving = isEditLoading || isAvatarUploading;
    const isEditSaveDisabled =
        isEditSaving ||
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

    const onEditCurtainCloseHandler = () => {
        if (isAvatarUploading || isEditLoading) {
            return;
        }

        setIsEditCurtainOpen(false);
    };

    const onAvatarSelectHandler = (e: ChangeEvent<HTMLInputElement>) => {
        const newPhoto = e.target.files?.[0];

        if (!newPhoto || isAvatarUploading || isEditLoading) {
            return;
        }

        setPendingAvatarFile(newPhoto);
        e.target.value = '';
    };

    const onBindClickHandler = (type: SocialRow['type']) => {
        if (type === 'max') {
            window.open(MAX_BIND_BOT_URL, '_blank', 'noopener,noreferrer');

            return;
        }

        window.open(TELEGRAM_BIND_BOT_URL, '_blank', 'noopener,noreferrer');
    };

    const openInvitePopup = (variant: InvitePopupVariant) => {
        setInvitePopupVariant(variant);
        setIsInvitePopupOpen(true);
    };

    const getInviteShareUrl = (type: SocialRow['type']) => {
        const sourceIdentities = inviteIdentities.length ? inviteIdentities : userData.identities || [];
        const identity = sourceIdentities.find((item) => getProviderType(item.providerName) === type);

        if (!identity && type === 'telegram') {
            toast.warn('Telegram не привязан');

            return null;
        }

        const baseUrl =
            type === 'telegram' ? getTelegramReferralUrl(identity || {}) : getMaxReferralUrl(identity || {});

        if (!baseUrl) {
            toast.warn('Не удалось сформировать реферальную ссылку');

            return null;
        }

        return withReferralUtm(baseUrl, type, userId);
    };

    const onShareInviteClickHandler = async (type: SocialRow['type']) => {
        const shareUrl = getInviteShareUrl(type);

        if (!shareUrl) {
            return;
        }

        const shareText = `Присоединяйся к Одной Крови: ${shareUrl}`;

        if (navigator.share) {
            try {
                await navigator.share({
                    title: 'Приглашение в Одной Крови',
                    text: shareText,
                });

                return;
            } catch (error) {
                if (error instanceof DOMException && error.name === 'AbortError') {
                    return;
                }
            }
        }

        if (navigator.clipboard?.writeText) {
            await navigator.clipboard.writeText(shareText);
            toast.success('Ссылка скопирована');

            return;
        }

        window.open(shareUrl, '_blank', 'noopener,noreferrer');
    };

    const editAvatarUrl = pendingAvatarPreviewUrl || avatarUrl;

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
                            onClick={() => openInvitePopup('default')}
                        >
                            Подробнее
                        </Button>
                        <img src={profilePhoto} alt='Питомцы' className={styles.bonusImage} />
                    </div>

                    {/* <div className={styles.bonusCardNew}>
                        <div className={styles.bonusCardNewTitle}>
                            Спасайте жизни
                            <br />
                            вместе
                        </div>
                        <Button
                            variant='contained'
                            className={styles.bonusCardNewButton}
                            onClick={() => openInvitePopup('bonusReceived')}
                        >
                            Пригласить друга
                        </Button>
                        <img src={profileBonus} alt='Питомцы-доноры' className={styles.bonusCardNewImage} />
                    </div> */}

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
                                        <button
                                            type='button'
                                            className={styles.bindButton}
                                            onClick={() => onBindClickHandler(type)}
                                        >
                                            {value}
                                        </button>
                                    ) : (
                                        <span className={styles.socialValue}>{value}</span>
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
                                {invitePopupVariant === 'bonusReceived'
                                    ? 'Бонус уже получен, но можете пригласить больше\u00A0друзей\u00A0и\u00A0вместе спасать жизни'
                                    : 'Пригласите друга в приложение - когда он проведет донацию, вы оба получите приоритетный поиск'}
                            </p>
                            <div className={styles.popupDivider} />
                            <button
                                type='button'
                                className={styles.popupShareButton}
                                onClick={() => {
                                    onShareInviteClickHandler('telegram');
                                }}
                            >
                                Поделиться
                            </button>
                            <div className={styles.popupSocials}>
                                <button
                                    type='button'
                                    className={styles.popupSocialButton}
                                    aria-label='Telegram'
                                    onClick={() => {
                                        onShareInviteClickHandler('telegram');
                                    }}
                                >
                                    <Tg />
                                </button>
                                <button
                                    type='button'
                                    className={styles.popupSocialButton}
                                    aria-label='MAX'
                                    onClick={() => {
                                        onShareInviteClickHandler('max');
                                    }}
                                >
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
                    onClose={onEditCurtainCloseHandler}
                    shouldCloseByWrapperClick
                    noRednerButtons
                    contentBorderRadius={0}
                >
                    <div className={styles.editCurtainWrapper}>
                        <div className={styles.editCurtainHeader}>
                            <div className={styles.editAvatarCircle}>
                                {editAvatarUrl ? (
                                    <img src={editAvatarUrl} alt='Фото профиля' className={styles.editAvatarImage} />
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
                                {isEditSaving ? (
                                    <span className={styles.editSaveButtonContent}>
                                        <span className={styles.editSaveButtonLoader} />
                                        <span>Сохраняем...</span>
                                    </span>
                                ) : (
                                    'Сохранить изменения'
                                )}
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
