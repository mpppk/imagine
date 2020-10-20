import actionCreatorFactory from 'typescript-fsa';
import {Asset, Direction, Tag} from "../models/models";
import {WSPayload} from "./workspace";
import {BoundingBoxModifyPayload} from "./box";

const indexActionCreatorFactory = actionCreatorFactory('INDEX');

export const indexActionCreators = {
  selectTag: indexActionCreatorFactory<Tag>('TAG/SELECT'),
  clickAddDirectoryButton: indexActionCreatorFactory<WSPayload>('ADD_DIRECTORY_BUTTON/CLICK'),
  clickAddTagButton: indexActionCreatorFactory<Tag>('ADD_TAG_BUTTON/CLICK'),
  clickEditTagButton: indexActionCreatorFactory<Tag>('EDIT_TAG_BUTTON/CLICK'),
  downNumberKey: indexActionCreatorFactory<number>('NUMBER_KEY/DOWN'),
  downArrowKey: indexActionCreatorFactory<Direction>('ARROW_KEY/DOWN'),
  assetSelect: indexActionCreatorFactory<Asset>('ASSET/SELECT'),
  dragResizeHandler: indexActionCreatorFactory<BoundingBoxModifyPayload>('RESIZE_HANDLER/DRAG'),
};

