import actionCreatorFactory from 'typescript-fsa';

const indexActionCreatorFactory = actionCreatorFactory('INDEX');

export const indexActionCreators = {
  clickAddDirectoryButton: indexActionCreatorFactory<void>(
    'CLICK_ADD_DIRECTORY_BUTTON'
  ),
};

