import actionCreatorFactory from 'typescript-fsa';
import {Tag} from "../models/models";

const serverActionCreatorFactory = actionCreatorFactory('SERVER');

export interface WSPayload {
  workSpaceName: string
}

interface TagScanPayload extends WSPayload {
  tags: Tag[]
}

export const serverActionCreators = {
  tagScan: serverActionCreatorFactory<TagScanPayload>('TAG/SCAN')
};
