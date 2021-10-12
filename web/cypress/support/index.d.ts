// load type definitions that come with Cypress module
// eslint-disable-next-line spaced-comment
/// <reference types="cypress" />

declare namespace Cypress {
  interface Chainable {
    registerUser(email: string, id: string): Chainable<AUTWindow>;
    loginUser(email: string): Chainable<AUTWindow>;
    sendRequest(id: string): Chainable<AUTWindow>;
    createGuild(id: string): Chainable<AUTWindow>;
    joinGuild(invite: string): Chainable<AUTWindow>;
    clickOnFirstGuild(): Chainable<AUTWindow>;
    openGuildMenu(): Chainable<AUTWindow>;
    getChat(): Chainable<AUTWindow>;
    firstMessage(): Chainable<AUTWindow>;
  }
}
