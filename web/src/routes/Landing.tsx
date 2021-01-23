import React from "react";
import { LandingLayout } from "../components/layouts/LandingLayout";
import { Hero } from "../components/sections/Hero";

export const Landing = () => {
  return (
    <LandingLayout>
      <Hero
        title="Valkyrie"
        subtitle="A Chat App using NestJS, REST and React"
        image={`${process.env.PUBLIC_URL}/logo.png`}
        ctaText="Get Started"
        ctaLink="/register"
      />
    </LandingLayout>
  );
};
