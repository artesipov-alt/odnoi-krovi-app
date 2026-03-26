import cn from 'classnames';
import onboardingScreen1 from 'imgs/onboardingScreen1.png';
import onboardingScreen2 from 'imgs/onboardingScreen2.png';
import onboardingScreen3 from 'imgs/onboardingScreen3.png';
import BackAngularArrow from 'imgs/svg/backAngularArrow';
import { FC, useEffect, useState } from 'react';
import { useNavigate } from 'react-router';

import styles from './PromoSlider.module.less';

const SLIDE_DURATION_MS = 5000;

const promoSlides = [
    {
        title: 'Найдите кровь\nдля своего питомца',
        description: 'Ищите среди запасов клиник и доноров',
        image: onboardingScreen1,
        route: '/owner',
    },
    {
        title: 'Спасайте\nжизни',
        description: 'Зарегистрируйте питомца донором и помогите тем, кто нуждается в переливании',
        image: onboardingScreen2,
        route: '/about/tech',
    },
    {
        title: 'Получайте\nнаграды',
        description: 'Для своего питомца-донора',
        image: onboardingScreen3,
        route: '/about/thanks',
    },
];

const PromoSlider: FC = () => {
    const navigate = useNavigate();
    const [activeSlideIndex, setActiveSlideIndex] = useState(0);

    useEffect(() => {
        const slideTimer = window.setTimeout(() => {
            setActiveSlideIndex((prev) => (prev + 1) % promoSlides.length);
        }, SLIDE_DURATION_MS);

        return () => {
            window.clearTimeout(slideTimer);
        };
    }, [activeSlideIndex]);

    const currentSlide = promoSlides[activeSlideIndex];

    return (
        <div className={styles.promoSliderCard}>
            <div className={styles.promoProgress}>
                {promoSlides.map((slide, index) => (
                    <span key={slide.title} className={styles.promoProgressTrack}>
                        <span
                            className={cn(styles.promoProgressFill, {
                                [styles.promoProgressFill_active]: index === activeSlideIndex,
                                [styles.promoProgressFill_done]: index < activeSlideIndex,
                            })}
                            style={
                                index === activeSlideIndex ? { animationDuration: `${SLIDE_DURATION_MS}ms` } : undefined
                            }
                        />
                    </span>
                ))}
            </div>
            <button
                type='button'
                className={styles.promoNextButton}
                onClick={() => navigate(currentSlide.route)}
                aria-label='Открыть раздел'
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
        </div>
    );
};

export default PromoSlider;
