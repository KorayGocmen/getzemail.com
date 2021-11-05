import React from "react";
import { 
  BrowserRouter as Router, 
  Route, 
  Switch, 
} from "react-router-dom";

import Index from "./Index/Index";
import Inbox from "./Inbox/Inbox";
import Message from "./Message/Message";

const Routes = () => {
  return (
    <Router>
      <Switch>
        <Route path="/messages/:messageId">
          <Message />
        </Route>
        <Route path="/:address">
          <Inbox />
        </Route>
        <Route path="/">
          <Index />
        </Route>
      </Switch>
    </Router>
  );
}

export default Routes;