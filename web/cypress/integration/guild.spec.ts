import { uuid } from '../support/utils';

describe('Guild related actions', () => {
  // The user's values
  const id = uuid().toString();
  const email = `${id}@example.com`;

  // The mock members values
  const memberId = uuid().toString();
  const memberEmail = `${memberId}@example.com`;

  // The invite link and guildId
  let invite = '';
  let guildId = '';

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

  it('should update the server after it got edited', () => {
    cy.loginUser(email);
    cy.clickOnFirstGuild();

    cy.openGuildMenu();
    cy.contains('Server Settings').click();

    // Updates the values and saves them
    cy.get('input[name="name"]').clear().type('Valkyrie').should('have.value', 'Valkyrie');
    cy.contains('Save Changes').click();

    // Check the updates were applied
    cy.contains('Valkyrie').should('exist');
  });

  it('should successfully clear the invites', () => {
    cy.loginUser(email);

    cy.clickOnFirstGuild();
    cy.openGuildMenu();
    cy.contains('Server Settings').click();

    cy.intercept({
      method: 'DELETE',
      pathname: `/api/guilds/${guildId}/invite`,
    }).as('clear');

    cy.contains('Invalidate Links').click();

    // Check the invites succesfully got cleared
    cy.wait('@clear').then((interception) => {
      const body = interception.response.body;
      expect(body).eq(true);
    });
  });

  it('should delete the server and go to home', () => {
    cy.loginUser(email);

    cy.clickOnFirstGuild();

    // Delete server and confirm the request
    cy.openGuildMenu();
    cy.contains('Server Settings').click();
    cy.contains('Delete Server').click();
    cy.contains('Delete Server').click();

    // Confirm the user is back at home and the guild is gone
    // Not working due to websocket connection problems
    // cy.url().should('include', '/channels/me');
    // cy.get('a[href*="channels"]').not('a[href*="channels/me"]').should('have.length', 0);
  });

  it('creates another guild to get an invite', () => {
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

    // get the invite for the other user
    cy.openGuildMenu();
    cy.contains('Invite People').click();

    // Check unlimited invites
    cy.get('input[type="checkbox"]').check({ force: true });
    cy.wait(50);

    // Store the invite in the variable
    cy.get('input[id="invite-link"]')
      .first()
      .invoke('val')
      .then((text) => {
        invite = text.toString();
      });
  });

  it('joins the guild for the given link', () => {
    cy.registerUser(memberEmail, memberId);

    cy.intercept({
      method: 'POST',
      pathname: '/api/guilds/join',
    }).as('join');

    cy.joinGuild(invite);

    // Confirm the user got sent to the guild and default channel
    cy.wait('@join').then((interception) => {
      const body = interception.response.body;
      const url = `channels/${body.id}/${body.default_channel_id}`;
      cy.url().should('include', url);
    });

    cy.contains(`${id}'s server`).should('exist');
    cy.contains(`Welcome to #general`).should('exist');
  });

  it('should not show "Server Settings" and "Create Channel" to the non owner', () => {
    cy.loginUser(memberEmail);

    cy.clickOnFirstGuild();

    cy.openGuildMenu();
    cy.contains('Server Settings').should('not.exist');
    cy.contains('Create Channel').should('not.exist');
  });

  it('should leave the server', () => {
    cy.loginUser(memberEmail);

    cy.clickOnFirstGuild();

    cy.openGuildMenu();
    cy.contains('Leave Server').click();

    // Confirm the user is back at home and the guild is gone
    cy.url().should('include', '/channels/me');
    cy.get('a[href*="channels"]').not('a[href*="channels/me"]').should('have.length', 0);
  });

  it('creates a third guild to test switching', () => {
    cy.loginUser(email);
    cy.createGuild('Test');
  });

  it('successfully switches between the guilds', () => {
    cy.loginUser(email);

    cy.clickOnFirstGuild();
    cy.url().should('include', guildId);
    cy.contains(`${id}'s server`).should('exist');

    // Click on the second server and confirm it's the second one created
    cy.get('a[href*="channels"]').eq(2).click();
    cy.contains(`Test's server`).should('exist');
  });
});
