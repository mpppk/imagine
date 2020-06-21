import {NextPage} from 'next';
import React, {useEffect} from 'react';
import {ImageGridList} from "../components/ImageGrid";
import {Button} from "@material-ui/core";

// tslint:disable-next-line variable-name
export const Index: NextPage = () => {
  useEffect(() => {
    const ws = new WebSocket('ws://localhost:1323/ws')

    ws.onopen = function () {
      console.log('Connected')
    }

    ws.onmessage = function (evt) {
      console.log(evt)
    }

    setInterval(function () {
      ws.send('Hello, Server!');
    }, 1000);
  }, [])
  return (
    <div>
      <ImageGridList/>
      <Button variant="outlined" color="primary">
        Edit Query
      </Button>
    </div>
  );
};

export default Index;
