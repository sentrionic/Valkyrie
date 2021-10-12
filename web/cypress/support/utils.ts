export const uuid = (): string => Cypress._.random(0, 1e6).toString();
