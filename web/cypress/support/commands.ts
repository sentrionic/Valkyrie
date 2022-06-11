Cypress.Commands.add('registerUser', (email, id) => {
  cy.visit('/');
  cy.contains('Valkyrie');
  cy.contains('Register').click();

  cy.url().should('include', '/register');
  cy.get('input[name="email"]').type(email).should('have.value', email);
  cy.get('input[name="username"]').type(id).should('have.value', id);
  cy.get('input[name="password"]').type('password').should('have.value', 'password');
  cy.contains('Register').click();

  cy.url().should('include', '/channels/me');
});

Cypress.Commands.add('loginUser', (email) => {
  cy.visit('/');
  cy.contains('Valkyrie');
  cy.contains('Login').click();

  cy.url().should('include', '/login');
  cy.get('input[name="email"]').type(email).should('have.value', email);
  cy.get('input[name="password"]').type('password').should('have.value', 'password');
  cy.contains('Login').click();

  cy.url().should('include', '/channels/me');
});

Cypress.Commands.add('sendRequest', (id) => {
  cy.contains('Add Friend').click();
  cy.get('input[name=id]').type(id).should('have.value', id);
  cy.get('[type=submit]').click();
});

Cypress.Commands.add('createGuild', (id) => {
  cy.get('[id="add-guild-icon"]').click();

  cy.contains('Create My Own').click();
  cy.get('input').clear().type(`${id}'s server`);
  cy.get('[type=submit]').click();

  cy.contains(`${id}'s server`).should('exist');
  cy.contains(`Welcome to #general`).should('exist');
});

Cypress.Commands.add('clickOnFirstGuild', () => {
  cy.wait(1000);
  cy.get('a[href*="channels"]').not('a[href*="channels/me"]').first().click();
});

Cypress.Commands.add('openGuildMenu', () => {
  cy.get('[id^="menu-button-"]').click();
});

Cypress.Commands.add('joinGuild', (invite) => {
  cy.get('[id="add-guild-icon"]').click();
  cy.contains('Join a Server').click();

  cy.get('input').type(invite).should('have.value', invite);
  cy.get('[type=submit]').click();
  cy.wait(1000);
});

Cypress.Commands.add('getChat', () => {
  cy.get('.infinite-scroll-component');
});

Cypress.Commands.add('firstMessage', () => {
  cy.getChat().children().first();
});
