import React from "react";
import Counter from "./components/Counter";
import Header from "./components/Header";

function App(props) {
  console.log("APP rendered", props);
  return (
    <div>
      <Header text={props.Name} />
      <Counter defaultNum={props.InitialNumber} />
    </div>
  );
}

export default App;
