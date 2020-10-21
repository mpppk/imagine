import {Theme} from "@material-ui/core";
import {makeStyles} from "@material-ui/core/styles";
import React, {useRef} from "react";
import {AssetWithIndex, BoundingBox} from "../models/models";
import {RectLayer} from "./svg/RectLayer";
import {Pixel} from "./svg/svg";

const useStyles = makeStyles((theme: Theme) => {
  return {
    rect: {
      fill: 'transparent',
      stroke: theme.palette.primary.light,
      strokeWidth: 4,
      cursor: 'move'
    }
  }
});

// export type BoundingBoxModifyHandler = (box: BoundingBox) => void;

interface Props {
  src: string
  asset: AssetWithIndex
  onMoveBoundingBox: (box: BoundingBox, dx: Pixel, dy: Pixel) => void
  onScaleBoundingBox: (box: BoundingBox, dx: Pixel, dy: Pixel) => void
}

// const useHandlers = (props: Props) => {
//   return useMemo(() => {
//     return {}
//   }, [props]);
// }

// const createRectProp = (onBoundingBoxModify: BoundingBoxModifyHandler, imageRef: RefObject<SVGImageElement>) => (box: BoundingBox) => {
//   console.log('update rect prop', box);
//   const rectLayerProp = {
//     onMove: (dx: Pixel, dy: Pixel) => {
//       console.log('move', box);
//       const imageWidth = imageRef.current?.width.baseVal.value ?? 0;
//       const imageHeight = imageRef.current?.height.baseVal.value ?? 0;
//       onBoundingBoxModify({
//         ...box,
//         x: Math.max(Math.min(box.x + dx, imageWidth - box.width), 0),
//         y: Math.max(Math.min(box.y + dy, imageHeight - box.height), 0),
//       });
//     },
//     onScale: (width: Pixel, height: Pixel) => {
//       console.log('scale', box);
//       onBoundingBoxModify({
//         ...box,
//         width, height,
//       });
//     },
//     id: box.id,
//     width: box.width as Pixel,
//     height: box.height as Pixel,
//     x: box.x as Pixel,
//     y: box.y as Pixel,
//     key: box.id
//   };
//
//   const w = imageRef.current?.width.baseVal.value ?? 0;
//   const h = imageRef.current?.height.baseVal.value ?? 0;
//   return isDefaultBox(box) ? {...rectLayerProp, width: w, height: h} : rectLayerProp;
// }

// const useViewState = (props: Props, imageRef: RefObject<SVGImageElement>) => {
//   const boxes = props.asset.boundingBoxes ?? [];
//   // const rectProps = boxes.map(createRectProp(props.onBoundingBoxModify, imageRef));
//   const rectProps = boxes.map((box) => {
//     const rectLayerProp = {
//       onMove: (dx: Pixel, dy: Pixel) => {
//         const imageWidth = imageRef.current?.width.baseVal.value ?? 0;
//         const imageHeight = imageRef.current?.height.baseVal.value ?? 0;
//         props.onBoundingBoxModify({
//           ...box,
//           x: Math.max(Math.min(box.x + dx, imageWidth - box.width), 0),
//           y: Math.max(Math.min(box.y + dy, imageHeight - box.height), 0),
//         });
//       },
//       onScale: (width: Pixel, height: Pixel) => {
//         props.onBoundingBoxModify({
//           ...box,
//           width, height,
//         });
//       },
//       id: box.id,
//       width: box.width as Pixel,
//       height: box.height as Pixel,
//       x: box.x as Pixel,
//       y: box.y as Pixel,
//       key: box.id
//     };
//
//     const w = imageRef.current?.width.baseVal.value ?? 0;
//     const h = imageRef.current?.height.baseVal.value ?? 0;
//     return isDefaultBox(box) ? {...rectLayerProp, width: w, height: h} : rectLayerProp;
//   });
//   return {rectProps};
// }

// tslint:disable-next-line:variable-name
export const ImagePreview: React.FC<Props> = (props) => {
  const imageRef = useRef<SVGImageElement>(null);

  const classes = useStyles();
  // const viewState = useViewState(props, imageRef);
  // const handlers = useHandlers(props);

  // const boxes = props.asset.boundingBoxes ?? [];
  // const rectProps = boxes.map(createRectProp(props.onBoundingBoxModify, imageRef));
  // const rectProps = boxes.map((box) => {
  //   const rectLayerProp = {
  //     onMove: (dx: Pixel, dy: Pixel) => {
  //       const imageWidth = imageRef.current?.width.baseVal.value ?? 0;
  //       const imageHeight = imageRef.current?.height.baseVal.value ?? 0;
  //       props.onBoundingBoxModify({
  //         ...box,
  //         x: Math.max(Math.min(box.x + dx, imageWidth - box.width), 0),
  //         y: Math.max(Math.min(box.y + dy, imageHeight - box.height), 0),
  //       });
  //     },
  //     onScale: (width: Pixel, height: Pixel) => {
  //       props.onBoundingBoxModify({
  //         ...box,
  //         width, height,
  //       });
  //     },
  //     id: box.id,
  //     width: box.width as Pixel,
  //     height: box.height as Pixel,
  //     x: box.x as Pixel,
  //     y: box.y as Pixel,
  //     key: box.id
  //   };
  //
  //   const w = imageRef.current?.width.baseVal.value ?? 0;
  //   const h = imageRef.current?.height.baseVal.value ?? 0;
  //   return isDefaultBox(box) ? {...rectLayerProp, width: w, height: h} : rectLayerProp;
  // });

  const boxes = props.asset.boundingBoxes ?? [];

  // const genMoveHandler = (box: BoundingBox) => (dx: Pixel, dy: Pixel) => {
  //   props.onBoundingBoxModify({
  //     ...box,
  //     x: dx, y: dy,
  //   })
  // };

  return (
    <div>
      <svg id="canvas" viewBox="0 0 500 500" width="500" height="500">
        <image href={props.src} width={'100%'} height={'100%'} ref={imageRef}/>
        {boxes.map((box) => {
          const handleScale = (dx: Pixel, dy: Pixel) => {
            props.onScaleBoundingBox(box, dx, dy)
          }
          const handleMove = (dx: Pixel, dy: Pixel) => {
            props.onMoveBoundingBox(box, dx, dy)
          }
          return <RectLayer
            className={classes.rect}
            key={box.id}
            onScale={handleScale}
            onMove={handleMove}
            height={box.height}
            id={box.id}
            width={box.width}
            x={box.x}
            y={box.y}
          />
        })}
        {/*{viewState.rectProps.map((rectProp) => {*/}
        {/*  return <RectLayer*/}
        {/*    src={rectProp}*/}
        {/*    key={rectProp.key}*/}
        {/*    className={classes.rect}*/}
        {/*    onScale={rectProp.onScale}*/}
        {/*  />*/}
        {/*})}*/}
      </svg>
    </div>
  )
}
