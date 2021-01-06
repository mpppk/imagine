const c = (l) => `[data-cy=${l}]`;

describe('asset and tag list', () => {
  it('show selected asset', () => {
    cy.visit('http://localhost:3000');

    cy.getBySel('image-grid-list').get('li:first-child').click();

    cy.getBySel('asset-information-table').within(() => {
      cy.getBySel('asset-id').should('have.text', 1);
      cy.getBySel('asset-path').should('have.text', 'path1');
      cy.getBySel('asset-tags').should('have.text', 'tag1');
    });

    cy.getBySel('tag-list-item').eq(0).should('have.text', '1: tag1');
    cy.getBySel('tag-list-item').eq(1).should('have.text', '2: tag2');
    cy.getBySel('tag-list-item').eq(2).should('have.text', '3: tag3');

    cy.getBySel('tag-list-item').first().click();

    cy.getBySel('tag-information-table').within(() => {
      cy.getBySel('tag-id').should('have.text', 1);
      cy.getBySel('tag-name').should('have.text', 'tag1');
    });

    cy.getBySel('image-grid-drawer').type('{downarrow}');

    // FIXME: move 2 grid
    cy.getBySel('asset-information-table').within(() => {
      cy.getBySel('asset-id').should('have.text', 3);
      cy.getBySel('asset-path').should('have.text', 'path3');
      cy.getBySel('asset-tags').should('have.text', 'tag3');
    });

    cy.getBySel('tag-list-item').eq(2).click();

    cy.getBySel('tag-information-table').within(() => {
      cy.getBySel('tag-id').should('have.text', 3);
      cy.getBySel('tag-name').should('have.text', 'tag3');
    });
  });
});
