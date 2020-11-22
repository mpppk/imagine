import actionCreatorFactory from 'typescript-fsa';
import { Asset, BoundingBox, BoundingBoxRequest } from '../models/models';
import { WSPayload } from './workspace';
import { Pixel } from '../components/svg/svg';

const boundingBoxActionCreatorFactory = actionCreatorFactory('BOUNDING_BOX');

export interface BoundingBoxAssignRequestPayload extends WSPayload {
  asset: Asset;
  box: BoundingBoxRequest;
}

export interface BoundingBoxUnAssignRequestPayload extends WSPayload {
  asset: Asset;
  boxID: number;
}

export interface BoundingBoxAssignPayload extends WSPayload {
  asset: Asset;
  box: BoundingBoxRequest;
}

export interface BoundingBoxUnAssignPayload extends WSPayload {
  asset: Asset;
  boxID: number;
}

export interface BoundingBoxModifyPayload extends WSPayload {
  asset: Asset;
  box: BoundingBox;
}

export interface BoundingBoxMovePayload extends WSPayload {
  assetID: number;
  boxID: number;
  dx: Pixel;
  dy: Pixel;
}

export interface BoundingBoxScalePayload extends WSPayload {
  assetID: number;
  boxID: number;
  dx: Pixel;
  dy: Pixel;
}

export interface BoundingBoxDeletePayload extends WSPayload {
  assetID: number;
  boxID: number;
}

export const boundingBoxActionCreators = {
  assignRequest: boundingBoxActionCreatorFactory<BoundingBoxAssignRequestPayload>(
    'ASSIGN/REQUEST'
  ),
  assign: boundingBoxActionCreatorFactory<BoundingBoxAssignPayload>('ASSIGN'),
  unAssignRequest: boundingBoxActionCreatorFactory<BoundingBoxUnAssignRequestPayload>(
    'UNASSIGN/REQUEST'
  ),
  unAssign: boundingBoxActionCreatorFactory<BoundingBoxUnAssignPayload>(
    'UNASSIGN'
  ),
  modify: boundingBoxActionCreatorFactory<BoundingBoxModifyPayload>('MODIFY'),
  modifyRequest: boundingBoxActionCreatorFactory<BoundingBoxModifyPayload>(
    'MODIFY/REQUEST'
  ),
  move: boundingBoxActionCreatorFactory<BoundingBoxScalePayload>('MOVE'),
  scale: boundingBoxActionCreatorFactory<BoundingBoxScalePayload>('SCALE'),
  startScale: boundingBoxActionCreatorFactory<BoundingBoxScalePayload>(
    'SCALE/START'
  ),
  doneScale: boundingBoxActionCreatorFactory<BoundingBoxScalePayload>(
    'SCALE/DONE'
  ),
  deleteRequest: boundingBoxActionCreatorFactory<BoundingBoxDeletePayload>(
    'DELETE/REQUEST'
  ),
};
