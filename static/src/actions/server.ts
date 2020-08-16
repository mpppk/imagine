import actionCreatorFactory from 'typescript-fsa';
import {Tag, WorkSpace} from "../models/models";

const serverActionCreatorFactory = actionCreatorFactory('SERVER');

export interface WSPayload {
  workSpaceName: string
}

interface TagScanPayload extends WSPayload {
  tags: Tag[]
}

export const serverActionCreators = {
  scanWorkSpaces: serverActionCreatorFactory<WorkSpace[]>(
    'SCAN_WORKSPACES'
  ),
  tagScan: serverActionCreatorFactory<TagScanPayload>('TAG/SCAN')
};
