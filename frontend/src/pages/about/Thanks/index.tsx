import mainAboutBg from 'imgs/mainAboutBg.png';
import BackAngularArrow from 'imgs/svg/backAngularArrow';
import Clinics from 'imgs/svg/clinics';
import Development from 'imgs/svg/development';
import Media from 'imgs/svg/media';
import Partners from 'imgs/svg/partners';
import { FC } from 'react';
import { useNavigate } from 'react-router-dom';

import Layout from 'components/Layout';

import aboutStyles from '../About.module.less';

const AboutThanks: FC = () => {
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
                    <div className={aboutStyles.title}>Благодарности</div>
                </div>

                <div className={`${aboutStyles.panel} ${aboutStyles.thanksPanel}`}>
                    <div className={aboutStyles.thanksContent}>
                        <div className={aboutStyles.thanksIntro}>
                            Благодарим партнеров и друзей, внесших свой вклад&nbsp;в&nbsp;создание Портала
                        </div>

                        <div className={aboutStyles.thanksSectionCard}>
                            <div className={aboutStyles.thanksSectionHeader}>
                                <span className={aboutStyles.thanksSectionIcon}>
                                    <Partners />
                                </span>
                                <span className={aboutStyles.thanksSectionTitle}>Партнеры</span>
                            </div>

                            <ul className={aboutStyles.thanksList}>
                                <li>Международная ассоциация «Лапотто Романьоло»</li>
                                <li>Благотворительный фонд «Благополучие животных»</li>
                                <li>АНО «Учено-кинологический центр «Собаки-помощники»</li>
                                <li>ГБУ «Мосприрода»</li>
                                <li>Проект «Кошки-не-птицы.рф»</li>
                                <li>ОПБФ «Счастливый лай»</li>
                                <li>АНО «ПРОДВИЖЕНИЕ»</li>
                            </ul>
                        </div>

                        <div className={aboutStyles.thanksSectionCard}>
                            <div className={aboutStyles.thanksSectionHeader}>
                                <span className={aboutStyles.thanksSectionIcon}>
                                    <Clinics />
                                </span>
                                <span className={aboutStyles.thanksSectionTitle}>Клиники</span>
                            </div>

                            <div className={aboutStyles.thanksClinics}>
                                <div>Национальная</div>
                                <div>Раденис</div>
                                <div>Ветеринарная Палата</div>
                                <div>Dr Pets</div>
                                <div>Белый Клык</div>
                                <div>Биоконтроль</div>
                                <div>Шанс Био</div>
                                <div>ВЕТГЕМ</div>
                                <div>Дуо Кор</div>
                            </div>
                        </div>

                        <div className={aboutStyles.thanksSectionCard}>
                            <div className={aboutStyles.thanksSectionHeader}>
                                <span className={aboutStyles.thanksSectionIcon}>
                                    <Media />
                                </span>
                                <span className={aboutStyles.thanksSectionTitle}>СМИ</span>
                            </div>
                            <div className={aboutStyles.thanksList}>
                                <div>Журнал «Питомцы»</div>
                                <div>Издание «Ветеринария и Жизнь»</div>
                                <div>Добро.media</div>
                            </div>
                        </div>

                        <div className={aboutStyles.thanksSectionCard}>
                            <div className={aboutStyles.thanksSectionHeader}>
                                <span className={aboutStyles.thanksSectionIcon}>
                                    <Development />
                                </span>
                                <span className={aboutStyles.thanksSectionTitle}>Разработка</span>
                            </div>

                            <div className={aboutStyles.thanksDevText}>
                                <div>Члены команды – Артем, Руслана, Сандра, Олег</div>
                                <div>Front-end – Артур</div>
                                <div>Back-end – Руслан</div>
                                <div>Дизайн приложения – Дарья</div>
                                <div>Дизайн информационных материалов - Елена</div>
                            </div>

                            <div className={aboutStyles.thanksDevSubtitle}>Тестирование:</div>
                            <div className={aboutStyles.thanksQaList}>
                                <span>Артур</span>
                                <span>Руслан</span>
                                <span>Дарья</span>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </Layout>
    );
};

export default AboutThanks;
