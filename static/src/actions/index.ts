import actionCreatorFactory from 'typescript-fsa';
import {Asset, Tag} from "../models/models";
import {WSPayload} from "./workspace";

const indexActionCreatorFactory = actionCreatorFactory('INDEX');

export const indexActionCreators = {
  clickAddDirectoryButton: indexActionCreatorFactory<WSPayload>('ADD_DIRECTORY_BUTTON/CLICK'),
  clickAddTagButton: indexActionCreatorFactory<Tag>('ADD_TAG_BUTTON/CLICK'),
  clickEditTagButton: indexActionCreatorFactory<Tag>('EDIT_TAG_BUTTON/CLICK'),
  downNumberKey: indexActionCreatorFactory<number>('NUMBER_KEY/DOWN'),
  assetSelect: indexActionCreatorFactory<Asset>('ASSET/SELECT'),
};

