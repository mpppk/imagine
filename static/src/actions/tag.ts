import actionCreatorFactory from 'typescript-fsa';
import {Tag} from "../models/models";
import {WSPayload} from "./workspace";

const tagActionCreatorFactory = actionCreatorFactory('TAG');

interface TagsPayload extends WSPayload {
  tags: Tag[]
}

interface TagPayload extends WSPayload {
  tag: Tag
}

export const tagActionCreators = {
  scanResult: tagActionCreatorFactory<TagsPayload>('SCAN/RESULT'),
  update: tagActionCreatorFactory<TagsPayload>('UPDATE'),
  rename: tagActionCreatorFactory<TagPayload>('RENAME'),
  save: tagActionCreatorFactory<TagsPayload>('SAVE'),
};

