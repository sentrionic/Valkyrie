import { uuid } from '../support/utils';

const waitTime = 1000;

describe('Message related actions', () => {
  // The user's values
  const id = uuid().toString();
  const email = `${id}@example.com`;

  // The mock members values
  const memberName = uuid().toString();
  const memberEmail = `${memberName}@example.com`;

  let userId = '';

  // The invite link and guildId
  let invite = '';

  it('should register the user and create a guild', () => {
    cy.intercept({
      method: 'POST',
      pathname: '/api/account/register',
    }).as('register');

    cy.registerUser(email, id);

    cy.wait('@register').then((interception) => {
      const body = interception.response.body;
      userId = body.id;
    });

    cy.createGuild(id);
  });

  it('should successfully post a message', () => {
    cy.loginUser(email);
    cy.clickOnFirstGuild();

    // Send message
    cy.get('textarea[name="text"]').type('Hello, World{enter}').should('have.value', '');

    // Confirm it got added to the chat
    // cy.wait(100);
    // cy.getChat().should('exist');
    // cy.getChat().children().should('have.length', 1);
    // cy.firstMessage().contains('Hello, World');
  });

  it('should confirm that the message got sent', () => {
    cy.loginUser(email);
    cy.clickOnFirstGuild();

    cy.wait(waitTime);
    cy.getChat().should('exist');
    cy.getChat().children().should('have.length', 1);
    cy.firstMessage().contains('Hello, World');
  });

  it('should post another message', () => {
    cy.loginUser(email);
    cy.clickOnFirstGuild();

    cy.get('textarea[name="text"]').type('Hello, Server{enter}').should('have.value', '');
  });

  it('should display newer messages under older', () => {
    cy.loginUser(email);
    cy.clickOnFirstGuild();

    // Confirm it got added under the first one
    cy.wait(waitTime);
    cy.getChat().children().should('have.length', 2);
    cy.firstMessage().contains('Hello, Server');
    cy.getChat().children().last().contains('Hello, World');
  });

  it('should successfully delete the message', () => {
    cy.loginUser(email);
    cy.clickOnFirstGuild();

    cy.firstMessage().rightclick();
    cy.contains('Delete Message').click();
    cy.get('button').contains('Delete').click();
  });

  // Manually check the message is gone because of websocket problems
  it('should confirm that the message got deleted', () => {
    cy.loginUser(email);
    cy.clickOnFirstGuild();

    cy.wait(waitTime);
    cy.getChat().children().should('have.length', 1);
    cy.getChat().children().contains('Hello, Server').should('not.exist');
  });

  it('should successfully edit the message', () => {
    cy.loginUser(email);
    cy.clickOnFirstGuild();

    cy.firstMessage().rightclick();
    cy.contains('Edit Message').click();
    cy.wait(waitTime);

    cy.get('input[id="editMessage"]').clear().type('Hello, Update');
    cy.get('button').contains('Save').click();

    // // Confirm it got edited and got the edit span
    // cy.firstMessage().contains('Hello, Update');
    // cy.firstMessage().contains('(edited)');
  });

  it('should confirm the message got edited', () => {
    cy.loginUser(email);
    cy.clickOnFirstGuild();

    // Confirm it got edited and got the edit span
    cy.firstMessage().contains('Hello, Update');
    cy.firstMessage().contains('(edited)');
  });

  it('should not display "Add Friend" or "Message" for own avatar', () => {
    cy.loginUser(email);
    cy.clickOnFirstGuild();

    cy.wait(waitTime);
    cy.firstMessage().get('img').eq(1).rightclick();
    cy.contains('Add Friend').should('not.exist');
    cy.contains('Message').should('not.exist');
  });

  it('should get an invite for the other member', () => {
    cy.loginUser(email);
    cy.clickOnFirstGuild();

    // get the invite for the other user
    cy.openGuildMenu();
    cy.contains('Invite People').click();

    // Check unlimited invites
    cy.get('input[type="checkbox"]').check({ force: true });
    cy.wait(waitTime);

    // Store the invite in the variable
    cy.get('input[id="invite-link"]')
      .first()
      .invoke('val')
      .then((text) => {
        invite = text.toString();
      });
  });

  it('should register the other member and join the guild', () => {
    cy.registerUser(memberEmail, memberName);
    cy.joinGuild(invite);
  });

  it('should not display message options for the member context menu when clicking on the owners message', () => {
    cy.loginUser(memberEmail);
    cy.clickOnFirstGuild();

    cy.getChat().children().last().rightclick();
    cy.contains('Edit Message').should('not.exist');
    cy.contains('Delete Message').should('not.exist');
  });

  it('should successfully add the user', () => {
    cy.loginUser(memberEmail);
    cy.clickOnFirstGuild();

    cy.intercept({
      method: 'POST',
      pathname: `/api/account/${userId}/friend`,
    }).as('addFriend');

    cy.wait(waitTime);
    cy.firstMessage().get('img').eq(1).rightclick();
    cy.contains('Add Friend').click();

    cy.wait('@addFriend').then((_) => {
      // Go to the pending tab to confirm the request
      cy.get('a[href="/channels/me"]').click();
      cy.contains('Pending').click();
      cy.contains(memberName).should('exist');
      cy.contains('Outgoing Friend Request').should('exist');
    });
  });

  it('should successfully go to the DMs with the user', () => {
    cy.loginUser(memberEmail);
    cy.clickOnFirstGuild();

    cy.intercept({
      method: 'POST',
      pathname: `/api/channels/${userId}/dm`,
    }).as('create');

    cy.wait(waitTime);
    cy.firstMessage().get('img').eq(1).rightclick();
    cy.contains('Message').click();

    cy.wait('@create').then(() => {
      // Confirm the user got moved to the correct DM
      cy.contains(id).should('exist');
      cy.contains(`This is the beginning of your direct message history with @${id}`).should('exist');
      cy.get('textarea[name="text"]').invoke('attr', 'placeholder').should('contain', `@${id}`);
    });
  });

  it('should successfully post a message', () => {
    cy.loginUser(memberEmail);
    cy.clickOnFirstGuild();

    cy.get('textarea[name="text"]').type('Hello, Owner{enter}').should('have.value', '');
  });

  it("should be able to delete the member's message as the owner", () => {
    cy.loginUser(email);
    cy.clickOnFirstGuild();

    // Owner can delete the other person's message, but cannot edit it
    cy.firstMessage().rightclick();
    cy.contains('Edit Message').should('not.exist');

    cy.contains('Delete Message').click();
    cy.get('button').contains('Delete').click();
  });

  it('should confirm that the message got deleted', () => {
    cy.loginUser(email);
    cy.clickOnFirstGuild();

    // Confirm the message got deleted
    cy.wait(100);
    cy.getChat().children().should('have.length', 1);
    cy.getChat().children().contains('Hello, Owner').should('not.exist');
  });
});
