import aboutBg from 'imgs/aboutBg.png';
import historyAvatar from 'imgs/historyAvatar.jpg';
import BackAngularArrow from 'imgs/svg/backAngularArrow';
import Quotes from 'imgs/svg/quotes';
import { FC } from 'react';
import { useNavigate } from 'react-router-dom';

import Layout from 'components/Layout';

import aboutStyles from '../About.module.less';

const AboutHistory: FC = () => {
    const navigate = useNavigate();

    return (
        <Layout>
            <div className={aboutStyles.wrapper} style={{ backgroundImage: `url(${aboutBg})` }}>
                <div className={aboutStyles.header}>
                    <button
                        type='button'
                        className={aboutStyles.backButton}
                        onClick={() => navigate('/about', { replace: true })}
                    >
                        <BackAngularArrow />
                    </button>
                    <div className={aboutStyles.title}>История проекта</div>
                </div>

                <div className={aboutStyles.panel}>
                    <div className={aboutStyles.content}>
                        <div className={aboutStyles.historyCard}>
                            <div className={aboutStyles.historyBigQuotes}>
                                <Quotes />
                            </div>
                            <img className={aboutStyles.historyAvatar} src={historyAvatar} alt='avatar' />
                            <div className={aboutStyles.historyMeta}>Основатель проекта, Ростислав Акадлович</div>
                            <div className={aboutStyles.historyLead}>
                                “Идея данного приложения родилась, когда мне пришлось искать кровь для своей кошки Рыси.
                                У Рыси была крупноклеточная лимфома, ей пророчили жить полгода, но мы боролись почти 2.
                            </div>
                            <div className={aboutStyles.historyMain}>
                                <div className={aboutStyles.historyParagraph}>
                                    Донорская кровь помогла Рысе прожить еще полгода, среди которых были счастливые
                                    моменты - а значит все было не зря.
                                </div>
                                <div className={aboutStyles.historyParagraph}>
                                    Надеемся, что Портал поможет хозяевам спасать своих любимых питомцев!”
                                </div>
                            </div>
                            <div className={aboutStyles.historyDivider} />
                            <div className={aboutStyles.historyFooter}>
                                Благодарим всех причастных доноров и их хозяев, а также клиники и банки крови, которые
                                помогали нам лечить Рысю – Белый Клык, Vetcity, Биоконтроль, Пастер, Вельвет, Сова, Свой
                                Доктор. Особая благодарность эндокринологу Белого Клыка Анне Юрьевне П. (Ч.), которая
                                заботилась о Рысе как о родной.
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </Layout>
    );
};

export default AboutHistory;
