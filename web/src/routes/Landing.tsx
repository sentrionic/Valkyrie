import React from 'react';
import { LandingLayout } from '../components/layouts/LandingLayout';
import { Hero } from '../components/sections/Hero';

export const Landing: React.FC = () => (
  <LandingLayout>
    <Hero
      title="Valkyrie"
      subtitle="Whether you’re part of a school club,
        gaming group, worldwide art community,
        or just a handful of friends that want to spend time together,
        Valkyrie makes it easy to talk every day and hang out more often"
      image={`${process.env.PUBLIC_URL}/logo.png`}
      ctaText="Get Started"
      ctaLink="/register"
    />
  </LandingLayout>
);
