import actionCreatorFactory from 'typescript-fsa';
import {Tag} from "../models/models";
import {WSPayload} from "./workspace";

const tagActionCreatorFactory = actionCreatorFactory('TAG');

interface TagScanPayload extends WSPayload {
  tags: Tag[]
}

export const tagActionCreators = {
  tagScan: tagActionCreatorFactory<TagScanPayload>('SCAN')
};

