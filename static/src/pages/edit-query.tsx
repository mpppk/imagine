import {NextPage} from 'next';
import React, {useState} from 'react';
import {DragDropContext, Draggable, Droppable, resetServerContext} from "react-beautiful-dnd";
import {ListItem} from "@material-ui/core";

interface Item {
  id: string
  content: string
}

// a little function to help us with reordering the result
const reorder = (list: any, startIndex: number, endIndex: number) => {
  const result = Array.from(list);
  const [removed] = result.splice(startIndex, 1);
  result.splice(endIndex, 0, removed);

  return result;
};

// fake data generator
const generateItems = (count: number) =>
  Array.from({ length: count }, (_, k) => k).map(k => ({
    id: `item-${k}`,
    content: `item ${k}`
  } as Item));

const grid = 8;

const getItemStyle = (isDragging: boolean, draggableStyle: any) => ({
  // some basic styles to make the items look a bit nicer
  userSelect: "none",
  padding: grid * 2,
  margin: `0 0 ${grid}px 0`,

  // change background colour if dragging
  background: isDragging ? "lightgreen" : "grey",

  // styles we need to apply on draggables
  ...draggableStyle
});

const getListStyle = (isDraggingOver: boolean) => ({
  background: isDraggingOver ? "lightblue" : "lightgrey",
  padding: grid,
  width: 250
});

// tslint:disable-next-line variable-name
export const EditQuery: NextPage = () => {
  const [items, setItems] = useState(generateItems(10));
  const onDragEnd = (result: any) => {
    // dropped outside the list
    if (!result.destination) {
      return;
    }

    const newItems = reorder(
      items,
      result.source.index,
      result.destination.index
    ) as Item[];

    setItems(newItems)
  }
  return (
    <div>
      <DragDropContext onDragEnd={onDragEnd}>
        <Droppable droppableId="droppable">
          {(provided, snapshot) => (
            <div
              {...provided.droppableProps}
              ref={provided.innerRef}
              style={getListStyle(snapshot.isDraggingOver)}
            >
              {items.map((item, index) => (
                <Draggable key={item.id} draggableId={item.id} index={index}>
                  {(provided, snapshot) => (
                    <ListItem
                      ref={provided.innerRef}
                      {...provided.draggableProps}
                      {...provided.dragHandleProps}
                      style={getItemStyle(
                        snapshot.isDragging,
                        provided.draggableProps.style
                      )}
                    >
                      {item.content}
                    </ListItem>
                  )}
                </Draggable>
              ))}
              {provided.placeholder}
            </div>
          )}
        </Droppable>
      </DragDropContext></div>
  );
};

export async function getServerSideProps() {
  resetServerContext()
  return {props: {}}
}

export default EditQuery;
