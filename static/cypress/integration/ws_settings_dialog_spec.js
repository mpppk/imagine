describe('workspace settings dialog', () => {
  it('can close by close button', () => {
    cy.visit('http://localhost:3000');

    cy.getBySel('workspace-settings-button').click();
    cy.getBySel('workspace-settings-dialog-title').should('exist');
    cy.getBySel('workspace-settings-dialog-close-button').click();
    cy.getBySel('workspace-settings-dialog-title').should('not.exist');
  });

  it('can change base path', () => {
    cy.visit('http://localhost:3000');

    cy.getBySel('workspace-settings-button').click();
    cy.getBySel('workspace-settings-dialog-title').should('exist');
    cy.getBySel('workspace-settings-dialog-change-base-path-button').click();
    cy.getBySel('workspace-settings-dialog-base-path').should(
      'have.text',
      'new-base-path'
    );
  });
});
