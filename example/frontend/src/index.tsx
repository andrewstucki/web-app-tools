import React from "react";
import ReactDOM from "react-dom";
import { Provider } from "react-redux";
import axios from "axios";

import initializeStore from "./store";
import { App } from "./App";

function initialize() {
  const store = initializeStore(axios.create());

  ReactDOM.render(
    <Provider store={store}>
      <React.StrictMode>
        <App />
      </React.StrictMode>
    </Provider>,
    document.getElementById("root")
  );
}

if (process.env.NODE_ENV === "development") {
  const targetHost = process.env.REACT_APP_BASE_URL;
  const currentHost = window.location.protocol + "//" + window.location.host;
  if (targetHost && currentHost !== targetHost) {
    window.location.href = targetHost;
  } else {
    initialize();
  }
} else {
  initialize();
}
