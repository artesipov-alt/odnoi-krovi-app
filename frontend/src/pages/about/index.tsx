import ArrowDown from 'imgs/svg/arrowDown';
import BackAngularArrow from 'imgs/svg/backAngularArrow';
import StoryAbout from 'imgs/svg/storyAbout';
import TechnicalAbout from 'imgs/svg/technicalAbout';
import ThanksAbout from 'imgs/svg/thanksAbout';
import { FC } from 'react';
import { useNavigate } from 'react-router-dom';

import styles from './About.module.less';

const About: FC = () => {
    const navigate = useNavigate();

    return (
        <div className={`${styles.wrapper} ${styles.aboutPageWrapper}`}>
            <div className={styles.header}>
                <button type='button' className={styles.backButton} onClick={() => navigate('/', { replace: true })}>
                    <BackAngularArrow />
                </button>
                <div className={styles.title}>О приложении</div>
            </div>

            <div className={`${styles.panel} ${styles.aboutPagePanel}`}>
                <div className={styles.menu}>
                    <button type='button' className={styles.menuItem} onClick={() => navigate('/about/history')}>
                        <span className={styles.menuItemLeft}>
                            <span className={styles.menuItemIcon}>
                                <StoryAbout />
                            </span>
                            <span className={styles.menuItemText}>История проекта</span>
                        </span>
                        <span className={styles.menuItemChevron}>
                            <ArrowDown />
                        </span>
                    </button>

                    <button type='button' className={styles.menuItem} onClick={() => navigate('/about/thanks')}>
                        <span className={styles.menuItemLeft}>
                            <span className={styles.menuItemIcon}>
                                <ThanksAbout />
                            </span>
                            <span className={styles.menuItemText}>Благодарности</span>
                        </span>
                        <span className={styles.menuItemChevron}>
                            <ArrowDown />
                        </span>
                    </button>

                    <button type='button' className={styles.menuItem} onClick={() => navigate('/about/tech')}>
                        <span className={styles.menuItemLeft}>
                            <span className={styles.menuItemIcon}>
                                <TechnicalAbout />
                            </span>
                            <span className={styles.menuItemText}>Техническая информация</span>
                        </span>
                        <span className={styles.menuItemChevron}>
                            <ArrowDown />
                        </span>
                    </button>
                </div>
            </div>
        </div>
    );
};

export default About;
