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

    // FIXME: check css
    cy.getBySel('tag-list-item').eq(2).click();

    cy.getBySel('tag-information-table').within(() => {
      cy.getBySel('tag-id').should('have.text', 3);
      cy.getBySel('tag-name').should('have.text', 'tag3');
    });

    // assign and unassign tag
    // FIXME: check css
    cy.getBySel('tag-list-item').eq(1).click();
    cy.getBySel('tag-list-item').eq(1).click();

    // assign and unassign tag by number key
    cy.getBySel('tag-list-drawer').type(1).type(1);

    // rename tag
    cy.getBySel('edit-tag-button').eq(1).click();
    cy.getBySel('tag-name-form').get('input').clear().type('tag2!');
    cy.getBySel('save-tag-name-button').click();
    cy.getBySel('tag-list-item').eq(1).contains('tag2!');

    // delete tag
    cy.getBySel('delete-tag-button').eq(1).click();
    cy.getBySel('tag-list-item').should('have.length', 2);
    cy.getBySel('tag-list-item').eq(1).should('not.have.text', 'tag2!');

    // add new tag
    cy.getBySel('add-new-tag-button').click();
    cy.getBySel('tag-name-form').get('input').type('tag4');
    cy.getBySel('save-tag-name-button').click();
    cy.getBySel('tag-list-item').first().should('have.text', '1: tag4');
  });
});
