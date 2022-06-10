import { uuid } from '../support/utils';

describe('Members related actions', () => {
  // Member values
  const memberName = uuid();
  const memberEmail = `${memberName}@example.com`;
  let invite = '';

  // AuthUser values
  const authName = uuid();
  const email = `${authName}@example.com`;

  it('registers the user', () => {
    cy.registerUser(email, authName);
  });

  it('creates a guild and get the invite', () => {
    cy.loginUser(email);
    cy.createGuild(authName);

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

  it('successfully changes the members settings', () => {
    cy.loginUser(email);
    cy.clickOnFirstGuild();

    cy.openGuildMenu();
    cy.contains('Change Appearance').click();

    // Sets the member values
    cy.get('input[name="nickname"]').clear().type('Tester').should('have.value', 'Tester');
    cy.get('div[title="#0693E3"]').click();
    cy.wait(100);

    cy.contains('Save').click();
  });

  it('joins the guild for the given link', () => {
    cy.registerUser(memberEmail, memberName);
    cy.joinGuild(invite);

    // Confirm the above changes
    cy.wait(100);
    cy.contains('Tester').should('exist');

    // Check that the two members + two labels are there
    cy.get('ul[id="member-list"]').children().should('have.length', 4);
  });

  it('should not show "Kick" and "Ban" options to the non owner', () => {
    cy.loginUser(memberEmail);
    cy.clickOnFirstGuild();

    cy.get('ul[id="member-list"]').contains('Tester').rightclick();
    cy.contains('Ban').should('not.exist');
    cy.contains('Kick').should('not.exist');
  });

  it('successfully resets the members settings', () => {
    cy.loginUser(email);
    cy.clickOnFirstGuild();

    cy.openGuildMenu();
    cy.contains('Change Appearance').click();

    // Reset the values
    cy.contains('Reset Nickname').click();
    cy.contains('Reset Color').click();
    cy.wait(100);

    cy.contains('Save').click();
  });

  it('should go to the members DMs on click', () => {
    cy.loginUser(email);
    cy.clickOnFirstGuild();

    cy.contains(memberName).rightclick();
    cy.contains('Message').click();

    // Confirm the user got moved to the correct DM
    cy.contains(memberName).should('exist');
    cy.contains(`This is the beginning of your direct message history with @${memberName}`).should('exist');
    cy.get('textarea[name="text"]').invoke('attr', 'placeholder').should('contain', `@${memberName}`);
  });

  it('should successfully sent a friends request', () => {
    cy.loginUser(email);
    cy.clickOnFirstGuild();

    cy.contains(memberName).rightclick();
    cy.contains('Add Friend').click();

    // Go to the pending tab to confirm the request
    cy.get('a[href="/channels/me"]').click();
    cy.contains('Pending').click();
    cy.contains(memberName).should('exist');
    cy.contains('Outgoing Friend Request').should('exist');
  });

  it('should kick and remove the member', () => {
    cy.loginUser(email);
    cy.clickOnFirstGuild();

    // Kick the user
    cy.contains(memberName).rightclick();
    cy.contains(`Kick ${memberName}`).click();
    cy.get('button').contains('Kick').click();

    // Confirm the member is gone
    cy.get('ul[id="member-list"]').children().should('have.length', 3);
    cy.contains(memberName).should('not.exist');
  });

  it('should successfully rejoin the server', () => {
    cy.loginUser(memberEmail);
    cy.joinGuild(invite);
  });

  it('should ban and remove the member', () => {
    cy.loginUser(email);
    cy.clickOnFirstGuild();

    // Ban the user
    cy.contains(memberName).rightclick();
    cy.contains(`Ban ${memberName}`).click();
    cy.get('button').contains('Ban').click();

    // Confirm the member is gone
    cy.get('ul[id="member-list"]').children().should('have.length', 3);
    cy.contains(memberName).should('not.exist');
  });

  it('should not be able to rejoin the server', () => {
    cy.loginUser(memberEmail);
    cy.joinGuild(invite);

    // Confirm the user did not join the guild
    cy.contains('You are banned from this server').should('exist');
    cy.url().should('include', '/channels/me');
  });

  it('should contain the member in the ban list', () => {
    cy.loginUser(email);
    cy.clickOnFirstGuild();

    // Go to the ban list modal
    cy.openGuildMenu();
    cy.contains('Server Settings').click();
    cy.contains('Bans').click();

    // Check the member exists there
    cy.wait(100);
    cy.contains(memberName).should('exist');
    cy.get('button[aria-label="unban user"]').click();

    // Confirm the member got unbanned
    cy.contains(memberName).should('not.exist');
  });

  it('should not display a context menu when clicking on oneself', () => {
    cy.loginUser(email);
    cy.clickOnFirstGuild();

    cy.contains(authName).rightclick();
    cy.contains('Add Friend').should('not.exist');
  });

  it('should toggle the member list on click', () => {
    cy.loginUser(email);
    cy.clickOnFirstGuild();

    // Toggle the member list
    cy.get('ul[id="member-list"]').should('be.visible');
    cy.get('[aria-label="toggle member list"]').click();
    cy.get('ul[id="member-list"]').should('not.exist');
    cy.get('[aria-label="toggle member list"]').click();
    cy.get('ul[id="member-list"]').should('be.visible');
  });
});
