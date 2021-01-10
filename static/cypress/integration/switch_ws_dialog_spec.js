describe('switch workspace dialog', () => {
  it('can close by close button', () => {
    cy.visit('http://localhost:3000');

    cy.getBySel('open-switch-workspace-dialog-button').click();
    cy.getBySel('switch-workspace-dialog-title').should('exist');
    cy.getBySel('cancel-switch-workspace-button').click();
    cy.getBySel('switch-workspace-dialog-title').should('not.exist');
  });
});
