describe('filter assets dialog', () => {
  it('can close by cancel button', () => {
    cy.visit('http://localhost:3000');

    cy.getBySel('open-filter-assets-dialog-button').click();
    cy.getBySel('filter-assets-dialog-title').should('exist');
    cy.getBySel('cancel-filter-assets-dialog').click();
    cy.getBySel('filter-assets-dialog-title').should('not.exist');
  });

  it('can filter assets', () => {
    cy.visit('http://localhost:3000');

    cy.getBySel('open-filter-assets-dialog-button').click();
    cy.getBySel('filter-assets-toggle-button').click();
    cy.getBySel('add-new-assets-query-button').click();
    cy.getBySel('assets-query-tag-name-input').type('tag1');
    cy.getBySel('filter-assets-dialog-apply-button').click();
    cy.getBySel('filter-assets-dialog-title').should('not.exist');
  });
});
