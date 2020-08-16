import actionCreatorFactory from 'typescript-fsa';
import {Tag} from "../models/models";
import {WSPayload} from "./workspace";

const indexActionCreatorFactory = actionCreatorFactory('INDEX');

export const indexActionCreators = {
  clickAddDirectoryButton: indexActionCreatorFactory<WSPayload>('ADD_DIRECTORY_BUTTON/CLICK'),
  clickAddTagButton: indexActionCreatorFactory<Tag>('ADD_TAG_BUTTON/CLICK'),
  clickEditTagButton: indexActionCreatorFactory<Tag>('EDIT_TAG_BUTTON/CLICK'),
};

