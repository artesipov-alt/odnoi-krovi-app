import cn from 'classnames';
import onboardingScreen1 from 'imgs/onboardingScreen1.png';
import onboardingScreen2 from 'imgs/onboardingScreen2.png';
import onboardingScreen3 from 'imgs/onboardingScreen3.png';
import BackAngularArrow from 'imgs/svg/backAngularArrow';
import { FC, useState } from 'react';

import Layout from 'components/Layout';

import styles from './PromoSlider.module.less';

const promoSlides = [
    {
        title: 'Найдите кровь\nдля своего питомца',
        description: 'Ищите среди запасов клиник и доноров',
        image: onboardingScreen1,
    },
    {
        title: 'Спасайте\nжизни',
        description: 'Зарегистрируйте питомца донором и помогите тем, кто нуждается в переливании',
        image: onboardingScreen2,
    },
    {
        title: 'Получайте\nнаграды',
        description: 'Для своего питомца-донора',
        image: onboardingScreen3,
    },
];

const PromoSlider: FC = () => {
    const [activeSlideIndex, setActiveSlideIndex] = useState(0);

    const handleNextSlide = () => {
        setActiveSlideIndex((prev) => (prev + 1) % promoSlides.length);
    };

    const currentSlide = promoSlides[activeSlideIndex];

    return (
        <Layout className={styles.promoSliderCard}>
            <div className={styles.promoProgress}>
                {promoSlides.map((slide, index) => (
                    <span key={slide.title} className={styles.promoProgressTrack}>
                        <span
                            className={cn(styles.promoProgressFill, {
                                [styles.promoProgressFill_active]: index === activeSlideIndex,
                                [styles.promoProgressFill_done]: index < activeSlideIndex,
                            })}
                        />
                    </span>
                ))}
            </div>
            <button
                type='button'
                className={styles.promoNextButton}
                onClick={handleNextSlide}
                aria-label='Следующий слайд'
            >
                <BackAngularArrow />
            </button>
            <div className={styles.promoTitle}>
                {currentSlide.title.split('\n').map((line) => (
                    <span key={line}>
                        {line}
                        <br />
                    </span>
                ))}
            </div>
            <div className={styles.promoDescription}>{currentSlide.description}</div>
            <img src={currentSlide.image} alt='Промо слайд' className={styles.promoImage} />
        </Layout>
    );
};

export default PromoSlider;
