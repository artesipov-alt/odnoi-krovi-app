import BackAngularArrow from 'imgs/svg/backAngularArrow';
import techBg from 'imgs/techBg.png';
import { FC } from 'react';
import { useNavigate } from 'react-router-dom';

import aboutStyles from '../About.module.less';

const AboutTech: FC = () => {
    const navigate = useNavigate();

    return (
        <div className={aboutStyles.wrapper} style={{ backgroundImage: `url(${techBg})` }}>
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
                        Идея данного приложения родилась, когда основателю проекта пришлось искать кровь для своей кошки
                        Рыси. У Рыси была крупноклеточная лимфома, ей пророчили жить полгода, но мы боролись почти 2.
                        <br />
                        Благодарим всех причастных доноров и их хозяев, а также клиники и банки крови, которые помогали
                        нам лечить Рысю – Белый Клык, Vetcity, Биоконтроль, Пастер, Вельвет, Сова, Свой Доктор.
                        <br />
                        Особая благодарность эндокринологу Белого Клыка Анне Юрьевне П. (Ч.), которая заботилась о Рысе
                        как о родной. Донорская кровь помогла Рысе прожить еще полгода, среди которых были счастливые
                        моменты, а значит все было не зря.
                        <br />
                        Надеемся, что Портал поможет хозяевам спасать своих любимых друзей
                    </div>
                </div>
            </div>
        </div>
    );
};

export default AboutTech;
