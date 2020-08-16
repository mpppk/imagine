import actionCreatorFactory from 'typescript-fsa';
import {Tag} from "../models/models";
import {WSPayload} from "./workspace";

const tagActionCreatorFactory = actionCreatorFactory('TAG');

interface TagScanPayload extends WSPayload {
  tags: Tag[]
}

interface UpdateTagsPayload extends WSPayload {
  tags: Tag[]
}

interface TagPayload extends WSPayload {
  tag: Tag
}

export const tagActionCreators = {
  scan: tagActionCreatorFactory<TagScanPayload>('SCAN'),
  update: tagActionCreatorFactory<UpdateTagsPayload>('UPDATE'),
  rename: tagActionCreatorFactory<TagPayload>('RENAME'),
};

