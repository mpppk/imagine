import {Theme} from "@material-ui/core";
import Card from "@material-ui/core/Card";
import CardContent from "@material-ui/core/CardContent";
import List from "@material-ui/core/List";
import {makeStyles} from "@material-ui/core/styles";
import Typography from "@material-ui/core/Typography";
import {NextPage} from 'next';
import React, {useState} from 'react';
import {DragDropContext, Draggable, Droppable, resetServerContext} from "react-beautiful-dnd";

const useStyles = makeStyles((theme: Theme) => {
  return {
    draggingItem: {
      background: "lightgreen",
      margin: `0 0 ${theme.spacing(1)}px 0`,
      padding: theme.spacing(2),
      userSelect: "none",
    },
    draggingList: {
      background: "lightblue",
      padding: theme.spacing(1),
      width: 250
    },
    item: {
      // background: "gray",
      margin: `0 0 ${theme.spacing(1)}px 0`,
      padding: theme.spacing(2),
      userSelect: "none",
    },
    list: {
      // background: "lightgray",
      padding: theme.spacing(1),
      width: 250
    },
  }
});

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
  Array.from({length: count}, (_, k) => k).map(k => ({
    content: `item ${k}`,
    id: `item-${k}`,
  } as Item));

// tslint:disable-next-line variable-name
export const Preferences: NextPage = () => {
  const classes = useStyles();
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
      Tags
      <DragDropContext onDragEnd={onDragEnd}>
        <Droppable droppableId="droppable">
          {(provided, snapshot) => (
            <List
              {...provided.droppableProps}
              ref={provided.innerRef}
              component="nav"
              className={snapshot.isDraggingOver ? classes.draggingList : classes.list}
            >
              {items.map((item, index) => (
                <Draggable key={item.id} draggableId={item.id} index={index}>
                  {(provided, snapshot) => (
                    <Card
                      ref={provided.innerRef}
                      {...provided.draggableProps}
                      {...provided.dragHandleProps}
                      className={snapshot.isDragging ? classes.draggingItem : classes.item}
                      style={{...provided.draggableProps.style}}
                    >
                      <CardContent>
                        <Typography color="textSecondary" gutterBottom={true}>
                          Word of the Day
                        </Typography>
                      </CardContent>
                    </Card>
                  )}
                </Draggable>
              ))}
              {provided.placeholder}
            </List>
          )}
        </Droppable>
      </DragDropContext>
    </div>
  );
};

export async function getServerSideProps() {
  resetServerContext()
  return {props: {}}
}

export default Preferences;
