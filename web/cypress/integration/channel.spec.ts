import { uuid } from '../support/utils';

describe('Channels related actions', () => {
  // The user's values
  const id = uuid().toString();
  const email = `${id}@example.com`;
  let guildId = '';
  let channelName = '';

  it('registers the user', () => {
    cy.registerUser(email, id);
  });

  it('creates a guild', () => {
    cy.loginUser(email);

    cy.intercept({
      method: 'POST',
      pathname: '/api/guilds/create',
    }).as('create');

    cy.createGuild(id);

    // Confirm the user got sent to the guild and default channel
    cy.wait('@create').then((interception) => {
      const body = interception.response.body;
      const url = `channels/${body.id}/${body.default_channel_id}`;
      cy.url().should('include', url);
      guildId = body.id;
    });
  });

  it('creates a channel for the guild', () => {
    cy.loginUser(email);

    cy.clickOnFirstGuild();
    cy.openGuildMenu();

    cy.intercept({
      method: 'POST',
      pathname: `/api/channels/${guildId}`,
    }).as('create');

    cy.contains('Create Channel').click();
    cy.get('input[name="name"]').type('random').should('have.value', 'random');
    cy.get('[type=submit]').click();

    // Confirm the user got sent to the newly created channel
    cy.wait('@create').then((interception) => {
      const body = interception.response.body;
      const url = `channels/${guildId}/${body.id}`;
      cy.url().should('include', url);
    });

    cy.contains('random').should('exist');
  });

  it('should successfully switch between channels', () => {
    cy.loginUser(email);
    cy.clickOnFirstGuild();
    cy.contains('Welcome to #general').should('exist');

    // Check the other channel is the correct one
    cy.get(`a[href^="/channels/${guildId}/"]`).last().click();
    cy.contains('Welcome to #random').should('exist');
  });

  it('should successfully edit the channel', () => {
    cy.loginUser(email);
    cy.clickOnFirstGuild();

    cy.contains('random').trigger('mouseover');
    cy.get('[aria-label="edit channel"]').click();

    // Edit the values
    cy.get('input[name="name"]').clear().type('secret').should('have.value', 'secret');
    cy.get('input[type="checkbox"]').check({ force: true });
    cy.get('[type=submit]').click();

    // Check that the edited channel exists
    cy.contains('secret').should('exist');
  });

  it('should successfully delete the channel', () => {
    cy.loginUser(email);
    cy.clickOnFirstGuild();

    cy.get(`a[href^="/channels/${guildId}/"]`).last().click();
    cy.contains('secret').last().trigger('mouseover');
    cy.get('[aria-label="edit channel"]').last().click();

    cy.contains('Delete Channel').click();
    cy.get('button').contains('Delete Channel').click();

    // Check that the channel is gone and that the user got moved
    cy.contains('secret').should('have.length.lte', 1);
    cy.contains('Welcome to #general').should('exist');
  });
});
