import actionCreatorFactory from 'typescript-fsa';
import {WSPayload} from "./server";
import {Tag} from "../models/models";

const indexActionCreatorFactory = actionCreatorFactory('INDEX');

interface UpdateTagsPayload extends WSPayload {
  tags: Tag[]
}

interface TagPayload extends WSPayload {
  tag: Tag
}

export const indexActionCreators = {
  clickAddDirectoryButton: indexActionCreatorFactory<WSPayload>('CLICK_ADD_DIRECTORY_BUTTON'),
  clickAddTagButton: indexActionCreatorFactory<Tag>('CLICK_ADD_TAG_BUTTON'),
  clickEditTagButton: indexActionCreatorFactory<Tag>('CLICK_EDIT_TAG_BUTTON'),
  renameTag: indexActionCreatorFactory<TagPayload>('RENAME_TAG'),
  updateTags: indexActionCreatorFactory<UpdateTagsPayload>('UPDATE_TAGS'),
};

