import actionCreatorFactory from 'typescript-fsa';
import {WorkSpace} from "../models/models";

const workspaceActionCreatorFactory = actionCreatorFactory('WORKSPACE');

export interface WSPayload {
  workSpaceName: string
}

export const workspaceActionCreators = {
  requestWorkSpaces: workspaceActionCreatorFactory<void>('REQUEST_WORKSPACES'),
  selectNewWorkSpace: workspaceActionCreatorFactory<WorkSpace>('SELECT_NEW_WORKSPACE'),
  scanWorkSpaces: workspaceActionCreatorFactory<WorkSpace[]>('SCAN_WORKSPACES'),
};

