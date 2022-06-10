import { uuid } from '../support/utils';

// Covers the account routes
describe('Account related pages', () => {
  const id = uuid();
  const email = `${id}@example.com`;

  it('redirects an unauthenticated user', () => {
    cy.visit('/channels/me');
    cy.url().should('include', '/login');
  });

  it('registers the user', () => {
    cy.registerUser(email, id);
  });

  it('signs in the user', () => {
    cy.loginUser(email);
  });

  it("checks the user's settings page", () => {
    cy.loginUser(email);
    cy.get('[aria-label=settings]').click();

    // Confirm account page got user's info
    cy.url().should('include', '/account');
    cy.contains('My Account'.toUpperCase());
    cy.get('input[name="email"]').should('have.value', email);
    cy.get('input[name="username"]').should('have.value', id);
  });

  it("updates the user's info", () => {
    cy.loginUser(email);
    cy.get('[aria-label=settings]').click();

    // Update the user's info and check for confirmation toast
    cy.url().should('include', '/account');
    cy.contains('My Account'.toUpperCase());
    cy.get('input[name="email"]').should('have.value', email);
    cy.get('input[name="username"]').clear().type('Test').should('have.value', 'Test');
    cy.get('[type=submit]').click();

    cy.wait(200);
    cy.contains('Account Updated.');
  });

  it("updates the user's password", () => {
    cy.loginUser(email);
    cy.get('[aria-label=settings]').click();

    cy.url().should('include', '/account');
    cy.contains('Change Password').click();

    // Change the user's password
    cy.contains('Change your password');
    cy.get('input[name="currentPassword"]').type('password').should('have.value', 'password');
    cy.get('input[name="newPassword"]').type('password').should('have.value', 'password');
    cy.get('input[name="confirmNewPassword"]').type('password').should('have.value', 'password');
    cy.contains('Done').click();

    cy.wait(200);
    cy.contains('Changed Password');
  });

  it('signs out the user', () => {
    cy.loginUser(email);
    cy.get('[aria-label=settings]').click();
    cy.url().should('include', '/account');

    cy.contains('Logout').click();

    cy.url().should('include', '');
  });
});
