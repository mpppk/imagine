import actionCreatorFactory from 'typescript-fsa';
import {Asset, BoundingBoxRequest} from "../models/models";
import {WSPayload} from "./workspace";

const boundingBoxActionCreatorFactory = actionCreatorFactory('BOUNDING_BOX');

export interface BoundingBoxAssignRequestPayload extends WSPayload {
  asset: Asset
  box: BoundingBoxRequest
}

export const boundingBoxActionCreators = {
  assignRequest: boundingBoxActionCreatorFactory<BoundingBoxAssignRequestPayload>('ASSIGN/REQUEST'),
  unAssignRequest: boundingBoxActionCreatorFactory<BoundingBoxAssignRequestPayload>('UNASSIGN/REQUEST'),
};

