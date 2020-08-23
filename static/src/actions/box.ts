import actionCreatorFactory from 'typescript-fsa';
import {Asset, BoundingBoxRequest} from "../models/models";

const boundingBoxActionCreatorFactory = actionCreatorFactory('BOUNDING_BOX');

export interface BoundingBoxAssignRequestPayload {
  asset: Asset
  box: BoundingBoxRequest
}

export const boundingBoxActionCreators = {
  assignRequest: boundingBoxActionCreatorFactory<BoundingBoxAssignRequestPayload>('ASSIGN/REQUEST'),
  unAssignRequest: boundingBoxActionCreatorFactory<BoundingBoxAssignRequestPayload>('UN_ASSIGN/REQUEST'),
};

