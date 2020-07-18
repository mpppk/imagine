import actionCreatorFactory from 'typescript-fsa';
import {WSPayload} from "./server";

const indexActionCreatorFactory = actionCreatorFactory('INDEX');

export const indexActionCreators = {
  clickAddDirectoryButton: indexActionCreatorFactory<WSPayload>(
    'CLICK_ADD_DIRECTORY_BUTTON'
  ),
};

