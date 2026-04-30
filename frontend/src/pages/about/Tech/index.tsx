import cn from 'classnames';
import mainAboutBg from 'imgs/mainAboutBg.png';
import BackAngularArrow from 'imgs/svg/backAngularArrow';
import { FC } from 'react';
import { useNavigate } from 'react-router-dom';

import Layout from 'components/Layout';

import aboutStyles from '../About.module.less';

const AboutTech: FC = () => {
    const navigate = useNavigate();

    return (
        <Layout>
            <div className={aboutStyles.wrapper} style={{ backgroundImage: `url(${mainAboutBg})` }}>
                <div className={aboutStyles.header}>
                    <button
                        type='button'
                        className={aboutStyles.backButton}
                        onClick={() => navigate('/about', { replace: true })}
                    >
                        <BackAngularArrow />
                    </button>
                    <div className={aboutStyles.title}>Техническая информация</div>
                </div>

                <div className={`${aboutStyles.panel} ${aboutStyles.techPagePanel}`}>
                    <div className={aboutStyles.techCard}>
                        <div className={aboutStyles.sectionText}>
                            {/* Идея данного приложения родилась, когда основателю проекта пришлось искать кровь для своей
                            кошки Рыси. У Рыси была крупноклеточная лимфома, ей пророчили жить полгода, но мы боролись
                            почти 2.
                            <br />
                            Благодарим всех причастных доноров и их хозяев, а также клиники и банки крови, которые
                            помогали нам лечить Рысю – Белый Клык, Vetcity, Биоконтроль, Пастер, Вельвет, Сова, Свой
                            Доктор.
                            <br />
                            Особая благодарность эндокринологу Белого Клыка Анне Юрьевне П. (Ч.), которая заботилась о
                            Рысе как о родной. Донорская кровь помогла Рысе прожить еще полгода, среди которых были
                            счастливые моменты, а значит все было не зря.
                            <br />
                            Надеемся, что Портал поможет хозяевам спасать своих любимых друзей */}
                        </div>
                        <div className={cn(aboutStyles.sectionText, { [aboutStyles.marginTop]: true })}>
                            Версия: 1.0.
                            <br />
                            <a
                                target='_blank'
                                rel='noreferrer'
                                className={aboutStyles.link}
                                href='https://однойкрови.рф/docs#n-a9dea2ae-b2a2-4bc6-b0ea-1eca71588ab0'
                            >
                                Пользовательское соглашение
                            </a>
                            <a
                                target='_blank'
                                rel='noreferrer'
                                className={aboutStyles.link}
                                href='https://однойкрови.рф/docs#n-80ae6549-954a-4c8c-bc54-357a93ec3dee'
                            >
                                Политика конфиденциальности
                            </a>
                            <br />© 2026. Одной Крови. Все права защищены.
                        </div>
                    </div>
                </div>
            </div>
        </Layout>
    );
};

export default AboutTech;
