import actionCreatorFactory from 'typescript-fsa';
import { WorkSpace } from '../models/models';

const workspaceActionCreatorFactory = actionCreatorFactory('WORKSPACE');

export interface WSPayload {
  workSpaceName: string;
}

export interface WorkSpaceScanResultPayload {
  basePath: string;
  workspaces: WorkSpace[];
}

export const workspaceActionCreators = {
  scanRequest: workspaceActionCreatorFactory<void>('SCAN/REQUEST'),
  select: workspaceActionCreatorFactory<WorkSpace>('SELECT'),
  scanResult: workspaceActionCreatorFactory<WorkSpaceScanResultPayload>(
    'SCAN/RESULT'
  ),
  updateRequest: workspaceActionCreatorFactory<WorkSpace>('UPDATE/REQUEST'),
};
