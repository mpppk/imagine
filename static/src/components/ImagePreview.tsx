import {Theme} from "@material-ui/core";
import {makeStyles} from "@material-ui/core/styles";
import React, {useMemo} from "react";
import {AssetWithIndex, BoundingBox} from "../models/models";
import {isDefaultBox} from "../util";
import {RectLayer} from "./svg/RectLayer";
import {Pixel} from "./svg/svg";

const useStyles = makeStyles((theme: Theme) => {
  return {
    rect: {
      fill: 'transparent',
      stroke: theme.palette.primary.light,
      strokeWidth: 10,
      cursor: 'move'
    }
  }
});

export type BoundingBoxModifyHandler = (box: BoundingBox) => void;

interface Props {
  src: string
  asset: AssetWithIndex
  onBoundingBoxModify: BoundingBoxModifyHandler
}

// const useHandlers = (props: Props) => {
//   return useMemo(() => {
//     return {}
//   }, [props]);
// }

const createRectProp = (onBoundingBoxModify: BoundingBoxModifyHandler) => (box: BoundingBox) => {
  const rectLayerProp = {
    onMove: (_dx: Pixel, _dy: Pixel) => {
      // onMove: (dx: Pixel, dy: Pixel) => {
      // FIXME
      onBoundingBoxModify({
        ...box,
        x: 0,
        y: 0,
        // x: box.x + dx,
        // y: box.y + dy,
      });
    },
    onScale: (width: Pixel, height: Pixel) => {
        onBoundingBoxModify({
          ...box,
          width, height,
        });
    },
    id: box.id,
    width: box.width as Pixel,
    height: box.height as Pixel,
    x: box.x as Pixel,
    y: box.y as Pixel,
    key: box.id
  };
  return isDefaultBox(box) ? {...rectLayerProp, width: '100%', height: '100%'} : rectLayerProp;
}

type Classes = ReturnType<typeof useStyles>;
const useViewState = (props: Props, classes: Classes) => {
  return useMemo(() => {
    const boxes = props.asset.boundingBoxes ?? [];
    const rectProps = boxes.map(createRectProp(props.onBoundingBoxModify));
    return {rectProps};
  }, [props, classes]);
}

// tslint:disable-next-line:variable-name
export const ImagePreview: React.FC<Props> = (props) => {
  const classes = useStyles();
  const viewState = useViewState(props, classes);
  // const handlers = useHandlers(props);

  return (
    <div>
      <svg id="canvas" viewBox="0 0 500 500" width="500" height="500">
        <image href={props.src} width={'100%'}/>
        {viewState.rectProps.map((rectProp) => {
          return <RectLayer
            src={rectProp}
            key={rectProp.key}
            className={classes.rect}
            onScale={rectProp.onScale}
          />
        })}
      </svg>
    </div>
  )
}
