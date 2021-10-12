import { uuid } from '../support/utils';

describe('Friends related actions', () => {
  // Friend values
  const friendName = uuid();
  const friendEmail = `${friendName}@example.com`;
  let friendId = '';

  // AuthUser values
  const authName = uuid();
  const email = `${authName}@example.com`;

  it('registers the friend and gets the id', () => {
    cy.registerUser(friendEmail, friendName);

    // Get the friend's ID and store it
    cy.contains('Add Friend').click();
    cy.get('input')
      .first()
      .invoke('val')
      .then((text) => {
        friendId = text.toString();
      });
  });

  it('signs up the user and gets the id', () => {
    cy.registerUser(email, authName);
  });

  it('sends a friend request to the member', () => {
    cy.loginUser(email);
    cy.sendRequest(friendId);

    // Confirm the 'Pending' tab got an outgoing friend request
    cy.contains('Pending').click();
    cy.contains(friendName);
    cy.contains('Outgoing Friend Request');
  });

  it('cancels the outgoing request', () => {
    cy.loginUser(email);

    cy.contains('Pending').click();
    cy.contains(friendName);
    // Confirm the accept button does not exist
    cy.get('[aria-label="accept request"]').should('not.exist');
    cy.get('[aria-label="decline request"]').click();

    // Confirm the user got removed and is not a friend.
    cy.contains(friendName).should('not.exist');
    cy.get('button').contains('Friends').click();
    cy.contains(friendName).should('not.exist');
  });

  it('sends another friend request to the member to be checked', () => {
    cy.loginUser(email);
    cy.sendRequest(friendId);
  });

  it('checks the incoming friend request and declines it', () => {
    cy.loginUser(friendEmail);

    // Confirm the user got a request and has the accept button
    cy.contains('Pending').click();
    cy.contains(authName);
    cy.contains('Incoming Friend Request');
    cy.get('[aria-label="accept request"]').should('exist');

    // Decline request
    cy.get('[aria-label="decline request"]').click();
    cy.contains(authName).should('not.exist');

    // Confirm the user got removed and is not a friend.
    cy.get('button').contains('Friends').click();
    cy.wait(100);
    cy.contains(authName).should('not.exist');
  });

  it('sends another friend request to the member to be accepted', () => {
    cy.loginUser(email);
    cy.sendRequest(friendId);
  });

  it('accepts the friend request', () => {
    cy.loginUser(friendEmail);

    // Accept request
    cy.contains('Pending').click();
    cy.contains(authName);
    cy.get('[aria-label="accept request"]').click();
    cy.contains(authName).should('not.exist');

    // Confirm the user is in the 'Friends' tab
    cy.get('button').contains('Friends').click();
    cy.wait(100);
    cy.contains(authName).should('exist');
  });

  it('directs to the friends dms when clicked', () => {
    cy.loginUser(email);

    cy.intercept({
      method: 'POST',
      pathname: `/api/channels/${friendId}/dm`,
    }).as('getDM');

    // Open the DM
    cy.contains(friendName).should('exist').click();

    // Confirm it's the friend's DM
    cy.contains(friendName).should('exist');
    cy.contains(`This is the beginning of your direct message history with @${friendName}`).should('exist');
    cy.get('textarea[name="text"]').invoke('attr', 'placeholder').should('contain', `@${friendName}`);

    // Confirm the DM's url is the correct one
    cy.wait('@getDM').then((interception) => {
      const body = interception.response.body;
      const url = `/channels/me/${body.id}`;
      cy.url().should('include', url);
    });
  });

  it('should successfully go to the DM when clicked on the item', () => {
    cy.loginUser(email);
    cy.get('ul[id="dm-list"]').children().contains(friendName).click();

    // Confirm it's the friend's DM
    cy.contains(friendName).should('exist');
    cy.contains(`This is the beginning of your direct message history with @${friendName}`).should('exist');
    cy.get('textarea[name="text"]').invoke('attr', 'placeholder').should('contain', `@${friendName}`);

    // Check that messaging is possible and the message gets added to the chat
    cy.get('textarea[name="text"]').type('Hello World{enter}').should('have.value', '');
    cy.wait(50);
    cy.contains('Hello World');
  });

  it('closes the dm when the close button is pressed', () => {
    cy.loginUser(email);

    // Check that the DM exists
    cy.get('ul[id="dm-list"]').children().contains(friendName).should('exist');
    cy.get('ul[id="dm-list"] li:first').trigger('mouseover');

    // Close the DM and confirm it's gone
    cy.get('[aria-label="close dm"]').click();
    cy.get('ul[id="dm-list"]').children().contains(friendName).should('not.exist');
  });

  it('removes the friend', () => {
    cy.loginUser(email);

    // Check that the friend exists in the tab
    cy.get('ul[id="friend-list"]').children().contains(friendName).should('exist');

    // Confirm and remove friend
    cy.get('[aria-label="remove friend"]').click();
    cy.contains('Remove Friend').click();

    // Confirm the friend got removed
    cy.get('ul[id="friend-list"]').should('not.exist');
  });
});
