import actionCreatorFactory from 'typescript-fsa';
import {Asset, BoundingBoxRequest} from "../models/models";
import {WSPayload} from "./workspace";

const boundingBoxActionCreatorFactory = actionCreatorFactory('BOUNDING_BOX');

export interface BoundingBoxAssignRequestPayload extends WSPayload {
  asset: Asset
  box: BoundingBoxRequest
}

export interface BoundingBoxUnAssignRequestPayload extends WSPayload {
  asset: Asset
  boxID: number
}

export interface BoundingBoxAssignPayload extends WSPayload {
  asset: Asset
  box: BoundingBoxRequest
}

export interface BoundingBoxUnAssignPayload extends WSPayload {
  asset: Asset
  boxID: number
}

export const boundingBoxActionCreators = {
  assignRequest: boundingBoxActionCreatorFactory<BoundingBoxAssignRequestPayload>('ASSIGN/REQUEST'),
  assign: boundingBoxActionCreatorFactory<BoundingBoxAssignPayload>('ASSIGN'),
  unAssignRequest: boundingBoxActionCreatorFactory<BoundingBoxUnAssignRequestPayload>('UNASSIGN/REQUEST'),
  unAssign: boundingBoxActionCreatorFactory<BoundingBoxUnAssignPayload>('UNASSIGN'),
};
