import actionCreatorFactory from 'typescript-fsa';
import {Tag} from "../models/models";
import {WSPayload} from "./workspace";

const indexActionCreatorFactory = actionCreatorFactory('INDEX');

export const indexActionCreators = {
  clickAddDirectoryButton: indexActionCreatorFactory<WSPayload>('CLICK_ADD_DIRECTORY_BUTTON'),
  clickAddTagButton: indexActionCreatorFactory<Tag>('CLICK_ADD_TAG_BUTTON'),
  clickEditTagButton: indexActionCreatorFactory<Tag>('CLICK_EDIT_TAG_BUTTON'),
};

