const baseUrl = 'http://localhost:3000/';
describe('hamburger menu', () => {
  it('has home button', () => {
    cy.visit(baseUrl);

    cy.getBySel('hamburger-menu-button').click();
    cy.getBySel('sidelist-home').click();

    cy.url().should('eq', baseUrl);
  });
});
